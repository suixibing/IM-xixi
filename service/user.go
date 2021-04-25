package service

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/suixibing/IM-xixi/model"
	"github.com/suixibing/IM-xixi/util"
)

// UserService 用于使用用户服务
type UserService struct{}

var userService UserService

// GetUserByMobile 检查用户是否存在
func (us UserService) GetUserByMobile(mobile string) (*model.User, error) {
	user := &model.User{}
	util.GetLog().Debug().Str("手机号", mobile).Msg("通过手机号查找用户")
	_, err := DBEngine.Where("mobile=?", mobile).Get(user)
	if err != nil || user.ID <= 0 {
		util.GetLog().Warn().Str("手机号", mobile).Msg("手机号未注册")
		return nil, fmt.Errorf("手机号未注册")
	}
	util.GetLog().Trace().Str("手机号", mobile).Msg("获取用户成功")
	return user, nil
}

// GetUserByID 通过id获取用户对象
func (us UserService) GetUserByID(userid int64) (*model.User, error) {
	user := &model.User{}
	util.GetLog().Debug().Int64("用户id", userid).Msg("通过id查找用户")
	_, err := DBEngine.Where("id=?", userid).Get(user)
	if err != nil || user.ID <= 0 {
		util.GetLog().Warn().Int64("用户id", userid).Msg("用户不存在")
		return nil, fmt.Errorf("用户不存在")
	}
	util.GetLog().Trace().Int64("用户id", userid).Msg("获取用户成功")
	return user, nil
}

// GetUsers 获得所有成员
func (us UserService) GetUsers(uids []int64) (
	[]*model.User, error) {
	users := make([]*model.User, 0)
	util.GetLog().Debug().Ints64("用户id", uids).Msg("获取用户组")
	for _, id := range uids {
		user, err := us.GetUserByID(id)
		if err != nil {
			util.GetLog().Warn().Err(err).Int64("用户id", id).Msg("获取用户失败")
			continue
		}
		users = append(users, user)
	}
	util.GetLog().Trace().Ints64("用户id", uids).Interface("用户组", users).Msg("获取用户组结果")
	return users, nil
}

// Login 登录服务
func (us UserService) Login(mobile, plainpwd string) (*model.User, error) {
	util.GetLog().Debug().Str("手机号", mobile).Msg("用户登陆")
	user, err := us.GetUserByMobile(mobile)
	if err != nil {
		util.GetLog().Error().Err(err).Int64("用户id", user.ID).Msg("登陆失败")
		return nil, fmt.Errorf("该手机号未注册")
	}

	if !util.ValidatePasswd(plainpwd, user.Salt, user.Passwd) {
		util.GetLog().Error().Str("error", "密码错误").Int64("用户id", user.ID).Stringer("用户", user).Msg("登陆失败")
		return nil, fmt.Errorf("密码错误")
	}
	// 更新登录状态 刷新token
	user.Online = model.IsOnline
	user.Token = fmt.Sprintf("%d", time.Now().Unix())
	util.GetLog().Trace().Int64("用户id", user.ID).Msg("更新登陆信息")
	err = us.updateUserInfo(user, []string{"online", "token"})
	if err != nil {
		util.GetLog().Error().Err(err).Int64("用户id", user.ID).Msg("更新登陆信息失败")
		return nil, err
	}
	util.GetLog().Trace().Int64("用户id", user.ID).Msg("更新登陆信息成功")
	return user, nil
}

// Logout 退出登录服务
func (us UserService) Logout(user *model.User) error {
	util.GetLog().Debug().Int64("用户id", user.ID).Stringer("用户", user).Msg("退出登陆")
	user, err := us.GetUserByID(user.ID)
	if err != nil {
		util.GetLog().Error().Err(err).Int64("用户id", user.ID).Msg("退出登陆失败")
		return err
	}
	// 更新登录状态 刷新token
	user.Online = model.NotOnline
	user.Token = fmt.Sprintf("%d", time.Now().Unix())
	util.GetLog().Trace().Int64("用户id", user.ID).Msg("更新登陆信息")
	err = us.updateUserInfo(user, []string{"online", "token"})
	if err != nil {
		util.GetLog().Error().Err(err).Int64("用户id", user.ID).Msg("更新登陆信息失败")
		return err
	}
	util.GetLog().Trace().Int64("用户id", user.ID).Msg("退出登陆成功")
	return nil
}

// Register 注册服务
func (us UserService) Register(mobile, plainpwd, nickname,
	avatar, sex string) (*model.User, error) {
	util.GetLog().Debug().Str("手机号", mobile).Msg("注册操作")
	if _, err := us.GetUserByMobile(mobile); err != nil {
		util.GetLog().Warn().Err(err).Msg("注册失败")
		return nil, fmt.Errorf("该手机号已经注册")
	}
	util.GetLog().Trace().Str("手机号", mobile).Msg("该手机号可以注册")
	user := &model.User{
		Mobile:   mobile,
		Nickname: nickname,
		Avatar:   avatar,
		Sex:      sex,
		Salt:     fmt.Sprintf("%6d", rand.Int31n(10000)),
		CreateAt: time.Now(),
		Token:    fmt.Sprintf("%8d", rand.Int31()),
	}
	user.Passwd = util.CreatePasswd(plainpwd, user.Salt)

	util.GetLog().Debug().Str("手机号", mobile).Stringer("用户信息", user).Msg("保存用户")
	_, err := DBEngine.InsertOne(user)
	if err != nil {
		util.GetLog().Error().Err(err).Str("手机号", mobile).Stringer("用户信息", user).Msg("保存用户失败")
		return nil, err
	}
	util.GetLog().Trace().Str("手机号", mobile).Stringer("用户信息", user).Msg("保存用户成功")
	return user, err
}

// UpdateInfo 更新用户信息
func (us UserService) UpdateInfo(user *model.User) (*model.User, error) {
	if user.ID == 0 {
		util.GetLog().Warn().Int64("用户id", user.ID).Msg("用户未登陆")
		return nil, fmt.Errorf("请先登录！")
	}

	util.GetLog().Debug().Int64("用户id", user.ID).Msg("获取用户原有数据")
	tmp, err := us.GetUserByID(user.ID)
	if err != nil {
		util.GetLog().Error().Err(err).Int64("用户id", user.ID).Msg("获取用户原有数据失败")
		return nil, err
	}
	util.GetLog().Trace().Int64("用户id", user.ID).Msg("获取用户原有数据成功")
	tmp.Avatar = user.Avatar
	tmp.Nickname = user.Nickname
	tmp.Sex = user.Sex
	tmp.Memo = user.Memo
	util.GetLog().Debug().Int64("用户id", user.ID).Stringer("用户", tmp).Msg("更新用户信息")
	_, err = DBEngine.ID(user.ID).Update(tmp)
	if err != nil {
		util.GetLog().Error().Err(err).Int64("用户id", user.ID).Stringer("用户", tmp).Msg("更新用户信息失败")
		return nil, err
	}
	util.GetLog().Trace().Int64("用户id", user.ID).Stringer("用户", tmp).Msg("更新用户信息成功")
	return tmp, nil
}

// updateUserInfo 更新用户在线信息
func (us UserService) updateUserInfo(user *model.User, cols []string) error {
	util.GetLog().Debug().Int64("用户id", user.ID).Stringer("用户", user).Msg("更新用户信息")
	_, err := DBEngine.ID(user.ID).Cols(cols...).Update(user)
	if err != nil {
		util.GetLog().Error().Err(err).Int64("用户id", user.ID).Stringer("用户", user).Msg("更新用户信息失败")
		return fmt.Errorf("更新用户信息失败")
	}
	util.GetLog().Trace().Int64("用户id", user.ID).Stringer("用户", user).Msg("更新用户信息成功")
	return nil

}

// CheckToken 检查token是否有效
func (us UserService) CheckToken(userid int64, token string) bool {
	util.GetLog().Trace().Int64("用户id", userid).Str("token", token).Msg("校验token")
	user, _ := us.GetUserByID(userid)
	return user.Token == token
}

// GetRelatedUserids 获取相关的所有用户的id
func (us UserService) GetRelatedUserids(userid int64) ([]int64, error) {
	contactService := ContactService{}
	util.GetLog().Debug().Int64("用户id", userid).Msg("获取用户的所有好友id")
	userids, err := contactService.GetFriendsID(userid)
	if err != nil {
		util.GetLog().Error().Err(err).Int64("用户id", userid).Msg("获取用户的所有好友id失败")
		return userids, err
	}
	util.GetLog().Trace().Int64("用户id", userid).Msg("获取用户的所有好友id成功")
	util.GetLog().Debug().Int64("用户id", userid).Msg("获取用户的所有群聊id")
	communityids, err := contactService.GetCommunitysID(userid)
	if err != nil {
		util.GetLog().Error().Err(err).Int64("用户id", userid).Msg("获取用户的所有群聊id失败")
		return userids, err
	}
	util.GetLog().Trace().Int64("用户id", userid).Msg("获取用户的所有群聊id成功")

	communityService := CommunityService{}
	for _, communityid := range communityids {
		util.GetLog().Debug().Int64("群聊id", communityid).Msg("获取群聊内的所有用户id")
		retids, err := communityService.GetAllMembersID(communityid)
		if err != nil {
			util.GetLog().Error().Err(err).Int64("群聊id", communityid).Msg("获取群聊内的所有用户id失败")
			return userids, err
		}
		util.GetLog().Trace().Int64("群聊id", communityid).Msg("获取群聊内的所有用户id成功")
		userids = append(userids, retids...)
	}
	userids = util.RemoveRep(userids)
	return userids, nil
}
