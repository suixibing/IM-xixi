package service

import (
	"IM-xixi/model"
	"fmt"
	"time"

	"github.com/prometheus/common/log"
)

// ContactService 用于使用用户通讯服务
type ContactService struct{}

var contactService ContactService

// GetContact 获取指定的关系对象
func (cs ContactService) GetContact(
	ownerid, dstobj int64, cate int) (*model.Contact, bool) {
	contact := &model.Contact{}
	_, err := DBEngine.Where("ownerid=? and dstobj=? and cate=?",
		ownerid, dstobj, cate).Get(contact)
	if err != nil || contact.Dstobj <= 0 {
		log.Error(err, "| contact.Dstobj: ", contact.Dstobj)
		return nil, false
	}
	return contact, true
}

// GetContactAll 获取用户的所有关系对象
func (cs ContactService) GetContactAll(ownerid int64, cate int) (
	[]*model.Contact, bool) {
	contacts := make([]*model.Contact, 0)
	err := DBEngine.Where("ownerid=? and cate=?",
		ownerid, cate).Find(&contacts)
	if err != nil {
		log.Error(err)
		return nil, false
	}
	return contacts, true
}

// GetFriendsID 获取用户的所有好友的id
func (cs ContactService) GetFriendsID(ownerid int64) ([]int64, error) {
	contacts, ok := cs.GetContactAll(ownerid, model.ContactCateUser)
	if !ok {
		return nil, fmt.Errorf("获取好友信息失败！")
	}
	ids := make([]int64, 0)
	for _, contact := range contacts {
		ids = append(ids, contact.Dstobj)
	}
	return ids, nil
}

// GetFriends 获取用户的所有好友
func (cs ContactService) GetFriends(ownerid int64) ([]*model.User, error) {
	users := make([]*model.User, 0)
	ids, err := cs.GetFriendsID(ownerid)
	if err != nil || len(ids) == 0 {
		return users, err
	}
	DBEngine.In("id", ids).Find(&users)
	return users, nil
}

// GetCommunitysID 获取用户的所有群组的id
func (cs ContactService) GetCommunitysID(ownerid int64) ([]int64, error) {
	contacts, ok := cs.GetContactAll(ownerid, model.ContactCateCommunity)
	if !ok {
		return nil, fmt.Errorf("获取群组信息失败！")
	}
	ids := make([]int64, 0)
	chatService := ChatService{}
	for _, contact := range contacts {
		ids = append(ids, contact.Dstobj)
		err := chatService.AddCommunityID(ownerid, contact.Dstobj)
		if err != nil {
			return nil, err
		}
	}
	return ids, nil
}

// GetCommunitys 获取用户的所有群组
func (cs ContactService) GetCommunitys(ownerid int64) ([]*model.Community, error) {
	communitys := make([]*model.Community, 0)
	ids, err := cs.GetCommunitysID(ownerid)
	if err != nil || len(ids) == 0 {
		return communitys, err
	}
	DBEngine.In("id", ids).Find(&communitys)
	return communitys, nil
}

// CanAddFriend 对方可以添加为好友
func (cs ContactService) CanAddFriend(ownerid, dstobj int64) error {
	if ownerid == dstobj {
		return fmt.Errorf("不能添加自己为好友！")
	}
	_, exist := userService.GetUserByID(dstobj)
	if !exist {
		return fmt.Errorf("该用户不存在！")
	}
	_, exist = cs.GetContact(ownerid, dstobj, model.ContactCateUser)
	if exist {
		return fmt.Errorf("该用户已经是您的好友！")
	}
	return nil
}

// AddFriend 添加好友
func (cs ContactService) AddFriend(ownerid, dstobj int64) error {
	if err := cs.CanAddFriend(ownerid, dstobj); err != nil {
		return err
	}
	session := DBEngine.NewSession()
	err := session.Begin()
	if err != nil {
		log.Error(err)
		return err
	}
	contact := &model.Contact{
		Ownerid:  ownerid,
		Dstobj:   dstobj,
		Cate:     model.ContactCateUser,
		Memo:     "",
		Createat: time.Now(),
	}

	if _, err = session.InsertOne(contact); err != nil {
		log.Error(err)
		if err := session.Rollback(); err != nil {
			log.Error(err)
			return err
		}
		return err
	}
	contact.ID = 0
	contact.Ownerid, contact.Dstobj = contact.Dstobj, contact.Ownerid
	if _, err = session.InsertOne(contact); err != nil {
		log.Error(err)
		if err := session.Rollback(); err != nil {
			log.Error(err)
			return err
		}
		return err
	}
	if err = session.Commit(); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

// CanJoinCommunity 用户可以加群
func (cs ContactService) CanJoinCommunity(ownerid, dstobj int64) error {
	_, exist := communityService.GetCommunityByID(dstobj)
	if !exist {
		return fmt.Errorf("该群组不存在！")
	}
	_, exist = cs.GetContact(ownerid, dstobj, model.ContactCateCommunity)
	if exist {
		return fmt.Errorf("已经加入该群组！")
	}
	return nil
}

// JoinCommunity 用户加群
func (cs ContactService) JoinCommunity(ownerid, dstobj int64) error {
	if err := cs.CanJoinCommunity(ownerid, dstobj); err != nil {
		return err
	}
	contact := &model.Contact{
		Ownerid:  ownerid,
		Dstobj:   dstobj,
		Cate:     model.ContactCateCommunity,
		Memo:     "",
		Createat: time.Now(),
	}
	if _, err := DBEngine.InsertOne(contact); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

// DropOutCommunity 用户退群
func (cs ContactService) DropOutCommunity(ownerid, dstobj int64) error {
	_, exist := cs.GetContact(ownerid, dstobj, model.ContactCateCommunity)
	if !exist {
		return fmt.Errorf("您没有加入该群组！")
	}

	contact := &model.Contact{
		Ownerid: ownerid,
		Dstobj:  dstobj,
		Cate:    model.ContactCateCommunity,
	}
	if _, err := DBEngine.Cols("ownerid", "dstobj", "cate").Delete(contact); err != nil {
		log.Error(err)
		return err
	}
	return nil
}
