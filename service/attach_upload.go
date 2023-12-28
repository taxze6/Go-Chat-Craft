package service

import (
	"GoChatCraft/common"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
)

func Image(ctx *gin.Context) {
	w := ctx.Writer
	req := ctx.Request
	srcFile, head, err := req.FormFile("file")
	if err != nil {
		common.RespFail(w, err.Error(), "Unable to read the file.")
		return
	}

	suffix := ".png"
	ofName := head.Filename
	tem := strings.Split(ofName, ".")
	if len(tem) > 1 {
		suffix = "." + tem[len(tem)-1]
	}
	fileName := fmt.Sprintf("%d%04d%s", time.Now().Unix(), rand.Int31(), suffix)
	dstFile, err := os.Create("./assets/upload/" + fileName)
	if err != nil {
		common.RespFail(w, err.Error(), "Failed to create the file.")
		return
	}
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		common.RespFail(w, err.Error(), "Upload failed.")
	}
	url := "http://192.168.31.135:8889/assets/upload/" + fileName
	common.RespOk(w, url, "Sent successfully.")
}
