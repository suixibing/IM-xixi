package ctrl

import (
	"IM-xixi/util"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

func init() {
	os.MkdirAll("./mnt", os.ModePerm)
}

// Upload 上传控制函数
func Upload(w http.ResponseWriter, r *http.Request) {
	UploadFile(w, r)
}

// UploadFile 上传文件
func UploadFile(w http.ResponseWriter, r *http.Request) {
	uploadfile, head, err := r.FormFile("file")
	if err != nil {
		util.RespFail(w, err.Error())
		return
	}

	suffix := ".jpg"
	filename := head.Filename
	tmp := strings.Split(filename, ".")
	if len(tmp) > 1 {
		suffix = "." + tmp[len(tmp)-1]
	}
	// formdata.append("filetype",".png") 前端可能指定了文件类型
	filetype := r.FormValue("filetype")
	if len(filetype) > 0 {
		suffix = filetype
	}

	filename = fmt.Sprintf("%d%04d%s",
		time.Now().Unix(), rand.Int31(), suffix)
	dstfile, err := os.Create("./mnt/" + filename)
	if err != nil {
		util.RespFail(w, err.Error())
		return
	}
	_, err = io.Copy(dstfile, uploadfile)
	if err != nil {
		util.RespFail(w, err.Error())
		return
	}

	url := "/mnt/" + filename
	util.RespOK(w, "", url)
}
