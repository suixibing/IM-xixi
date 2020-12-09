package model

import (
	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
)

// ChatNode chat节点
type ChatNode struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
	GroupSet  set.Interface
}
