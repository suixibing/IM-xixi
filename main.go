package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"text/template"

	"IM-xixi/ctrl"
)

func main() {
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		file, err := os.Open("./favicon.icon")
		if err != nil {
			return
		}
		_, err = io.Copy(w, file)
		if err != nil {
			return
		}
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", "/user/login.shtml")
		w.WriteHeader(301)
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
	http.ListenAndServe(":80", nil)
	// http.ListenAndServeTLS(":8081", "server.crt", "server.key", nil)
}

// RegisterViews 自动注册模板
func RegisterViews() {
	tpls, err := template.New("root").ParseGlob("view/**/*")
	if err != nil {
		log.Fatal(err)
	}
	for _, tpl := range tpls.Templates() {
		tplname := tpl.Name()
		http.HandleFunc(tplname, func(w http.ResponseWriter, r *http.Request) {
			err = tpls.ExecuteTemplate(w, tplname, nil)
			if err != nil {
				log.Fatal(err)
			}
		})
		log.Printf("register tpl:> %s", tplname)
	}
}
