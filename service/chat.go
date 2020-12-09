package service

import (
	"IM-xixi/model"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// ChatService 用于使用通讯服务
type ChatService struct{}

var chatService ChatService

var clientMap map[int64]*model.ChatNode = make(map[int64]*model.ChatNode)
var rwlocker sync.RWMutex

// AddCommunityID 为用户添加群组id
func (cs ChatService) AddCommunityID(userid, dstid int64) error {
	rwlocker.Lock()
	defer rwlocker.Unlock()
	node, ok := clientMap[userid]
	if !ok {
		return fmt.Errorf("userid[%d] websocket未生效！", userid)
	}
	node.GroupSet.Add(dstid)
	return nil
}

// AddChatNode 添加节点
func (cs ChatService) AddChatNode(userid int64, node *model.ChatNode) {
	rwlocker.Lock()
	clientMap[userid] = node
	rwlocker.Unlock()
}

// Dispatch 调度函数
func (cs ChatService) Dispatch(data []byte) {
	msg := model.Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		log.Println(err.Error())
		return
	}
	// 持久化存储聊天记录
	if msg.Cmd == 10 || msg.Cmd == 11 {
		ms := MessageService{}
		err = ms.SaveMsg(&msg)
		if err != nil {
			log.Println(err.Error())
			return
		}
	}

	switch msg.Cmd {
	case model.CmdHeart:
	case model.CmdSingleMsg:
		cs.SendMsg(msg.Dstid, data)
	case model.CmdRoomMsg:
		for userid, node := range clientMap {
			if node.GroupSet.Has(msg.Dstid) && msg.Userid != userid {
				cs.SendMsg(userid, data)
			}
		}
	case model.CmdDropOutGroup:
		contactService := ContactService{}
		err = contactService.DropOutCommunity(msg.Userid, msg.Dstid)
		if err != nil {
			msg.Cmd = model.CmdRequestFail
			msg.Content = err.Error()
		} else {
			msg.Cmd = model.CmdDropOutGroupOK
		}
		cs.SendMsgJSON(msg.Userid, msg)
	case model.CmdGetChatHistory:
		messageService := MessageService{}
		end, err := time.Parse(time.RFC3339, msg.Content)
		if err != nil {
			msg.Cmd = model.CmdRequestFail
			log.Println(err)
			return
		}
		start := end.AddDate(0, 0, -1)

		msgs, err := messageService.LoadChatMsg(msg.Userid, msg.Dstid, msg.Amount, start, end)
		if err != nil {
			msg.Cmd = model.CmdRequestFail
		} else {
			msg.Cmd = model.CmdGetChatHistoryOK
		}
		// 如果找不到目标内的记录的话就查找过去一年内是否有记录
		if len(msgs) == 0 {
			start = start.AddDate(-1, 0, 0)
			msgs, err = messageService.LoadChatMsg(msg.Userid, msg.Dstid, msg.Amount, start, end)
			if err != nil {
				msg.Cmd = model.CmdRequestFail
			} else {
				msg.Cmd = model.CmdGetChatHistoryOK
			}
			// 如果有的话返回最新的n条(默认为20，当最旧的一条记录同一秒内有多条记录时返回可能大于20)记录
			n := 20
			for len(msgs) > n && msgs[len(msgs)-n-1].CreateAt == msgs[len(msgs)-n].CreateAt {
				n++
			}
			if len(msgs) > n {
				msgs = msgs[len(msgs)-n:]
			}
		}
		data, _ = json.Marshal(msgs)
		msg.Content = string(data)
		cs.SendMsgJSON(msg.Userid, msg)
	case model.CmdApplyInfo:
		msg.Cmd = model.CmdApplyInfoOK
		// 检查是否可以申请
		if msg.Amount == 10 {
			err = contactService.CanAddFriend(msg.Userid, msg.Dstid)
		} else if msg.Amount == 11 {
			err = contactService.CanJoinCommunity(msg.Userid, msg.Dstid)
			if err == nil {
				// 群聊的时候先用url存储下目标群id
				msg.URL = fmt.Sprintf("%d", msg.Dstid)
				// 如果可以申请加群的话，获取群主的id
				msg.Dstid, err = communityService.GetOwnerID(msg.Dstid)
			}
		} else {
			err = fmt.Errorf("未知的msg.amount！！")
		}
		if err != nil {
			msg.Cmd = model.CmdRequestFail
			msg.Content = err.Error()
		} else {
			// 如果没有问题就持久化存储msg
			err = messageService.SaveMsg(&msg)
			if err != nil {
				msg.Cmd = model.CmdRequestFail
				msg.Content = err.Error()
				// msg.Content = "申请未成功，请确认无误后再次尝试！"
			} else {
				// 通知被申请人(amount-10时为申请对象,amount-11时为群主)
				cs.SendMsgJSON(msg.Dstid, msg)
			}
		}
		// 通知申请人
		cs.SendMsgJSON(msg.Userid, msg)
	case model.CmdUpdateMsg:
		msg.Cmd = model.CmdUpdateMsgOK
		err = messageService.UpdateMsg(&msg, []string{"memo"})
		if err != nil {
			msg.Cmd = model.CmdRequestFail
			msg.Content = "更新申请信息错误！"
		}
		// 通知批准人
		cs.SendMsgJSON(msg.Dstid, msg)
		// 通知申请人
		cs.SendMsgJSON(msg.Userid, msg)
	case model.CmdDealApply:
		if msg.Amount == 10 {
			err = contactService.AddFriend(msg.Userid, msg.Dstid)
		} else if msg.Amount == 11 {
			var dstid int64
			dstid, err = strconv.ParseInt(msg.URL, 10, 64)
			if err == nil {
				err = contactService.JoinCommunity(msg.Userid, dstid)
			}
		} else {
			err = fmt.Errorf("未知的msg.amount！！")
		}
		if err != nil {
			msg.Cmd = model.CmdRequestFail
			msg.Content = err.Error()
		} else {
			msg.Cmd = model.CmdDealApplyOK
			if msg.Amount == 11 {
				// 对方可能不在线，此时不需要处理错误判断
				chatService.AddCommunityID(msg.Userid, msg.Dstid)
			}
			// 通知申请人
			cs.SendMsgJSON(msg.Userid, msg)
		}
		// 通知批准人(用户/群主)
		cs.SendMsgJSON(msg.Dstid, msg)
	}
}

// SendMsgJSON 发送json格式的消息
func (cs ChatService) SendMsgJSON(userid int64, v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		return
	}
	cs.SendMsg(userid, data)
}

// SendMsgJSONToUsers 向指定的用户们发送json格式的消息
func (cs ChatService) SendMsgJSONToUsers(userids []int64, v interface{}) {
	for _, userid := range userids {
		cs.SendMsgJSON(userid, v)
	}
}

// SendMsg 发送消息
func (cs ChatService) SendMsg(userid int64, msg []byte) {
	rwlocker.Lock()
	node, ok := clientMap[userid]
	rwlocker.Unlock()

	if ok {
		fmt.Printf("send to user[%d] %s\n", userid, msg)
		node.DataQueue <- msg
	}
}

// SendMsgToUsers 向指定的用户们发送msg
func (cs ChatService) SendMsgToUsers(userids []int64, msg []byte) {
	for _, userid := range userids {
		cs.SendMsg(userid, msg)
	}
}

// Sendproc 后台管理发送数据的函数
func (cs ChatService) Sendproc(node *model.ChatNode) {
	for {
		select {
		case data := <-node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Print(err)
				return
			}
		}
	}
}

// Recvproc 后台管理接收数据的函数
func (cs ChatService) Recvproc(node *model.ChatNode) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			log.Print(err)
			return
		}
		log.Printf("recv %s\n", data)
		cs.Dispatch(data)
		// node.DataQueue <- []byte(fmt.Sprintf("recv<=%s", data))
	}
}
