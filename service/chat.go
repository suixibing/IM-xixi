package service

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/suixibing/IM-xixi/model"
	"github.com/suixibing/IM-xixi/util"
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
		util.GetLog().Error().Str("error", "websocket未生效").Int64("用户id", userid).Msg("获取用户节点失败")
		return fmt.Errorf("websocket未生效")
	}
	node.GroupSet.Add(dstid)
	util.GetLog().Debug().Int64("用户id", userid).Msg("为用户节点添加群组id成功")
	return nil
}

// DeleteGroupID 为用户删除群组id
func (cs ChatService) DeleteGroupID(userid, dstid int64) error {
	rwlocker.Lock()
	defer rwlocker.Unlock()
	node, ok := clientMap[userid]
	if !ok {
		util.GetLog().Error().Str("error", "websocket未生效").Int64("用户id", userid).Msg("获取用户节点失败")
		return fmt.Errorf("websocket未生效")
	}
	node.GroupSet.Remove(dstid)
	util.GetLog().Debug().Int64("用户id", userid).Msg("为用户删除群组id成功")
	return nil
}

// AddChatNode 添加节点
func (cs ChatService) AddChatNode(userid int64, node *model.ChatNode) {
	rwlocker.Lock()
	clientMap[userid] = node
	rwlocker.Unlock()
	util.GetLog().Debug().Int64("用户id", userid).Stringer("聊天节点", node).Msg("添加聊天节点")
}

// Dispatch 调度函数
func (cs ChatService) Dispatch(data []byte) {
	util.GetLog().Info().Str("data", string(data)).Msg("开始数据分发")

	msg := model.Message{}
	util.GetLog().Debug().Str("data", string(data)).Msg("开始数据json解码")
	err := json.Unmarshal(data, &msg)
	if err != nil {
		util.GetLog().Error().Err(err).Str("data", string(data)).Msg("json解码失败")
		return
	}

	util.GetLog().Debug().Int64("用户id", msg.Userid).Stringer("data", msg).Msg("显示分发的msg数据")
	// 持久化存储聊天记录
	if msg.Cmd == 10 || msg.Cmd == 11 {
		ms := MessageService{}
		err = ms.SaveMsg(&msg)
		if err != nil {
			util.GetLog().Error().Err(err).Int64("用户id", msg.Userid).Msg("聊天记录存储失败")
			return
		}
		util.GetLog().Debug().Int64("用户id", msg.Userid).Msg("聊天记录存储成功")
	}

	switch msg.Cmd {
	case model.CmdHeart:
		// todo
		// 心跳检测用户是否在线
	case model.CmdSingleMsg:
		util.GetLog().Info().Int64("用户id", msg.Userid).Int64("对方id", msg.Dstid).Msg("处理私聊数据")
		cs.SendMsg(msg.Dstid, data)
	case model.CmdRoomMsg:
		util.GetLog().Info().Int64("用户id", msg.Userid).Int64("群聊id", msg.Dstid).Msg("处理群聊数据")
		for userid, node := range clientMap {
			if node.GroupSet.Has(msg.Dstid) && msg.Userid != userid {
				util.GetLog().Debug().Int64("用户id", msg.Userid).
					Int64("群成员id", userid).Msg("发送群聊聊天数据")
				cs.SendMsg(userid, data)
			}
		}
	case model.CmdDropOutGroup:
		util.GetLog().Info().Int64("用户id", msg.Userid).Int64("群聊id", msg.Dstid).Msg("处理退群请求")
		contactService := ContactService{}
		err = contactService.DropOutCommunity(msg.Userid, msg.Dstid)
		if err != nil {
			util.GetLog().Error().Err(err).Int64("退群用户id", msg.Userid).Int64("群聊id", msg.Dstid).
				Msg("退群请求处理失败")
			msg.Cmd = model.CmdRequestFail
			msg.Content = "退群请求处理失败"
		} else {
			util.GetLog().Debug().Int64("退群用户id", msg.Userid).Int64("群聊id", msg.Dstid).Msg("退群请求处理成功")
			msg.Cmd = model.CmdDropOutGroupOK
			for userid, node := range clientMap {
				if node.GroupSet.Has(msg.Dstid) {
					util.GetLog().Debug().Int64("退群用户id", msg.Userid).
						Int64("群成员id", userid).Msg("发送群员退群提醒")
					cs.SendMsg(userid, data)
				}
			}
		}
		util.GetLog().Debug().Int64("退群用户id", msg.Userid).Int64("群聊id", msg.Dstid).Msg("反馈退群请求处理结果")
		cs.SendMsgJSON(msg.Userid, msg)
	case model.CmdGetChatHistory:
		// 这里的Amount暂存了Cmd来决定是私聊还是群聊
		util.GetLog().Info().Int64("用户id", msg.Userid).Int64("聊天对象id", msg.Dstid).
			Int("聊天种类", msg.Amount).Msg("获取聊天历史记录请求")
		messageService := MessageService{}
		end, err := time.Parse(time.RFC3339, msg.Content)
		if err != nil {
			util.GetLog().Error().Err(err).Int64("用户id", msg.Userid).Int64("聊天对象id", msg.Dstid).
				Int("聊天种类", msg.Amount).Msg("时间解析失败")
			msg.Cmd = model.CmdRequestFail
			return
		}
		start := end.AddDate(0, 0, -1)

		msgs, err := messageService.LoadChatMsg(msg.Userid, msg.Dstid, msg.Amount, start, end)
		if err != nil {
			util.GetLog().Error().Err(err).Int64("用户id", msg.Userid).Int64("聊天对象id", msg.Dstid).
				Int("聊天种类", msg.Amount).Msg("加载聊天历史失败")
			msg.Cmd = model.CmdRequestFail
		} else {
			util.GetLog().Debug().Int64("用户id", msg.Userid).Int64("聊天对象id", msg.Dstid).
				Int("聊天种类", msg.Amount).Msg("加载聊天历史成功")
			msg.Cmd = model.CmdGetChatHistoryOK
		}
		// 如果找不到目标内的记录的话就查找过去一年内是否有记录
		if len(msgs) == 0 {
			util.GetLog().Debug().Int64("用户id", msg.Userid).Int64("聊天对象id", msg.Dstid).
				Int("聊天种类", msg.Amount).Msg("获取过去一年内的记录")
			start = start.AddDate(-1, 0, 0)
			msgs, err = messageService.LoadChatMsg(msg.Userid, msg.Dstid, msg.Amount, start, end)
			if err != nil {
				util.GetLog().Error().Err(err).Int64("用户id", msg.Userid).Int64("聊天对象id", msg.Dstid).
					Int("聊天种类", msg.Amount).Msg("加载聊天历史失败")
				msg.Cmd = model.CmdRequestFail
			} else {
				util.GetLog().Debug().Int64("用户id", msg.Userid).Int64("聊天对象id", msg.Dstid).
					Int("聊天种类", msg.Amount).Msg("加载聊天历史成功")
				msg.Cmd = model.CmdGetChatHistoryOK
			}
			// 如果有的话返回最新的n条(默认为20，当最旧的一条记录同一秒内有多条记录时返回可能大于20)记录
			n := 20
			for len(msgs) > n && msgs[len(msgs)-n-1].CreateAt == msgs[len(msgs)-n].CreateAt {
				n++
			}
			util.GetLog().Debug().Int64("用户id", msg.Userid).Int64("聊天对象id", msg.Dstid).
				Int("聊天种类", msg.Amount).Ints("返回数/总数", []int{n, len(msgs)}).Msg("")
			if len(msgs) > n {
				msgs = msgs[len(msgs)-n:]
			}
		}
		data, _ = json.Marshal(msgs)
		msg.Content = string(data)
		util.GetLog().Debug().Int64("用户id", msg.Userid).Int64("聊天对象id", msg.Dstid).
			Int("聊天种类", msg.Amount).Msg("返回聊天记录")
		cs.SendMsgJSON(msg.Userid, msg)
	case model.CmdApplyInfo:
		msg.Cmd = model.CmdApplyInfoOK
		// 检查是否可以申请
		if msg.Amount == 10 {
			util.GetLog().Info().Int64("用户id", msg.Userid).Int64("被申请用户id", msg.Dstid).
				Str("申请种类", "加好友").Msg("申请加好友")
			err = contactService.CanAddFriend(msg.Userid, msg.Dstid)
		} else if msg.Amount == 11 {
			util.GetLog().Info().Int64("用户id", msg.Userid).Int64("群聊id", msg.Dstid).
				Str("申请种类", "加群").Msg("申请加群")
			err = contactService.CanJoinCommunity(msg.Userid, msg.Dstid)
			if err == nil {
				// 群聊的时候先用url存储下目标群id
				msg.URL = fmt.Sprintf("%d", msg.Dstid)
				util.GetLog().Debug().Int64("用户id", msg.Userid).Int64("群聊id", msg.Dstid).
					Str("申请种类", "加群").Msg("开始查询群主id")
				// 如果可以申请加群的话，获取群主的id
				msg.Dstid, err = communityService.GetOwnerID(msg.Dstid)
			}
		} else {
			err = fmt.Errorf("未知的msg.amount！！")
		}
		if err != nil {
			util.GetLog().Error().Err(err).Int64("用户id", msg.Userid).Int64("审核用户id", msg.Dstid).Msg("申请处理失败")
			msg.Cmd = model.CmdRequestFail
			msg.Content = "申请处理失败"
		} else {
			// 如果没有问题就持久化存储msg
			util.GetLog().Debug().Int64("用户id", msg.Userid).Int64("审核用户id", msg.Dstid).
				Msg("可以加对方为好友/获取申请审核用户成功")
			err = messageService.SaveMsg(&msg)
			if err != nil {
				util.GetLog().Error().Err(err).Int64("用户id", msg.Userid).Int64("审核用户id", msg.Dstid).Msg("申请处理失败")
				msg.Cmd = model.CmdRequestFail
				msg.Content = "申请处理失败"
			} else {
				// 通知被申请人(amount-10时为申请对象,amount-11时为群主)
				util.GetLog().Debug().Int64("用户id", msg.Userid).Int64("审核用户id", msg.Dstid).
					Msg("申请请求处理成功，向被申请人发送审核信息")
				cs.SendMsgJSON(msg.Dstid, msg)
			}
		}
		// 通知申请人
		util.GetLog().Debug().Int64("用户id", msg.Userid).Msg("向申请人反馈申请请求处理结果")
		cs.SendMsgJSON(msg.Userid, msg)
	case model.CmdUpdateApplyInfo:
		util.GetLog().Info().Int64("用户id", msg.Userid).Int64("消息id", msg.ID).Msg("更新消息信息")
		msg.Cmd = model.CmdUpdateApplyInfoOK
		err = messageService.UpdateMsg(&msg, []string{"memo"})
		if err != nil {
			util.GetLog().Error().Err(err).Int64("用户id", msg.Userid).Int64("消息id", msg.ID).Msg("更新申请信息错误")
			msg.Cmd = model.CmdRequestFail
			msg.Content = "更新申请信息错误！"
		}
		// 通知审核人审核状态变化
		util.GetLog().Debug().Int64("用户id", msg.Userid).Int64("审核用户id", msg.Dstid).Msg("通知审核人审核状态变化")
		cs.SendMsgJSON(msg.Dstid, msg)
		// 通知申请人审核状态变化
		util.GetLog().Debug().Int64("用户id", msg.Userid).Int64("审核用户id", msg.Dstid).Msg("通知被申请人审核状态变化")
		cs.SendMsgJSON(msg.Userid, msg)
	case model.CmdDealApply:
		if msg.Amount == 10 {
			util.GetLog().Info().Int64("审核人id", msg.Userid).Int64("申请用户id", msg.Dstid).Msg("加好友操作")
			err = contactService.AddFriend(msg.Userid, msg.Dstid)
		} else if msg.Amount == 11 {
			var dstid int64
			util.GetLog().Info().Int64("审核人id", msg.Userid).Int64("申请用户id", msg.Dstid).Msg("加群操作")
			dstid, err = strconv.ParseInt(msg.URL, 10, 64)
			if err == nil {
				util.GetLog().Debug().Int64("用户id", msg.Userid).Int64("申请用户id", msg.Dstid).
					Int64("群聊id", dstid).Msg("加群操作")
				err = contactService.JoinCommunity(msg.Userid, dstid)
			}
		} else {
			err = fmt.Errorf("未知的msg.amount！！")
		}
		if err != nil {
			util.GetLog().Error().Err(err).Int64("用户id", msg.Userid).Int64("申请用户id", msg.Dstid).Msg("加好友/群操作失败")
			msg.Cmd = model.CmdRequestFail
			msg.Content = "加好友/群操作失败"
		} else {
			msg.Cmd = model.CmdDealApplyOK
			if msg.Amount == 11 {
				util.GetLog().Debug().Int64("用户id", msg.Userid).Int64("群聊id", msg.Dstid).
					Msg("将群聊id加入申请人节点")
				// 申请人下线的话才会报错，此时没成功添加群聊id也没事
				chatService.AddCommunityID(msg.Userid, msg.Dstid)
			}
			// 通知申请人
			util.GetLog().Debug().Int64("用户id", msg.Userid).Int64("申请用户id", msg.Dstid).
				Msg("通知申请人申请已通过")
			cs.SendMsgJSON(msg.Userid, msg)
		}
		// 通知审批人(用户/群主)
		util.GetLog().Debug().Int64("用户id", msg.Userid).Int64("申请用户id", msg.Dstid).
			Msg("通知审批人审批处理结果")
		cs.SendMsgJSON(msg.Dstid, msg)
	}
}

// SendMsgJSON 发送json格式的消息
func (cs ChatService) SendMsgJSON(userid int64, v interface{}) {
	util.GetLog().Debug().Int64("目标用户id", userid).Interface("data", v).Msg("json编码")
	data, err := json.Marshal(v)
	if err != nil {
		util.GetLog().Error().Err(err).Int64("目标用户id", userid).Interface("data", v).Msg("json编码失败")
		return
	}
	util.GetLog().Debug().Int64("目标用户id", userid).Str("data", string(data)).Msg("json编码成功，发送json数据")
	cs.SendMsg(userid, data)
}

// SendMsgJSONToUsers 向指定的用户们发送json格式的消息
func (cs ChatService) SendMsgJSONToUsers(userids []int64, v interface{}) {
	util.GetLog().Debug().Ints64("目标用户id", userids).Interface("data", v).Msg("向一组目标用户发送json数据")
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
		util.GetLog().Debug().Int64("接收者id", userid).Str("data", string(msg)).Msg("发送数据")
		node.DataQueue <- msg
	}
}

// SendMsgToUsers 向指定的用户们发送msg
func (cs ChatService) SendMsgToUsers(userids []int64, data []byte) {
	util.GetLog().Debug().Ints64("目标用户id", userids).Str("data", string(data)).Msg("向一组目标用户发送数据")
	for _, userid := range userids {
		cs.SendMsg(userid, data)
	}
}

// Sendproc 后台管理发送数据的函数
func (cs ChatService) Sendproc(node *model.ChatNode) {
	for {
		select {
		case data := <-node.DataQueue:
			util.GetLog().Debug().Str("data", string(data)).Msg("开始发送数据")
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				util.GetLog().Error().Err(err).Str("data", string(data)).Msg("发送数据失败")
				return
			}
			util.GetLog().Debug().Str("data", string(data)).Msg("发送数据成功")
		}
	}
}

// Recvproc 后台管理接收数据的函数
func (cs ChatService) Recvproc(node *model.ChatNode) {
	for {
		_, data, err := node.Conn.ReadMessage()
		util.GetLog().Debug().Msg("开始接收数据")
		if err != nil {
			util.GetLog().Error().Err(err).Msg("接收数据失败")
			return
		}
		util.GetLog().Debug().Str("data", string(data)).Msg("接收数据成功")
		cs.Dispatch(data)
	}
}
