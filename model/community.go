package model

import "time"

// Community 群组结构体
type Community struct {
	ID       int64     `xorm:"pk autoincr bigint(20) 'id'" form:"id" json:"id"`
	Name     string    `xorm:"varchar(30)" form:"name" json:"name"`
	Ownerid  int64     `xorm:"bigint(20)" form:"ownerid" json:"ownerid"`
	Icon     string    `xorm:"varchar(250)" form:"icon" json:"icon"`
	Cate     int       `xorm:"int(11)" form:"cate" json:"cate"`
	Memo     string    `xorm:"varchar(120)" form:"memo" json:"memo"`
	CreateAt time.Time `xorm:"datetime" form:"createat" json:"createat"`
}

// CommunityData 群组的相关信息
type CommunityData struct {
	*Community
	Memids []int64 `form:"memids" json:"memids"`
}

const (
	// CommunityCateDefault 默认
	CommunityCateDefault = 0x00
	// CommunityCateHobby 兴趣爱好
	CommunityCateHobby = 0x01
	// CommunityCateBusiness 行业交流
	CommunityCateBusiness = 0x02
	// CommunityCateLife 生活休闲
	CommunityCateLife = 0x03
	// CommunityCateStudy 学习考试
	CommunityCateStudy = 0x04
)
