package model

import "time"

const (
	// SexMan 男性
	SexMan = "M"
	// SexWomen 女性
	SexWomen = "W"
	// SexUnknow 性别未知
	SexUnknow = "U"

	// NotOnline 离线状态
	NotOnline = 0
	// IsOnline 在线状态
	IsOnline = 1
)

// User 用户结构体
type User struct {
	ID       int64     `xorm:"pk autoincr bigint(20) 'id'" form:"id" json:"id"`
	Mobile   string    `xorm:"varchar(20)" form:"mobile" json:"mobile"`
	Passwd   string    `xorm:"varchar(40)" form:"passwd" json:"-"`
	Nickname string    `xorm:"varchar(20)" form:"nickname" json:"nickname"`
	Avatar   string    `xorm:"varchar(150)" form:"avatar" json:"avatar"`
	Sex      string    `xorm:"varchar(2)" form:"sex" json:"sex"`
	Birthday time.Time `xorm:"datetime" form:"birthday" json:"birthday"`
	// 加盐
	Salt   string `xorm:"varchar(10)" form:"salt" json:""`
	Online int    `xorm:"int(10)" form:"online" json:"online"`
	// 前端鉴权因子
	Token    string    `xorm:"varchar(40)" form:"token" json:"token"`
	Memo     string    `xorm:"varchar(140)" form:"memo" json:"memo"`
	CreateAt time.Time `xorm:"datetime" form:"createat" json:"createat"`
}
