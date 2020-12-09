package service

import (
	"IM-xixi/model"
	"fmt"
	"time"

	"github.com/prometheus/common/log"
)

// CommunityService 用于使用群组服务
type CommunityService struct{}

var communityService CommunityService

// GetCommunityByID 通过id获取用户对象
func (cs CommunityService) GetCommunityByID(id int64) (*model.Community, bool) {
	community := &model.Community{}
	_, err := DBEngine.Where("id=?", id).Get(community)
	if err != nil || community.ID <= 0 {
		return nil, false
	}
	return community, true
}

// GetOwnerID 获得指定群聊群主的id
func (cs CommunityService) GetOwnerID(commid int64) (int64, error) {
	community := &model.Community{}
	_, err := DBEngine.ID(commid).Get(community)
	if err != nil || community.ID <= 0 {
		log.Error(err, "| community.ID: ", community.ID)
		return 0, fmt.Errorf("获取群主id失败！")
	}
	return community.Ownerid, nil
}

// GetAllMembersID 获得群组中所有成员的id
func (cs CommunityService) GetAllMembersID(communityid int64) ([]int64, error) {
	uids := make([]int64, 0)
	contacts := make([]*model.Contact, 0)
	err := DBEngine.Where("dstobj=? and cate=?",
		communityid, model.ContactCateCommunity).Find(&contacts)
	if err != nil {
		log.Error(err)
		return uids, err
	}
	for _, contact := range contacts {
		uids = append(uids, contact.Ownerid)
	}
	return uids, nil
}

// CreateCommunity1 创建群聊old
func (cs CommunityService) CreateCommunity1(ownerid int64, cate int,
	name, pic, memo string) (*model.Community, error) {
	community := &model.Community{
		Name:     name,
		Ownerid:  ownerid,
		Icon:     pic,
		Cate:     cate,
		Memo:     memo,
		CreateAt: time.Now(),
	}
	_, err := DBEngine.InsertOne(community)
	return community, err
}

// CreateCommunity 创建群聊
func (cs CommunityService) CreateCommunity(comm *model.Community) (
	*model.Community, error) {
	if comm.Ownerid == 0 {
		return nil, fmt.Errorf("请先登录！")
	}
	if len(comm.Name) == 0 {
		return nil, fmt.Errorf("缺少群名称！")
	}

	comCount := &model.Community{
		Ownerid: comm.Ownerid,
	}
	cnt, err := DBEngine.Count(comCount)
	if err != nil {
		return nil, fmt.Errorf("创建群聊失败！")
	}
	if cnt >= 5 {
		return nil, fmt.Errorf("每个用户最多只能成为5个群的群主！")
	}

	comm.CreateAt = time.Now()
	session := DBEngine.NewSession()
	if err = session.Begin(); err != nil {
		return nil, err
	}
	if _, err = session.InsertOne(comm); err != nil {
		if err = session.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}
	contact := &model.Contact{
		Ownerid:  comm.Ownerid,
		Dstobj:   comm.ID,
		Cate:     model.ContactCateCommunity,
		Memo:     "",
		Createat: comm.CreateAt,
	}
	if _, err = session.InsertOne(contact); err != nil {
		if err = session.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}
	if err = session.Commit(); err != nil {
		return nil, err
	}
	return comm, nil
}
