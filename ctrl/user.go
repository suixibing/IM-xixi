package ctrl

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/suixibing/IM-xixi/model"
	"github.com/suixibing/IM-xixi/service"
	"github.com/suixibing/IM-xixi/util"
)

var userService service.UserService
var msgService service.MessageService

// UserLogin 用户登录函数
func UserLogin(w http.ResponseWriter, r *http.Request) {
	util.GetLog().Info().Msg("用户登陆")
	err := r.ParseForm()
	if err != nil {
		util.GetLog().Error().Err(err).Msg("获取用户数据失败")
		util.RespFail(w, "获取用户数据失败")
		return
	}
	util.GetLog().Trace().Msg("获取用户数据成功")

	mobile := r.PostForm.Get("mobile")
	passwd := r.PostForm.Get("passwd")
	util.GetLog().Debug().Str("手机号", mobile).Msg("登陆操作")
	user, err := userService.Login(mobile, passwd)
	if err != nil {
		util.GetLog().Error().Err(err).Str("手机号", mobile).Msg("登陆操作失败")
		util.RespFail(w, fmt.Sprintf("登录失败: %s", err.Error()))
	} else {
		util.GetLog().Trace().Str("手机号", mobile).Msg("登陆操作成功")
		util.RespOK(w, "登陆成功", user)
	}
}

// UserLogout 用户退出函数
func UserLogout(w http.ResponseWriter, r *http.Request) {
	util.GetLog().Info().Msg("用户退出")
	var arg model.User
	util.GetLog().Debug().Msg("参数对象绑定")
	err := util.Bind(r, &arg)
	if err != nil {
		util.GetLog().Error().Err(err).Msg("参数对象绑定失败")
		util.RespFail(w, "参数用户数据异常")
		return
	}
	util.GetLog().Trace().Int64("用户id", arg.ID).Msg("参数对象绑定成功")

	util.GetLog().Debug().Int64("用户id", arg.ID).Msg("退出操作")
	if err := userService.Logout(&arg); err != nil {
		util.GetLog().Error().Err(err).Int64("用户id", arg.ID).Msg("退出操作失败")
		util.RespFail(w, err.Error())
	} else {
		util.GetLog().Trace().Int64("用户id", arg.ID).Msg("退出操作成功")
		util.RespOK(w, "退出成功", nil)
	}
}

// UserRegister 用户注册函数
func UserRegister(w http.ResponseWriter, r *http.Request) {
	util.GetLog().Info().Msg("用户注册")
	err := r.ParseForm()
	if err != nil {
		util.GetLog().Error().Err(err).Msg("获取用户数据失败")
		util.RespFail(w, "获取用户数据失败")
		return
	}
	util.GetLog().Trace().Msg("获取用户数据成功")

	mobile := r.PostForm.Get("mobile")
	passwd := r.PostForm.Get("passwd")
	nickname := fmt.Sprintf("user%s", mobile)
	avatar := "/asset/images/user.jpg"
	sex := model.SexUnknow

	util.GetLog().Debug().Str("手机号", mobile).Msg("注册操作")
	user, err := userService.Register(mobile, passwd, nickname, avatar, sex)
	if err != nil {
		util.GetLog().Error().Err(err).Str("手机号", mobile).Msg("注册操作失败")
		util.RespFail(w, fmt.Sprintf("注册失败: %s", err.Error()))
	} else {
		util.GetLog().Trace().Str("手机号", mobile).Msg("注册操作成功")
		util.RespOK(w, "注册成功", user)
	}
}

// UserFind 用户查找函数
func UserFind(w http.ResponseWriter, r *http.Request) {
	util.GetLog().Info().Msg("用户查找")
	err := r.ParseForm()
	if err != nil {
		util.GetLog().Error().Err(err).Msg("获取用户数据失败")
		util.RespFail(w, "获取用户数据失败")
		return
	}
	util.GetLog().Trace().Msg("获取用户数据成功")

	id := r.PostForm.Get("dstid")
	util.GetLog().Debug().Str("用户id", id).Msg("解析id为数字类型")
	userid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		util.GetLog().Error().Err(err).Str("用户id", id).Msg("使用了不合法的id")
		util.RespFail(w, "使用了不合法的id")
		return
	}
	util.GetLog().Trace().Int64("用户id", userid).Msg("解析id为数字类型成功")

	util.GetLog().Debug().Int64("用户id", userid).Msg("通过id获取用户对象")
	user, err := userService.GetUserByID(userid)
	if err != nil {
		util.GetLog().Error().Err(err).Int64("用户id", userid).Msg("通过id获取用户对象失败")
		util.RespFail(w, "获取用户对象失败")
	} else {
		util.GetLog().Trace().Int64("用户id", userid).Msg("通过id获取用户对象成功")
		util.RespOK(w, "", user)
	}
}

// UserUpdateInfo 更新用户信息
func UserUpdateInfo(w http.ResponseWriter, r *http.Request) {
	util.GetLog().Info().Msg("更新用户信息")
	var arg model.User
	util.GetLog().Debug().Msg("参数对象绑定")
	err := util.Bind(r, &arg)
	if err != nil {
		util.GetLog().Error().Err(err).Msg("参数对象绑定失败")
		util.RespFail(w, "参数用户数据异常")
		return
	}
	util.GetLog().Trace().Int64("用户id", arg.ID).Msg("参数对象绑定成功")

	util.GetLog().Debug().Int64("用户id", arg.ID).Msg("更新用户信息")
	user, err := userService.UpdateInfo(&arg)
	if err != nil {
		util.GetLog().Error().Err(err).Int64("用户id", arg.ID).Msg("更新用户信息失败")
		util.RespFail(w, "更新用户信息失败")
	} else {
		// 失败就暂时先不用更新其他用户的本地缓存
		util.GetLog().Debug().Int64("用户id", arg.ID).Msg("尝试更新其他用户的本用户缓存")
		userids, err := userService.GetRelatedUserids(user.ID)
		if err != nil {
			util.GetLog().Warn().Err(err).Int64("用户id", arg.ID).
				Msg("获取相关用户失败, 放弃更新其他用户的本用户缓存")
			userids = nil
		}
		msg := model.Message{
			Userid:   user.ID,
			Cmd:      model.CmdUpdateUserInfo,
			CreateAt: time.Now(),
		}
		chatService.SendMsgJSONToUsers(userids, msg)
		util.GetLog().Trace().Int64("用户id", arg.ID).Msg("更新用户信息成功")
		util.RespOK(w, "更新用户信息成功", user)
	}
}

// LoadNotificationMsg 加载通知
func LoadNotificationMsg(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	util.GetLog().Debug().Str("通知id", q.Get("dstid")).Msg("解析id为数字类型")
	dstid, err := strconv.ParseInt(q.Get("dstid"), 10, 64)
	if err != nil {
		util.GetLog().Error().Err(err).Str("通知id", q.Get("dstid")).Msg("使用了不合法的id")
		util.RespFail(w, "使用了不合法的id")
		return
	}
	util.GetLog().Trace().Int64("通知id", dstid).Msg("解析id为数字类型成功")

	util.GetLog().Debug().Int64("通知id", dstid).Str("数量", q.Get("num")).Msg("解析目标加载通知条数")
	num, err := strconv.Atoi(q.Get("num"))
	if err != nil {
		util.GetLog().Warn().Err(err).Int64("通知id", dstid).Str("数量", q.Get("num")).
			Msg("解析目标加载通知条数失败, 使用默认值(获取所有)")
		num = 0
	}
	util.GetLog().Trace().Int64("通知id", dstid).Int("数量", num).Msg("解析目标加载通知条数成功")

	util.GetLog().Debug().Int64("通知id", dstid).Int("数量", num).Msg("加载通知")
	msgs, err := msgService.LoadNotificationMsg(dstid, num)
	if err != nil {
		util.GetLog().Error().Err(err).Int64("通知id", dstid).Int("数量", num).Msg("加载通知失败")
		util.RespFail(w, err.Error())
	} else {
		util.GetLog().Trace().Int64("通知id", dstid).Int("数量", num).Msg("加载通知成功")
		util.RespListOK(w, msgs, len(msgs))
	}
}
