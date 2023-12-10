package service

import (
	"GoChatCraft/common"
	"GoChatCraft/models"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"strconv"
)

func SendUserMsg(ctx *gin.Context) {
	models.Chat(ctx.Writer, ctx.Request)
}

func GetRedisMsg(ctx *gin.Context) {
	getData, _ := ctx.GetRawData()
	var body map[string]string
	_ = json.Unmarshal(getData, &body)
	user := ctx.GetHeader("UserId")
	userId, _ := strconv.Atoi(user)
	targetId, _ := strconv.Atoi(body["targetId"])
	start, _ := strconv.Atoi(body["start"])
	end, _ := strconv.Atoi(body["end"])
	isRev, _ := strconv.ParseBool(body["isRev"])
	res := models.RedisMsg(int64(userId), int64(targetId), int64(start), int64(end), isRev)
	common.RespOkList(ctx.Writer, res, "Get Msg Ok", len(res))
}
