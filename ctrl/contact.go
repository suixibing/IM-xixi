package ctrl

import (
	"net/http"

	"github.com/suixibing/IM-xixi/model"
	"github.com/suixibing/IM-xixi/service"
	"github.com/suixibing/IM-xixi/util"
)

var contactService service.ContactService
var communityService service.CommunityService

// AddFriend 添加好友
func AddFriend(w http.ResponseWriter, r *http.Request) {
	var arg model.ContactArg
	util.GetLog().Debug().Msg("参数对象绑定")
	err := util.Bind(r, &arg)
	if err != nil {
		util.GetLog().Error().Err(err).Msg("参数对象绑定失败")
		util.RespFail(w, "参数用户数据异常")
		return
	}
	util.GetLog().Trace().Int64("用户id", arg.Userid).Msg("参数对象绑定成功")

	util.GetLog().Debug().Int64("用户id", arg.Userid).Int64("目标id", arg.Dstid).Msg("通过id获取用户对象")
	tmp, err := userService.GetUserByID(arg.Dstid)
	if err != nil {
		util.GetLog().Warn().Err(err).Int64("用户id", arg.Userid).Int64("目标id", arg.Dstid).Msg("通过id获取用户对象失败")
		util.GetLog().Debug().Int64("用户id", arg.Userid).Str("手机号", arg.DstMobile).Msg("通过手机号获取用户对象")
		tmp, err = userService.GetUserByMobile(arg.DstMobile)
		if err != nil {
			util.GetLog().Error().Err(err).Int64("用户id", arg.Userid).Str("手机号", arg.DstMobile).Msg("通过手机号获取用户对象失败")
			util.RespFail(w, err.Error())
			return
		}
		util.GetLog().Trace().Int64("用户id", arg.Userid).Str("手机号", arg.DstMobile).Msg("通过手机号获取用户对象成功")
	}
	util.GetLog().Trace().Int64("用户id", arg.Userid).Int64("目标id", tmp.ID).Msg("通过id获取用户对象成功")
	util.GetLog().Debug().Int64("用户id", arg.Userid).Int64("目标id", tmp.ID).Msg("添加好友")
	err = contactService.AddFriend(arg.Userid, tmp.ID)
	if err != nil {
		util.GetLog().Error().Err(err).Int64("用户id", arg.Userid).Int64("目标id", tmp.ID).Msg("添加好友失败")
		util.RespFail(w, err.Error())
	} else {
		util.GetLog().Trace().Int64("用户id", arg.Userid).Int64("目标id", tmp.ID).Msg("添加好友成功")
		util.RespOK(w, "添加好友成功", nil)
	}
}

// JoinCommunity 用户加群
func JoinCommunity(w http.ResponseWriter, r *http.Request) {
	var arg model.ContactArg
	util.GetLog().Debug().Msg("参数对象绑定")
	err := util.Bind(r, &arg)
	if err != nil {
		util.GetLog().Error().Err(err).Msg("参数对象绑定失败")
		util.RespFail(w, "参数用户数据异常")
		return
	}
	util.GetLog().Trace().Int64("用户id", arg.Userid).Msg("参数对象绑定成功")

	util.GetLog().Debug().Int64("用户id", arg.Userid).Int64("群聊id", arg.Dstid).Msg("加群")
	err = contactService.JoinCommunity(arg.Userid, arg.Dstid)
	if err != nil {
		util.GetLog().Error().Err(err).Int64("用户id", arg.Userid).Int64("群聊id", arg.Dstid).Msg("加群失败")
		util.RespFail(w, err.Error())
	} else {
		util.GetLog().Trace().Int64("用户id", arg.Userid).Int64("群聊id", arg.Dstid).Msg("加群成功")
		util.GetLog().Debug().Int64("用户id", arg.Userid).Int64("群聊id", arg.Dstid).Msg("为用户节点添加群组id")
		err = chatService.AddCommunityID(arg.Userid, arg.Dstid)
		if err != nil {
			util.GetLog().Warn().Err(err).Int64("用户id", arg.Userid).Int64("群聊id", arg.Dstid).Msg("为用户节点添加群组id失败")
			util.RespOK(w, "加群成功，但是更新数据异常，请重新登陆以便获取最新数据", nil)
			return
		}
		util.GetLog().Trace().Int64("用户id", arg.Userid).Int64("群聊id", arg.Dstid).Msg("为用户节点添加群组id成功")
		util.RespOK(w, "加群成功", nil)
	}
}

// LoadFriend 加载好友
func LoadFriend(w http.ResponseWriter, r *http.Request) {
	var arg model.ContactArg
	util.GetLog().Debug().Msg("参数对象绑定")
	err := util.Bind(r, &arg)
	if err != nil {
		util.GetLog().Error().Err(err).Msg("参数对象绑定失败")
		util.RespFail(w, "参数用户数据异常")
		return
	}
	util.GetLog().Trace().Int64("用户id", arg.Userid).Msg("参数对象绑定成功")

	util.GetLog().Debug().Int64("用户id", arg.Userid).Msg("获取好友列表")
	users, err := contactService.GetFriends(arg.Userid)
	if err != nil {
		util.GetLog().Error().Err(err).Int64("用户id", arg.Userid).Msg("获取好友列表失败")
		util.RespFail(w, err.Error())
	} else {
		util.GetLog().Trace().Int64("用户id", arg.Userid).Msg("获取好友列表成功")
		util.RespListOK(w, users, len(users))
	}
}

// CreateCommunity 创建群聊
func CreateCommunity(w http.ResponseWriter, r *http.Request) {
	var arg model.Community
	util.GetLog().Debug().Msg("参数对象绑定")
	err := util.Bind(r, &arg)
	if err != nil {
		util.GetLog().Error().Err(err).Msg("参数对象绑定失败")
		util.RespFail(w, "参数用户数据异常")
		return
	}
	util.GetLog().Trace().Int64("群聊id", arg.ID).Msg("参数对象绑定成功")

	util.GetLog().Debug().Stringer("群聊id", arg).Msg("创建群聊")
	community, err := communityService.CreateCommunity(&arg)
	if err != nil {
		util.GetLog().Error().Err(err).Msg("创建群聊失败")
		util.RespFail(w, err.Error())
	} else {
		util.GetLog().Trace().Msg("创建群聊成功")
		util.RespOK(w, "创建群聊成功", community)
	}
}

// LoadCommunity 加载群聊
func LoadCommunity(w http.ResponseWriter, r *http.Request) {
	var arg model.ContactArg
	util.GetLog().Debug().Msg("参数对象绑定")
	err := util.Bind(r, &arg)
	if err != nil {
		util.GetLog().Error().Err(err).Msg("参数对象绑定失败")
		util.RespFail(w, "参数用户数据异常")
		return
	}
	util.GetLog().Trace().Int64("用户id", arg.Userid).Msg("参数对象绑定成功")

	util.GetLog().Debug().Int64("用户id", arg.Userid).Msg("获取群聊列表")
	communitys, err := contactService.GetCommunitys(arg.Userid)
	if err != nil {
		util.GetLog().Error().Err(err).Int64("用户id", arg.Userid).Msg("获取群聊列表失败")
		util.RespFail(w, err.Error())
	} else {
		util.GetLog().Trace().Int64("用户id", arg.Userid).Msg("获取群聊列表成功")
		ret := make([]model.CommunityData, len(communitys))
		for i, community := range communitys {
			ret[i].Community = community
			util.GetLog().Debug().Int64("用户id", arg.Userid).Msg("获取群员列表")
			ret[i].Memids, err = communityService.GetAllMembersID(community.ID)
			if err != nil {
				util.GetLog().Warn().Err(err).Int64("用户id", arg.Userid).Int64("群聊id", community.ID).Msg("获取群员列表失败")
				continue
			}
			util.GetLog().Trace().Int64("用户id", arg.Userid).Int64("群聊id", community.ID).Msg("获取群员列表成功")
		}
		util.RespListOK(w, ret, len(ret))
	}
}
