package util

import (
	"encoding/json"
	"log"
	"net/http"
)

// Rsp 这是返回辅助结构体
type Rsp struct {
	Code  int         `json:"code"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data,omitempty"`
	Rows  interface{} `json:"rows,omitempty"`
	Total int         `json:"total,omitempty"`
}

// RespFail 失败时的返回
func RespFail(w http.ResponseWriter, msg string) {
	Resp(w, -1, msg, nil)
}

// RespListFail 失败时的数组返回
func RespListFail(w http.ResponseWriter) {
	RespList(w, -1, nil, 0)
}

// RespOK 正常时的返回
func RespOK(w http.ResponseWriter, msg string, data interface{}) {
	Resp(w, 0, msg, data)
}

// RespListOK 正常时的数组返回
func RespListOK(w http.ResponseWriter, rows interface{}, total int) {
	RespList(w, 0, rows, total)
}

// Resp 通用的返回函数
func Resp(w http.ResponseWriter, code int, msg string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	h := Rsp{
		Code: code,
		Msg:  msg,
		Data: data,
	}
	ret, err := json.Marshal(h)
	if err != nil {
		log.Println(err.Error())
	}
	w.Write(ret)
}

// RespList 通用的数组返回函数
func RespList(w http.ResponseWriter, code int, rows interface{}, total int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	h := Rsp{
		Code:  code,
		Rows:  rows,
		Total: total,
	}
	ret, err := json.Marshal(h)
	if err != nil {
		log.Println(err.Error())
	}
	w.Write(ret)
}
