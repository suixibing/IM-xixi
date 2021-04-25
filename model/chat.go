package model

import (
	"fmt"

	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
)

// ChatNode chat节点
type ChatNode struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
	GroupSet  set.Interface
}

// String 实现fmt.Stringer接口
func (c ChatNode) String() string {
	return fmt.Sprintf("[%v %v]", c.Conn, c.GroupSet)
}
