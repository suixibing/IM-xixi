package service

import (
	"IM-xixi/model"
	"fmt"
	"time"
)

// MessageService 用于使用消息管理服务
type MessageService struct{}

var messageService MessageService

// LoadMsg 加载消息
func (ms MessageService) LoadMsg(dstid int64, cmd, num int) (
	[]*model.Message, error) {
	msgs := make([]*model.Message, 0)
	err := DBEngine.Where("dstid=? and cmd=?",
		dstid, cmd).Desc("create_at").Find(&msgs)
	// num=0 表示加载所有的消息
	if num == 0 {
		num = len(msgs)
	}
	if len(msgs) > num {
		return msgs[:num], err
	}
	return msgs, err
}

// LoadNotificationMsg 加载通知
func (ms MessageService) LoadNotificationMsg(dstid int64, num int) (
	[]*model.Message, error) {
	return ms.LoadMsg(dstid, model.CmdNotificationMsg, num)
}

// LoadChatMsg 加载通信记录
func (ms MessageService) LoadChatMsg(userid, dstid int64, cmd int, starttime, endtime time.Time) (
	[]*model.Message, error) {
	msgs := make([]*model.Message, 0)
	err := fmt.Errorf("还不支持加载记录cmd(%d)", cmd)
	if cmd == model.CmdSingleMsg {
		err = DBEngine.Where(`(((userid=? and dstid=?) or (userid=? and dstid=?)) 
			and cmd=? and (create_at between ? and ?))`,
			userid, dstid, dstid, userid, cmd, starttime.Format("2006-01-02 15:04:05"),
			endtime.Format("2006-01-02 15:04:05")).Asc("create_at").Find(&msgs)
	} else if cmd == model.CmdRoomMsg {
		err = DBEngine.Where("dstid=? and cmd=? and (create_at between ? and ?)",
			dstid, cmd, starttime.Format("2006-01-02 15:04:05"),
			endtime.Format("2006-01-02 15:04:05")).Asc("create_at").Find(&msgs)
	}
	return msgs, err
}

// SaveMsg 保存消息记录
func (ms MessageService) SaveMsg(msg *model.Message) error {
	// msg.CreateAt = time.Now()
	_, err := DBEngine.InsertOne(msg)
	return err
}

// UpdateMsg 更新消息记录
func (ms MessageService) UpdateMsg(msg *model.Message, cols []string) error {
	_, err := DBEngine.ID(msg.ID).Cols(cols...).Update(msg)
	return err
}
