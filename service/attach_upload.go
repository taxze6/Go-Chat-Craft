package service

import (
	"GoChatCraft/common"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"math/rand"
	"os"
	"time"
)

func File(ctx *gin.Context) {
	w := ctx.Writer
	req := ctx.Request
	srcFile, head, err := req.FormFile("file")
	fileType := req.FormValue("type")
	if err != nil {
		common.RespFail(w, err.Error(), "Unable to read the file.")
		return
	}
	//default folder name
	folderName := "image/"
	switch fileType {
	case "102":
		// processing picture types
		folderName = "image/"
	case "103":
		// processing speech type
		folderName = "voice/"
	case "104":
		// processing video type
		folderName = "video/"
	default:
		// handling unknown types
		folderName = "image/"
	}
	//suffix := ".png"
	//ofName := head.Filename
	//tem := strings.Split(ofName, ".")
	//if len(tem) > 1 {
	//	suffix = "." + tem[len(tem)-1]
	//}
	fileName := fmt.Sprintf("%d%04d%s", time.Now().Unix(), rand.Int31(), head.Filename)

	//fileName := fmt.Sprintf("%s", head.Filename)
	dstFile, err := os.Create("./assets/upload/" + folderName + fileName)
	if err != nil {
		common.RespFail(w, err.Error(), "Failed to create the file.")
		return
	}
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		common.RespFail(w, err.Error(), "Upload failed.")
	}
	url := "http://192.168.31.123:8889/assets/upload/" + folderName + fileName
	common.RespOk(w, url, "Sent successfully.")
}
