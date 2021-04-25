package service

import (
	"fmt"
	"time"

	"github.com/suixibing/IM-xixi/model"
	"github.com/suixibing/IM-xixi/util"
)

// MessageService 用于使用消息管理服务
type MessageService struct{}

var messageService MessageService

// LoadMsg 加载消息
func (ms MessageService) LoadMsg(dstid int64, cmd, num int) (
	[]*model.Message, error) {
	msgs := make([]*model.Message, 0)
	util.GetLog().Debug().Int64("目标id", dstid).Int("消息类型", cmd).Int("数量", num).Msg("加载消息")
	err := DBEngine.Where("dstid=? and cmd=?", dstid, cmd).Desc("create_at").Find(&msgs)
	if err != nil {
		util.GetLog().Error().Err(err).Int64("目标id", dstid).Int("消息类型", cmd).Int("数量", num).
			Msg("加载消息失败")
		return nil, fmt.Errorf("加载消息失败")
	}
	util.GetLog().Trace().Int64("目标id", dstid).Int("消息类型", cmd).Int("数量", num).Msg("加载消息成功")
	// num=0 表示加载所有的消息
	if num == 0 {
		num = len(msgs)
	}
	if len(msgs) > num {
		return msgs[:num], nil
	}
	return msgs, nil
}

// LoadNotificationMsg 加载通知
func (ms MessageService) LoadNotificationMsg(dstid int64, num int) (
	[]*model.Message, error) {
	util.GetLog().Debug().Int64("目标id", dstid).Int("数量", num).Msg("加载通知消息")
	return ms.LoadMsg(dstid, model.CmdNotificationMsg, num)
}

// LoadChatMsg 加载聊天记录
func (ms MessageService) LoadChatMsg(userid, dstid int64, cmd int, starttime, endtime time.Time) (
	[]*model.Message, error) {
	util.GetLog().Debug().Int64("用户id", userid).Int64("目标id", dstid).
		Times("起止时间", []time.Time{starttime, endtime}).Msg("加载聊天记录")
	msgs := make([]*model.Message, 0)
	err := fmt.Errorf("不支持加载记录类型cmd(%d)", cmd)
	if cmd == model.CmdSingleMsg {
		util.GetLog().Debug().Int64("用户id", userid).Int64("目标id", dstid).
			Times("起止时间", []time.Time{starttime, endtime}).Msg("加载私聊记录")
		err = DBEngine.Where(`(((userid=? and dstid=?) or (userid=? and dstid=?)) 
		        and cmd=? and (create_at between ? and ?))`,
			userid, dstid, dstid, userid, cmd, starttime.Format(util.DefaultTimeFormat),
			endtime.Format(util.DefaultTimeFormat)).Asc("create_at").Find(&msgs)
	} else if cmd == model.CmdRoomMsg {
		util.GetLog().Debug().Int64("用户id", userid).Int64("目标id", dstid).
			Times("起止时间", []time.Time{starttime, endtime}).Msg("加载群聊记录")
		err = DBEngine.Where("dstid=? and cmd=? and (create_at between ? and ?)",
			dstid, cmd, starttime.Format(util.DefaultTimeFormat),
			endtime.Format(util.DefaultTimeFormat)).Asc("create_at").Find(&msgs)
	}
	if err != nil {
		util.GetLog().Error().Err(err).Int64("用户id", userid).Int64("目标id", dstid).
			Times("起止时间", []time.Time{starttime, endtime}).Msg("加载聊天记录失败")
		return nil, fmt.Errorf("加载聊天记录失败")
	}
	util.GetLog().Trace().Int64("用户id", userid).Int64("目标id", dstid).
		Times("起止时间", []time.Time{starttime, endtime}).Msg("加载聊天记录成功")
	return msgs, nil
}

// SaveMsg 保存消息记录
func (ms MessageService) SaveMsg(msg *model.Message) error {
	// msg.CreateAt = time.Now()
	util.GetLog().Debug().Stringer("data", msg).Msg("保存消息记录操作")
	_, err := DBEngine.InsertOne(msg)
	if err != nil {
		util.GetLog().Error().Err(err).Stringer("data", msg).Msg("保存消息记录失败")
		return fmt.Errorf("保存消息记录失败")
	}
	util.GetLog().Trace().Stringer("data", msg).Msg("保存消息记录成功")
	return nil
}

// UpdateMsg 更新消息记录
func (ms MessageService) UpdateMsg(msg *model.Message, cols []string) error {
	util.GetLog().Debug().Stringer("data", msg).Msg("更新消息记录操作")
	_, err := DBEngine.ID(msg.ID).Cols(cols...).Update(msg)
	if err != nil {
		util.GetLog().Error().Err(err).Stringer("data", msg).Msg("更新消息记录失败")
		return fmt.Errorf("更新消息记录失败")
	}
	util.GetLog().Debug().Stringer("data", msg).Msg("更新消息记录成功")
	return nil
}
