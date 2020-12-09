package model

import "time"

const (
	// CmdHeart 心跳
	CmdHeart = 0
	// CmdNotificationMsg 通知公告(dstid表示通知的群聊，0表示是系统通知)
	CmdNotificationMsg = 1
	// CmdSingleMsg 私聊(dstid是私聊对象id)
	CmdSingleMsg = 10
	// CmdRoomMsg 群聊(dstid是群id)
	CmdRoomMsg = 11
	// CmdUpdateUserInfo 要求客户端重新获取指定用户信息(userid是要获得的用户id)
	CmdUpdateUserInfo = 20

	// CmdRequestOK 请求成功
	CmdRequestOK = 200
	// CmdDropOutGroupOK 退群成功(dstid是要退出的群id)
	CmdDropOutGroupOK = 201
	// CmdGetChatHistoryOK 获取用户聊天记录成功(dstid是要获取的对象的id,amount区分请求对象是群-11还是用户-10)
	CmdGetChatHistoryOK = 202
	// CmdApplyInfoOK 申请加群或加好友成功(dstid表示目标id[加群时为群主id],url存储群id,
	// amount区分请求对象是群-11还是用户-10,content为申请理由)
	CmdApplyInfoOK = 203
	// CmdUpdateMsgOK 更新消息信息成功
	CmdUpdateMsgOK = 204
	// CmdDealApplyOK 加好友/群成功
	CmdDealApplyOK = 205

	// CmdDropOutGroup 用户请求退群(userid是用户id,dstid是想要退的群id)
	CmdDropOutGroup = 301
	// CmdGetChatHistory 用户请求获取聊天记录信息
	// (userid是用户id,dstid是想要获取的对象的id,amount区分请求对象是群-11还是用户-10,content包括开始时间和结束时间)
	CmdGetChatHistory = 302
	// CmdApplyInfo 申请加群或加好友(dstid表示目标id,amount区分请求对象是群-11还是用户-10,content为申请理由)
	CmdApplyInfo = 303
	// CmdUpdateMsg 更新消息信息(通过id更新,所以必须要有id)
	CmdUpdateMsg = 304
	// CmdDealApply 处理请求，加好友(直接修改203为305)，加群(直接修改203为306)
	CmdDealApply = 305

	// CmdRequestFail 请求失败
	CmdRequestFail = 400
)

// Message 消息结构体
type Message struct {
	ID       int64     `xorm:"pk autoincr bigint(20) 'id'" json:"id,omitempty" form:"id"` //消息ID
	Userid   int64     `xorm:"bigint(20)" json:"userid,omitempty" form:"userid"`          //发送者id
	Cmd      int       `xorm:"int(11)" json:"cmd,omitempty" form:"cmd"`                   //消息类型，群聊还是私聊
	Dstid    int64     `xorm:"bigint(20)" json:"dstid,omitempty" form:"dstid"`            //对端用户ID/群ID
	Media    int       `xorm:"int(11)" json:"media,omitempty" form:"media"`               //消息按照什么样式展示
	Content  string    `xorm:"varchar(120)" json:"content,omitempty" form:"content"`      //消息的内容
	Pic      string    `xorm:"varchar(120)" json:"pic,omitempty" form:"pic"`              //预览图片
	URL      string    `xorm:"varchar(120) 'url'" json:"url,omitempty" form:"url"`        //服务的URL
	Memo     string    `xorm:"varchar(120)" json:"memo,omitempty" form:"memo"`            //简单描述
	Amount   int       `xorm:"int(11)" json:"amount,omitempty" form:"amount"`             //其他和数字相关的
	CreateAt time.Time `xorm:"datetime" form:"createat" json:"createat"`                  //消息时间
}
