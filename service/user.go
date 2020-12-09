package service

import (
	"IM-xixi/model"
	"IM-xixi/util"
	"fmt"
	"math/rand"
	"time"

	"github.com/prometheus/common/log"
)

// UserService 用于使用用户服务
type UserService struct{}

var userService UserService

// GetUserByMobile 检查用户是否存在
func (us UserService) GetUserByMobile(mobile string) (*model.User, bool) {
	user := &model.User{}
	_, err := DBEngine.Where("mobile=?", mobile).Get(user)
	if err != nil || user.ID <= 0 {
		return nil, false
	}
	return user, true
}

// GetUserByID 通过id获取用户对象
func (us UserService) GetUserByID(userid int64) (*model.User, bool) {
	user := &model.User{}
	_, err := DBEngine.Where("id=?", userid).Get(user)
	if err != nil || user.ID <= 0 {
		return nil, false
	}
	return user, true
}

// GetUsers 获得所有成员
func (us UserService) GetUsers(uids []int64) (
	[]*model.User, error) {
	users := make([]*model.User, 0)
	for _, id := range uids {
		user, ok := us.GetUserByID(id)
		if !ok {
			log.Error(id, ": 用户不存在")
			return users, nil
		}
		users = append(users, user)
	}
	return users, nil
}

// Login 登录服务
func (us UserService) Login(mobile, plainpwd string) (*model.User, error) {
	user, exist := us.GetUserByMobile(mobile)
	if !exist {
		return nil, fmt.Errorf("该手机号未注册")
	}

	if !util.ValidatePasswd(plainpwd, user.Salt, user.Passwd) {
		return nil, fmt.Errorf("密码错误")
	}
	// 更新登录状态 刷新token
	user.Online = model.IsOnline
	user.Token = fmt.Sprintf("%d", time.Now().Unix())
	_, err := DBEngine.ID(user.ID).Cols("online", "token").Update(user)
	return user, err
}

// Logout 退出登录服务
func (us UserService) Logout(user *model.User) error {
	user, exist := us.GetUserByID(user.ID)
	if !exist {
		return fmt.Errorf("该手机号未注册")
	}
	// 更新登录状态 刷新token
	user.Online = model.NotOnline
	_, err := DBEngine.ID(user.ID).Cols("online").Update(user)
	return err
}

// Register 注册服务
func (us UserService) Register(mobile, plainpwd, nickname,
	avatar, sex string) (*model.User, error) {
	if _, exist := us.GetUserByMobile(mobile); exist {
		return nil, fmt.Errorf("该手机号已经注册")
	}
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

	_, err := DBEngine.InsertOne(user)
	return user, err
}

// UpdateInfo 更新用户信息
func (us UserService) UpdateInfo(user *model.User) (*model.User, error) {
	if user.ID == 0 {
		return nil, fmt.Errorf("请先登录！")
	}

	tmp, exist := us.GetUserByID(user.ID)
	if !exist {
		return nil, fmt.Errorf("该手机号未注册")
	}
	tmp.Avatar = user.Avatar
	tmp.Nickname = user.Nickname
	tmp.Sex = user.Sex
	tmp.Memo = user.Memo
	_, err := DBEngine.ID(user.ID).Update(tmp)
	return tmp, err
}

// CheckToken 检查token是否有效
func (us UserService) CheckToken(userid int64, token string) bool {
	user, _ := us.GetUserByID(userid)
	return user.Token == token
}

// GetRelatedUserids 获取相关的所有用户的id
func (us UserService) GetRelatedUserids(userid int64) ([]int64, error) {
	contactService := ContactService{}
	userids, err := contactService.GetFriendsID(userid)
	if err != nil {
		return userids, err
	}
	communityids, err := contactService.GetCommunitysID(userid)
	if err != nil {
		return userids, err
	}

	communityService := CommunityService{}
	for _, communityid := range communityids {
		retids, err := communityService.GetAllMembersID(communityid)
		if err != nil {
			return userids, err
		}
		userids = append(userids, retids...)
	}
	userids = util.RemoveRep(userids)
	return userids, nil
}
