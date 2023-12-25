package service

import (
	"GoChatCraft/common"
	"GoChatCraft/models"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

func SendUserMsg(ctx *gin.Context) {
	models.Chat(ctx.Writer, ctx.Request, ctx.GetHeader("UserId"))
}

func GetRedisMsg(ctx *gin.Context) {
	getData, _ := ctx.GetRawData()
	var body map[string]interface{}
	_ = json.Unmarshal(getData, &body)
	test := body["targetId"]
	fmt.Println(test)
	user := ctx.GetHeader("UserId")
	userId, _ := strconv.Atoi(user)
	targetId := 0
	if val, ok := body["targetId"].(float64); ok {
		targetId = int(val)
	}

	start := 0
	if val, ok := body["start"].(float64); ok {
		start = int(val)
	}

	end := 0
	if val, ok := body["end"].(float64); ok {
		end = int(val)
	}

	isRev := false
	if val, ok := body["isRev"].(bool); ok {
		isRev = val
	}
	res := models.RedisMsg(int64(userId), int64(targetId), int64(start), int64(end), isRev)
	fmt.Println(res)
	common.RespOkList(ctx.Writer, res, "Get Msg Ok", len(res))
}
