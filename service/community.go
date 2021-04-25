package service

import (
	"fmt"
	"time"

	"github.com/suixibing/IM-xixi/model"
	"github.com/suixibing/IM-xixi/util"
)

// CommunityService 用于使用群组服务
type CommunityService struct{}

var communityService CommunityService

// GetCommunityByID 通过id获取用户对象
func (cs CommunityService) GetCommunityByID(id int64) (*model.Community, error) {
	community := &model.Community{}
	util.GetLog().Debug().Int64("群聊id", id).Msg("获取群聊对象")
	_, err := DBEngine.Where("id=?", id).Get(community)
	if err != nil || community.ID <= 0 {
		util.GetLog().Error().Err(err).Int64("群聊id", id).Msg("获取群聊对象失败")
		return nil, err
	}
	util.GetLog().Trace().Int64("群聊id", id).Msg("获取群聊对象成功")
	return community, nil
}

// GetOwnerID 获得指定群聊群主的id
func (cs CommunityService) GetOwnerID(gid int64) (int64, error) {
	util.GetLog().Debug().Int64("群聊id", gid).Msg("获取群聊的群主id")
	community := &model.Community{}
	_, err := DBEngine.ID(gid).Get(community)
	if err != nil || community.ID <= 0 {
		util.GetLog().Error().Err(err).Int64("群聊id", gid).Msg("获取群聊信息失败")
		return 0, err
	}
	util.GetLog().Trace().Int64("群聊id", gid).Int64("群主id", community.Ownerid).Msg("获取群主id成功")
	return community.Ownerid, nil
}

// GetAllMembersID 获得群组中所有成员的id
func (cs CommunityService) GetAllMembersID(gid int64) ([]int64, error) {
	uids := make([]int64, 0)
	contacts := make([]*model.Contact, 0)
	util.GetLog().Debug().Int64("群聊id", gid).Msg("获取群组中所有的成员对象")
	err := DBEngine.Where("dstobj=? and cate=?", gid, model.ContactCateCommunity).Find(&contacts)
	if err != nil {
		util.GetLog().Error().Err(err).Int64("群聊id", gid).Msg("获取群成员失败")
		return uids, err
	}
	for _, contact := range contacts {
		uids = append(uids, contact.Ownerid)
	}
	util.GetLog().Trace().Int64("群聊id", gid).Ints64("群员id", uids).Msg("获取群成员成功")
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
		util.GetLog().Warn().Int64("用户id", comm.Ownerid).Msg("用户未登陆")
		return nil, fmt.Errorf("用户未登陆")
	}
	if len(comm.Name) == 0 {
		util.GetLog().Warn().Int64("用户id", comm.Ownerid).Msg("缺少群名称")
		return nil, fmt.Errorf("缺少群名称")
	}

	comCount := &model.Community{
		Ownerid: comm.Ownerid,
	}
	util.GetLog().Debug().Int64("用户id", comm.Ownerid).Msg("获取用户已经拥有的群数量")
	cnt, err := DBEngine.Count(comCount)
	if err != nil {
		util.GetLog().Error().Err(err).Int64("用户id", comm.Ownerid).Msg("获取用户已经拥有的群数量失败")
		return nil, fmt.Errorf("创建群聊失败！")
	}
	util.GetLog().Trace().Int64("用户id", comm.Ownerid).Msg("获取用户已经拥有的群数量成功")
	if cnt >= 5 {
		util.GetLog().Warn().Int64("用户id", comm.Ownerid).Msg("拥有的群数量超过限制")
		return nil, fmt.Errorf("拥有的群数量超过限制")
	}
	util.GetLog().Trace().Int64("用户id", comm.Ownerid).Msg("可以继续创建群聊")

	comm.CreateAt = time.Now()
	util.GetLog().Debug().Int64("用户id", comm.Ownerid).Msg("创建数据库会话")
	session := DBEngine.NewSession()
	if err = session.Begin(); err != nil {
		util.GetLog().Error().Err(err).Int64("用户id", comm.Ownerid).Msg("创建数据库会话失败")
		return nil, err
	}
	util.GetLog().Debug().Int64("用户id", comm.Ownerid).Msg("[会话]插入群聊信息")
	if _, err = session.InsertOne(comm); err != nil {
		util.GetLog().Error().Err(err).Int64("用户id", comm.Ownerid).Msg("[会话]插入群聊信息失败")
		if err = session.Rollback(); err != nil {
			util.GetLog().Error().Err(err).Int64("用户id", comm.Ownerid).Msg("[会话]数据库回滚失败")
			return nil, err
		}
		util.GetLog().Trace().Int64("用户id", comm.Ownerid).Msg("[会话]数据库回滚成功")
		return nil, err
	}
	util.GetLog().Trace().Int64("用户id", comm.Ownerid).Msg("[会话]插入群聊信息成功")
	contact := &model.Contact{
		Ownerid:  comm.Ownerid,
		Dstobj:   comm.ID,
		Cate:     model.ContactCateCommunity,
		Memo:     "",
		Createat: comm.CreateAt,
	}
	util.GetLog().Debug().Int64("用户id", comm.Ownerid).Msg("[会话]插入群主关系数据")
	if _, err = session.InsertOne(contact); err != nil {
		util.GetLog().Error().Err(err).Int64("用户id", comm.Ownerid).Msg("[会话]插入群主关系失败")
		if err = session.Rollback(); err != nil {
			util.GetLog().Error().Err(err).Int64("用户id", comm.Ownerid).Msg("[会话]数据库回滚失败")
			return nil, err
		}
		util.GetLog().Trace().Int64("用户id", comm.Ownerid).Msg("[会话]数据库回滚成功")
		return nil, err
	}
	util.GetLog().Trace().Int64("用户id", comm.Ownerid).Msg("[会话]插入群主关系成功")
	if err = session.Commit(); err != nil {
		util.GetLog().Error().Err(err).Int64("用户id", comm.Ownerid).Msg("[会话]提交会话失败")
		return nil, err
	}
	util.GetLog().Trace().Int64("用户id", comm.Ownerid).Msg("[会话]提交会话成功")
	return comm, nil
}
