package ctrl

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/suixibing/IM-xixi/util"
)

func init() {
	if err := os.MkdirAll("./mnt", os.ModePerm); err != nil {
		util.GetLog().Error().Err(err).Str("path", "./mnt").Msg("创建目录出错")
		panic("创建目录出错")
	}
}

// Upload 上传控制函数
func Upload(w http.ResponseWriter, r *http.Request) {
	UploadFile(w, r)
}

// UploadFile 上传文件
func UploadFile(w http.ResponseWriter, r *http.Request) {
	util.GetLog().Debug().Msg("获取文件")
	uploadfile, head, err := r.FormFile("file")
	if err != nil {
		util.GetLog().Error().Err(err).Msg("获取文件失败")
		util.RespFail(w, err.Error())
		return
	}
	util.GetLog().Trace().Msg("获取文件成功")

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

	filename = fmt.Sprintf("%d%04d%s", time.Now().Unix(), rand.Int31(), suffix)
	util.GetLog().Debug().Str("文件名", filename).Msg("创建文件")
	dstfile, err := os.Create("./mnt/" + filename)
	if err != nil {
		util.GetLog().Error().Err(err).Str("文件名", filename).Msg("创建文件失败")
		util.RespFail(w, err.Error())
		return
	}
	util.GetLog().Trace().Str("文件名", filename).Msg("创建文件成功")
	util.GetLog().Debug().Str("文件名", filename).Msg("保存文件")
	_, err = io.Copy(dstfile, uploadfile)
	if err != nil {
		util.GetLog().Error().Err(err).Str("文件名", filename).Msg("保存文件失败")
		util.RespFail(w, err.Error())
		return
	}
	util.GetLog().Trace().Str("文件名", filename).Msg("保存文件成功")

	url := "/mnt/" + filename
	util.RespOK(w, "", url)
}
