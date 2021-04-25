package service

import (
	"fmt"
	"time"

	"github.com/suixibing/IM-xixi/model"
	"github.com/suixibing/IM-xixi/util"
)

// ContactService 用于使用用户通讯服务
type ContactService struct{}

var contactService ContactService

// GetContact 获取指定的关系对象
func (cs ContactService) GetContact(
	ownerid, dstobj int64, cate int) (*model.Contact, error) {
	contact := &model.Contact{}
	util.GetLog().Debug().Int64("用户id", ownerid).Int64("目标id", dstobj).
		Int("关系类别", cate).Msg("获取指定的关系对象")
	_, err := DBEngine.Where("ownerid=? and dstobj=? and cate=?", ownerid, dstobj, cate).Get(contact)
	if err != nil || contact.Dstobj <= 0 {
		util.GetLog().Error().Err(err).Int64("用户id", ownerid).Int64("目标id", contact.Dstobj).
			Int("关系类别", cate).Msg("获取关系对象失败")
		return nil, fmt.Errorf("获取关系对象失败")
	}
	util.GetLog().Trace().Int64("用户id", ownerid).Int64("目标id", dstobj).
		Int("关系类别", cate).Msg("获取关系对象成功")
	return contact, nil
}

// GetContactAll 获取用户的所有关系对象
func (cs ContactService) GetContactAll(ownerid int64, cate int) (
	[]*model.Contact, error) {
	contacts := make([]*model.Contact, 0)
	util.GetLog().Debug().Int64("用户id", ownerid).Int("关系类别", cate).Msg("获取用户的所有关系对象")
	err := DBEngine.Where("ownerid=? and cate=?", ownerid, cate).Find(&contacts)
	if err != nil {
		util.GetLog().Error().Err(err).Int64("用户id", ownerid).Int("关系类别", cate).
			Msg("获取用户所有关系对象失败")
		return nil, fmt.Errorf("获取用户所有关系对象失败")
	}
	util.GetLog().Trace().Int64("用户id", ownerid).Int("关系类别", cate).Msg("获取用户所有关系对象成功")
	return contacts, nil
}

// GetFriendsID 获取用户的所有好友的id
func (cs ContactService) GetFriendsID(ownerid int64) ([]int64, error) {
	util.GetLog().Debug().Int64("用户id", ownerid).Msg("获取所有的好友关系")
	contacts, err := cs.GetContactAll(ownerid, model.ContactCateUser)
	if err != nil {
		util.GetLog().Error().Err(err).Int64("用户id", ownerid).Msg("获取所有的好友关系失败")
		return nil, fmt.Errorf("获取所有的好友关系失败")
	}
	ids := make([]int64, 0)
	for _, contact := range contacts {
		ids = append(ids, contact.Dstobj)
	}
	util.GetLog().Trace().Int64("用户id", ownerid).Ints64("好友id", ids).Msg("获取所有的好友关系成功")
	return ids, nil
}

// GetFriends 获取用户的所有好友
func (cs ContactService) GetFriends(ownerid int64) ([]*model.User, error) {
	users := make([]*model.User, 0)
	util.GetLog().Debug().Int64("用户id", ownerid).Msg("获取所有的好友对象")
	ids, err := cs.GetFriendsID(ownerid)
	if err != nil || len(ids) == 0 {
		util.GetLog().Error().Err(err).Int64("用户id", ownerid).Msg("获取所有的好友id失败")
		return users, fmt.Errorf("获取所有的好友id失败")
	}
	util.GetLog().Trace().Int64("用户id", ownerid).Ints64("好友id", ids).Msg("获取所有的好友id成功")
	util.GetLog().Debug().Int64("用户id", ownerid).Ints64("好友id", ids).Msg("获取所有的好友对象")
	if err = DBEngine.In("id", ids).Find(&users); err != nil {
		util.GetLog().Error().Err(err).Int64("用户id", ownerid).Ints64("好友id", ids).Msg("获取所有的好友对象失败")
		return nil, fmt.Errorf("获取所有的好友对象失败")
	}
	util.GetLog().Trace().Int64("用户id", ownerid).Ints64("好友id", ids).Msg("获取所有的好友对象成功")
	return users, nil
}

// GetCommunitysID 获取用户的所有群聊的id
func (cs ContactService) GetCommunitysID(ownerid int64) ([]int64, error) {
	util.GetLog().Debug().Int64("用户id", ownerid).Msg("获取用户的所有群聊的id")
	contacts, err := cs.GetContactAll(ownerid, model.ContactCateCommunity)
	if err != nil {
		util.GetLog().Error().Err(err).Int64("用户id", ownerid).Msg("获取用户的所有群聊的id失败")
		return nil, fmt.Errorf("获取用户的所有群聊的id失败")
	}
	util.GetLog().Trace().Int64("用户id", ownerid).Msg("获取用户的所有群聊的id成功")
	ids := make([]int64, 0)
	chatService := ChatService{}
	for _, contact := range contacts {
		ids = append(ids, contact.Dstobj)
		util.GetLog().Debug().Int64("用户id", ownerid).Int64("群聊id", contact.Dstobj).Msg("向用户节点添加群聊id")
		err := chatService.AddCommunityID(ownerid, contact.Dstobj)
		if err != nil {
			util.GetLog().Error().Err(err).Int64("用户id", ownerid).Msg("向用户节点添加群聊id失败")
			return nil, fmt.Errorf("向用户节点添加群聊id失败")
		}
		util.GetLog().Trace().Int64("用户id", ownerid).Int64("群聊id", contact.Dstobj).Msg("向用户节点添加群聊id成功")
	}
	return ids, nil
}

// GetCommunitys 获取用户的所有群组
func (cs ContactService) GetCommunitys(ownerid int64) ([]*model.Community, error) {
	communitys := make([]*model.Community, 0)
	util.GetLog().Debug().Int64("用户id", ownerid).Msg("获取用户的所有群聊")
	ids, err := cs.GetCommunitysID(ownerid)
	if err != nil || len(ids) == 0 {
		util.GetLog().Error().Err(err).Int64("用户id", ownerid).Msg("获取用户的所有群聊失败")
		return communitys, fmt.Errorf("获取用户的所有群聊失败")
	}
	util.GetLog().Trace().Int64("用户id", ownerid).Msg("获取用户所有的群聊成功")
	if err = DBEngine.In("id", ids).Find(&communitys); err != nil {
		util.GetLog().Error().Err(err).Int64("用户id", ownerid).Msg("获取用户所有的群聊失败")
		return nil, fmt.Errorf("获取用户所有的群聊失败")
	}
	util.GetLog().Trace().Int64("用户id", ownerid).Msg("获取用户所有的群聊成功")
	return communitys, nil
}

// CanAddFriend 对方可以添加为好友
func (cs ContactService) CanAddFriend(ownerid, dstobj int64) error {
	util.GetLog().Debug().Int64("用户id", ownerid).Int64("对方id", dstobj).Msg("判断能否加对方为好友")
	if ownerid == dstobj {
		util.GetLog().Warn().Int64("用户id", ownerid).Int64("对方id", dstobj).Msg("不能添加自己为好友")
		return fmt.Errorf("不能添加自己为好友")
	}
	_, err := userService.GetUserByID(dstobj)
	if err != nil {
		util.GetLog().Error().Err(err).Int64("用户id", ownerid).Int64("对方id", dstobj).Msg("该用户不存在")
		return fmt.Errorf("该用户不存在")
	}
	_, err = cs.GetContact(ownerid, dstobj, model.ContactCateUser)
	if err != nil {
		util.GetLog().Error().Err(err).Int64("用户id", ownerid).Int64("对方id", dstobj).Msg("该用户已经是您的好友")
		return fmt.Errorf("该用户已经是您的好友")
	}
	util.GetLog().Trace().Int64("用户id", ownerid).Int64("对方id", dstobj).Msg("能加对方为好友")
	return nil
}

// AddFriend 添加好友
func (cs ContactService) AddFriend(ownerid, dstobj int64) error {
	util.GetLog().Debug().Int64("用户id", ownerid).Int64("对象id", dstobj).Msg("加好友操作")
	if err := cs.CanAddFriend(ownerid, dstobj); err != nil {
		util.GetLog().Error().Err(err).Int64("用户id", ownerid).Int64("对象id", dstobj).Msg("不能加对方为好友")
		return fmt.Errorf("不能加对方为好友")
	}
	util.GetLog().Debug().Int64("用户id", ownerid).Int64("对象id", dstobj).Msg("创建数据库会话")
	session := DBEngine.NewSession()
	err := session.Begin()
	if err != nil {
		util.GetLog().Error().Err(err).Int64("用户id", ownerid).Int64("对象id", dstobj).Msg("创建数据库会话失败")
		return fmt.Errorf("创建数据库会话失败")
	}
	contact := &model.Contact{
		Ownerid:  ownerid,
		Dstobj:   dstobj,
		Cate:     model.ContactCateUser,
		Memo:     "",
		Createat: time.Now(),
	}

	util.GetLog().Debug().Int64("用户id", ownerid).Int64("对象id", dstobj).Msg("[会话]插入好友正向关系数据")
	if _, err = session.InsertOne(contact); err != nil {
		util.GetLog().Error().Err(err).Int64("用户id", ownerid).Int64("对象id", dstobj).Msg("[会话]插入好友正向关系数据失败")
		if err := session.Rollback(); err != nil {
			util.GetLog().Error().Err(err).Msg("[会话]数据库回滚失败")
			return fmt.Errorf("[会话]插入好友正向关系数据失败")
		}
		util.GetLog().Trace().Int64("用户id", ownerid).Int64("对象id", dstobj).Msg("[会话]插入好友正向关系数据成功")
		return fmt.Errorf("[会话]插入好友正向关系数据失败")
	}
	contact.ID = 0
	contact.Ownerid, contact.Dstobj = contact.Dstobj, contact.Ownerid
	util.GetLog().Debug().Int64("用户id", ownerid).Int64("对象id", dstobj).Msg("[会话]插入好友逆向关系数据")
	if _, err = session.InsertOne(contact); err != nil {
		util.GetLog().Error().Err(err).Int64("用户id", ownerid).Int64("对象id", dstobj).Msg("[会话]插入好友逆向关系数据失败")
		if err := session.Rollback(); err != nil {
			util.GetLog().Error().Err(err).Msg("[会话]数据库回滚失败")
			return fmt.Errorf("[会话]插入好友逆向关系数据失败")
		}
		util.GetLog().Trace().Int64("用户id", ownerid).Int64("对象id", dstobj).Msg("[会话]插入好友逆向关系数据成功")
		return fmt.Errorf("[会话]插入好友逆向关系数据失败")
	}
	if err = session.Commit(); err != nil {
		util.GetLog().Error().Err(err).Int64("用户id", ownerid).Int64("对象id", dstobj).Msg("[会话]提交会话失败")
		return fmt.Errorf("[会话]提交会话失败")
	}
	util.GetLog().Trace().Int64("用户id", ownerid).Int64("对象id", dstobj).Msg("[会话]提交会话成功")
	return nil
}

// CanJoinCommunity 用户可以加群
func (cs ContactService) CanJoinCommunity(ownerid, dstobj int64) error {
	util.GetLog().Debug().Int64("用户id", ownerid).Int64("群聊id", dstobj).Msg("检查是否可以申请加群")
	_, err := communityService.GetCommunityByID(dstobj)
	if err != nil {
		util.GetLog().Error().Err(err).Int64("用户id", ownerid).Int64("群聊id", dstobj).Msg("该群聊不存在")
		return fmt.Errorf("该群聊不存在")
	}
	_, err = cs.GetContact(ownerid, dstobj, model.ContactCateCommunity)
	if err != nil {
		util.GetLog().Error().Err(err).Int64("用户id", ownerid).Int64("群聊id", dstobj).Msg("已经加入该群聊")
		return fmt.Errorf("已经加入该群聊")
	}
	util.GetLog().Trace().Int64("用户id", ownerid).Int64("群聊id", dstobj).Msg("可以申请加入该群聊")
	return nil
}

// JoinCommunity 用户加群
func (cs ContactService) JoinCommunity(ownerid, dstobj int64) error {
	util.GetLog().Debug().Int64("用户id", ownerid).Int64("群聊id", dstobj).Msg("用户加群操作")
	if err := cs.CanJoinCommunity(ownerid, dstobj); err != nil {
		util.GetLog().Error().Err(err).Int64("用户id", ownerid).Int64("群聊id", dstobj).Msg("用户无法申请加入该群")
		return fmt.Errorf("用户无法申请加入该群")
	}
	contact := &model.Contact{
		Ownerid:  ownerid,
		Dstobj:   dstobj,
		Cate:     model.ContactCateCommunity,
		Memo:     "",
		Createat: time.Now(),
	}
	util.GetLog().Debug().Int64("用户id", ownerid).Int64("群聊id", dstobj).Msg("插入用户群聊关系数据")
	if _, err := DBEngine.InsertOne(contact); err != nil {
		util.GetLog().Error().Err(err).Int64("用户id", ownerid).Int64("群聊id", dstobj).Msg("插入用户群聊关系数据失败")
		return fmt.Errorf("插入用户群聊关系数据失败")
	}
	util.GetLog().Trace().Int64("用户id", ownerid).Int64("群聊id", dstobj).Msg("插入用户群聊关系数据成功")
	return nil
}

// DropOutCommunity 用户退群
func (cs ContactService) DropOutCommunity(ownerid, dstobj int64) error {
	util.GetLog().Debug().Int64("用户id", ownerid).Int64("群聊id", dstobj).Msg("用户退群操作")
	_, err := cs.GetContact(ownerid, dstobj, model.ContactCateCommunity)
	if err != nil {
		util.GetLog().Error().Err(err).Int64("用户id", ownerid).Int64("群聊id", dstobj).Msg("获取群关系出错")
		return fmt.Errorf("获取群关系出错")
	}
	util.GetLog().Trace().Int64("用户id", ownerid).Int64("群聊id", dstobj).Msg("获取群关系成功")

	contact := &model.Contact{
		Ownerid: ownerid,
		Dstobj:  dstobj,
		Cate:    model.ContactCateCommunity,
	}
	util.GetLog().Debug().Int64("用户id", ownerid).Int64("群聊id", dstobj).Msg("删除用户群聊关系数据")
	if _, err := DBEngine.Cols("ownerid", "dstobj", "cate").Delete(contact); err != nil {
		util.GetLog().Error().Err(err).Int64("用户id", ownerid).Int64("群聊id", dstobj).Msg("删除用户群聊关系数据失败")
		return fmt.Errorf("删除用户群聊关系数据失败")
	}
	util.GetLog().Trace().Int64("用户id", ownerid).Int64("群聊id", dstobj).Msg("删除用户群聊关系数据成功")
	//todo 将群id从用户节点中删除
	return nil
}
