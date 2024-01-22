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

func FriendList(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.GetHeader("UserId"))
	users, err := dao.FriendList(uint(id))
	if err != nil {
		zap.S().Info("failed to retrieve friend list.", err)
		common.RespFail(ctx.Writer, "Friend list is empty.", "Friend list is empty.")
		return
	}
	infos := make([]models.UserResponse, 0)
	for _, v := range *users {
		info := models.UserResponse{
			ID:         v.ID,
			Name:       v.Name,
			Email:      v.Email,
			Phone:      v.Phone,
			Avatar:     v.Avatar,
			Motto:      v.Motto,
			Identity:   v.Identity,
			ClientIp:   v.ClientIp,
			ClientPort: v.ClientPort,
		}
		infos = append(infos, info)
	}
	common.RespOkList(ctx.Writer, infos, "Successfully retrieved friend list.", len(infos))
}

func AddFriendByName(ctx *gin.Context) {
	getData, _ := ctx.GetRawData()
	var body map[string]string
	_ = json.Unmarshal(getData, &body)
	user := ctx.GetHeader("UserId")
	userId, err := strconv.Atoi(user)
	if err != nil {
		zap.S().Info("failed to convert data type", err)
		return
	}
	tar := body["targetName"]
	code, err := dao.AddFriendByName(uint(userId), tar)
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
	common.RespOk(ctx.Writer, "Successfully added friend.", "Successfully added friend.")
}

func AddFriendByUserId(ctx *gin.Context) {
	getData, _ := ctx.GetRawData()
	var body map[string]int
	_ = json.Unmarshal(getData, &body)
	user := ctx.GetHeader("UserId")
	userId, err := strconv.Atoi(user)
	if err != nil {
		zap.S().Info("failed to convert data type", err)
		return
	}
	tar := body["targetUserId"]
	code, err := dao.AddFriendByUserId(uint(userId), uint(tar))
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
	common.RespOk(ctx.Writer, "Successfully added friend.", "Successfully added friend.")
}
