package service

import (
	"GoChatCraft/common"
	"GoChatCraft/dao"
	"GoChatCraft/models"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

func NewGroup(ctx *gin.Context) {
	getData, _ := ctx.GetRawData()
	var body map[string]string
	_ = json.Unmarshal(getData, &body)
	img := body["group_img"]
	name := body["group_name"]
	desc := body["group_desc"]
	user := ctx.GetHeader("UserId")
	userId, err := strconv.Atoi(user)
	if err != nil {
		zap.S().Info("failed to convert data type", err)
		return
	}
	community := models.Community{}

	community.Name = name
	community.Type = 2
	community.Image = img
	community.Desc = desc
	community.OwnerId = uint(userId)
	code, err := dao.CreateCommunity(community)
	if err != nil {
		var res string
		switch code {
		case -1:
			res = err.Error()
		case 0:
			res = "The friend already exists."
		case -2:
			res = "You cannot add yourself."
		default:
			res = "Unknown error."
		}
		common.RespFail(ctx.Writer, res, res)
		return
	}

	ctx.JSON(200, gin.H{
		"code":    0, //  0成功   -1失败
		"message": "键群成功",
	})
}
