package main

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/suixibing/IM-xixi/ctrl"
	"github.com/suixibing/IM-xixi/util"
)

func main() {
	conf := util.LoadConfig(util.DefaultConfig, false)
	util.GetLog().Info().Msg("设置服务日志为一号服务的日志")
	util.SetLog(conf.Services[0].Log)
	util.GetLog().Info().Stringer("服务日志级别", util.GetLog().GetLevel()).Msg("服务设置")

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", "/asset/images/favicon.ico")
		w.WriteHeader(301)
		util.Log.Trace().Msg("图标重定向")
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", "/user/login.shtml")
		w.WriteHeader(301)
		util.Log.Trace().Msg("根目录重定向到注册页面")
	})

	http.HandleFunc("/user/login", ctrl.UserLogin)
	http.HandleFunc("/user/logout", ctrl.UserLogout)
	http.HandleFunc("/user/register", ctrl.UserRegister)
	http.HandleFunc("/user/updateinfo", ctrl.UserUpdateInfo)
	http.HandleFunc("/user/find", ctrl.UserFind)
	http.HandleFunc("/user/addfriend", ctrl.AddFriend)
	http.HandleFunc("/user/joincommunity", ctrl.JoinCommunity)

	http.HandleFunc("/contact/loadfriend", ctrl.LoadFriend)
	http.HandleFunc("/contact/createcommunity", ctrl.CreateCommunity)
	http.HandleFunc("/contact/loadcommunity", ctrl.LoadCommunity)

	http.HandleFunc("/chat", ctrl.Chat)
	http.HandleFunc("/attach/upload", ctrl.Upload)
	http.HandleFunc("/loadnotifications", ctrl.LoadNotificationMsg)

	// 设置静态资源目录
	http.Handle("/asset/", http.FileServer(http.Dir(".")))
	http.Handle("/mnt/", http.FileServer(http.Dir(".")))

	RegisterViews()

	port := fmt.Sprintf(":%d", conf.Services[0].Port)
	if conf.Global.UseHttps {
		if err := http.ListenAndServeTLS(port, conf.Global.HttpsCrt, conf.Global.HttpsKey, nil); err != nil {
			util.Log.Error().Err(err).Msg("https服务启动失败")
		}
	} else {
		if err := http.ListenAndServe(port, nil); err != nil {
			util.Log.Error().Err(err).Msg("http服务启动失败")
		}
	}
}

// RegisterViews 自动注册模板
func RegisterViews() {
	tpls, err := template.New("root").ParseGlob("view/**/*")
	if err != nil {
		util.Log.Fatal().Err(err).Msg("ParseGlob失败")
	}
	for _, tpl := range tpls.Templates() {
		tplname := tpl.Name()
		http.HandleFunc(tplname, func(w http.ResponseWriter, r *http.Request) {
			err = tpls.ExecuteTemplate(w, tplname, nil)
			if err != nil {
				util.Log.Fatal().Err(err).Msg("模板执行失败")
			}
		})
		util.Log.Debug().Str("模板名称", tplname).Msg("模板注册")
	}
	util.Log.Info().Msg("views注册结束")
}
