package ctrl

import (
	"IM-xixi/model"
	"IM-xixi/service"
	"IM-xixi/util"
	"net/http"
)

var contactService service.ContactService
var communityService service.CommunityService

// AddFriend 添加好友
func AddFriend(w http.ResponseWriter, r *http.Request) {
	var arg model.ContactArg
	// 对象绑定
	util.Bind(r, &arg)

	tmp, exist := userService.GetUserByID(arg.Dstid)
	if !exist {
		tmp, exist = userService.GetUserByMobile(arg.DstMobile)
		if !exist {
			util.RespFail(w, "该用户不存在！")
			return
		}
	}
	err := contactService.AddFriend(arg.Userid, tmp.ID)
	if err != nil {
		util.RespFail(w, err.Error())
	} else {
		util.RespOK(w, "添加好友成功", nil)
	}
}

// JoinCommunity 用户加群
func JoinCommunity(w http.ResponseWriter, r *http.Request) {
	var arg model.ContactArg
	// 对象绑定
	util.Bind(r, &arg)
	err := contactService.JoinCommunity(arg.Userid, arg.Dstid)
	if err != nil {
		util.RespFail(w, err.Error())
	} else {
		err = chatService.AddCommunityID(arg.Userid, arg.Dstid)
		if err != nil {
			util.RespFail(w, err.Error())
			return
		}
		util.RespOK(w, "加群成功", nil)
	}
}

// LoadFriend 加载好友
func LoadFriend(w http.ResponseWriter, r *http.Request) {
	var arg model.ContactArg
	// 对象绑定
	util.Bind(r, &arg)
	users, err := contactService.GetFriends(arg.Userid)
	if err != nil {
		util.RespFail(w, err.Error())
	} else {
		util.RespListOK(w, users, len(users))
	}
}

// CreateCommunity 创建群聊
func CreateCommunity(w http.ResponseWriter, r *http.Request) {
	var arg model.Community
	util.Bind(r, &arg)
	community, err := communityService.CreateCommunity(&arg)
	if err != nil {
		util.RespFail(w, err.Error())
	} else {
		util.RespOK(w, "创建群聊成功", community)
	}
}

// LoadCommunity 加载群聊
func LoadCommunity(w http.ResponseWriter, r *http.Request) {
	var arg model.ContactArg
	// 对象绑定
	util.Bind(r, &arg)
	communitys, err := contactService.GetCommunitys(arg.Userid)
	if err != nil {
		util.RespFail(w, err.Error())
	} else {
		ret := make([]model.CommunityData, len(communitys))
		for i, community := range communitys {
			ret[i].Community = community
			ret[i].Memids, err = communityService.GetAllMembersID(community.ID)
			if err != nil {
				util.RespFail(w, err.Error())
				return
			}
		}
		util.RespListOK(w, ret, len(ret))
	}
}
