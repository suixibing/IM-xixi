package ctrl

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"IM-xixi/model"
	"IM-xixi/service"
	"IM-xixi/util"
)

var userService service.UserService
var msgService service.MessageService

// UserLogin 用户登录函数
func UserLogin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	mobile := r.PostForm.Get("mobile")
	passwd := r.PostForm.Get("passwd")
	user, err := userService.Login(mobile, passwd)
	if err != nil {
		util.RespFail(w, fmt.Sprintf("登录失败: %s", err.Error()))
	} else {
		util.RespOK(w, "", user)
	}
}

// UserLogout 用户退出函数
func UserLogout(w http.ResponseWriter, r *http.Request) {
	var arg model.User
	util.Bind(r, &arg)
	if err := userService.Logout(&arg); err != nil {
		util.RespFail(w, err.Error())
	} else {
		util.RespOK(w, "退出成功", nil)
	}
}

// UserRegister 用户注册函数
func UserRegister(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	mobile := r.PostForm.Get("mobile")
	passwd := r.PostForm.Get("passwd")
	nickname := fmt.Sprintf("user%s", mobile)
	avatar := "/asset/images/user.jpg"
	sex := model.SexUnknow

	user, err := userService.Register(mobile, passwd, nickname, avatar, sex)
	if err != nil {
		util.RespFail(w, fmt.Sprintf("注册失败: %s", err.Error()))
	} else {
		util.RespOK(w, "", user)
	}
}

// UserFind 用户查找函数
func UserFind(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	id := r.PostForm.Get("dstid")
	userid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		util.RespFail(w, fmt.Sprintf("查找用户失败: %s", err.Error()))
	}
	user, exist := userService.GetUserByID(userid)
	if !exist {
		util.RespFail(w, "用户不存在！")
	} else {
		util.RespOK(w, "", user)
	}
}

// UserUpdateInfo 更新用户信息
func UserUpdateInfo(w http.ResponseWriter, r *http.Request) {
	var arg model.User
	util.Bind(r, &arg)
	user, err := userService.UpdateInfo(&arg)
	if err != nil {
		util.RespFail(w, err.Error())
	} else {
		// 失败就暂时先不用更新其他用户的本地缓存
		userids, _ := userService.GetRelatedUserids(user.ID)
		msg := model.Message{
			Userid:   user.ID,
			Cmd:      model.CmdUpdateUserInfo,
			CreateAt: time.Now(),
		}
		chatService.SendMsgJSONToUsers(userids, msg)
		util.RespOK(w, "更新用户信息成功", user)
	}
}

// LoadNotificationMsg 加载通知
func LoadNotificationMsg(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	dstid, err := strconv.ParseInt(q.Get("dstid"), 10, 64)
	if err != nil {
		util.RespFail(w, err.Error())
		return
	}
	num, err := strconv.Atoi(q.Get("num"))
	if err != nil {
		util.RespFail(w, err.Error())
		return
	}
	msgs, err := msgService.LoadNotificationMsg(dstid, num)
	if err != nil {
		util.RespFail(w, err.Error())
	} else {
		util.RespListOK(w, msgs, len(msgs))
	}
}
