package ctrl

import (
	"log"
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
	userid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Print(err)
		util.RespFail(w, "使用了不合法的id！")
		return
	}

	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return userService.CheckToken(userid, token)
		}}).Upgrade(w, r, nil)
	if err != nil {
		log.Print(err)
		util.RespFail(w, "ws: chat 失败！")
		return
	}

	node := &model.ChatNode{
		Conn:      conn,
		DataQueue: make(chan []byte),
		GroupSet:  set.New(set.ThreadSafe),
	}
	chatService.AddChatNode(userid, node)

	go chatService.Sendproc(node)
	go chatService.Recvproc(node)
	chatService.SendMsg(userid, []byte("hello world"))
	util.RespOK(w, "ws: chat 成功", nil)
}
