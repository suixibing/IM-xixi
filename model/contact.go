package model

import (
	"fmt"
	"time"

	"github.com/suixibing/IM-xixi/util"
)

const (
	// ContactCateUser 用户
	ContactCateUser = 0x01
	// ContactCateCommunity 群组
	ContactCateCommunity = 0x02
)

// Contact 通讯记录表
type Contact struct {
	ID       int64     `xorm:"pk autoincr bigint(20) 'id'" form:"id" json:"id"`
	Ownerid  int64     `xorm:"bigint(20)" form:"ownerid" json:"ownerid"`
	Dstobj   int64     `xorm:"bigint(20)" form:"dstobj" json:"dstobj"`
	Cate     int       `xorm:"int(11)" form:"cate" json:"cate"`
	Memo     string    `xorm:"varchar(120)" form:"memo" json:"memo"`
	Createat time.Time `xorm:"datetime" form:"createat" json:"createat"`
}

// String 实现fmt.Stringer接口
func (c Contact) String() string {
	return fmt.Sprintf("{id: %d, ownerid: %d, dstobj: %d, cate: %d, memo: %s, createat: %v}",
		c.ID, c.Ownerid, c.Dstobj, c.Cate, c.Memo, c.Createat.Format(util.DefaultTimeFormat))
}
