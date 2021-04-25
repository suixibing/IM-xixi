package ctrl

import (
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/suixibing/IM-xixi/model"
	"github.com/suixibing/IM-xixi/service"
	"github.com/suixibing/IM-xixi/util"
	"gopkg.in/fatih/set.v0"
)

var chatService service.ChatService

// Chat 聊天功能函数
func Chat(w http.ResponseWriter, r *http.Request) {
	querys := r.URL.Query()
	id, token := querys.Get("id"), querys.Get("token")
	util.GetLog().Debug().Str("用户id", id).Str("token", token).Msg("解析id为数字类型")
	userid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		util.GetLog().Error().Err(err).Str("用户id", id).Msg("使用了不合法的id")
		util.RespFail(w, "使用了不合法的id")
		return
	}
	util.GetLog().Trace().Int64("用户id", userid).Msg("解析id为数字类型成功")

	util.GetLog().Info().Int64("用户id", userid).Msg("使用websocket协议进行通讯")
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			util.GetLog().Debug().Int64("用户id", userid).Str("token", token).Msg("token校验")
			ok := userService.CheckToken(userid, token)
			if !ok {
				util.GetLog().Warn().Int64("用户id", userid).Str("token", token).Msg("token校验不通过")
				return false
			}
			util.GetLog().Trace().Int64("用户id", userid).Str("token", token).Msg("token校验通过")
			return true
		}}).Upgrade(w, r, nil)
	if err != nil {
		util.GetLog().Error().Err(err).Int64("用户id", userid).Msg("使用websocket协议进行通讯失败")
		util.RespFail(w, "使用websocket协议进行通讯失败")
		return
	}
	util.GetLog().Trace().Int64("用户id", userid).Msg("使用websocket协议进行通讯成功")

	node := &model.ChatNode{
		Conn:      conn,
		DataQueue: make(chan []byte),
		GroupSet:  set.New(set.ThreadSafe),
	}
	chatService.AddChatNode(userid, node)

	util.GetLog().Debug().Int64("用户id", userid).Msg("启动发送数据goroutine")
	go chatService.Sendproc(node)
	util.GetLog().Debug().Int64("用户id", userid).Msg("启动接收数据goroutine")
	go chatService.Recvproc(node)
	util.GetLog().Debug().Int64("用户id", userid).Msg("发送测试数据")
	chatService.SendMsg(userid, []byte("hello world"))
}
