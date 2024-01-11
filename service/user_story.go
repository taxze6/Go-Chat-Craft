package service

import (
	"GoChatCraft/common"
	"GoChatCraft/dao"
	"GoChatCraft/models"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetStoryList(ctx *gin.Context) {
	getData, _ := ctx.GetRawData()
	var body map[string]string
	_ = json.Unmarshal(getData, &body)
	userId, _ := strconv.Atoi(body["userId"])
	page, _ := strconv.Atoi(body["page"])
	pageSize, _ := strconv.Atoi(body["pageSize"])

	storyList, err := dao.GetStoryList(uint(userId), page, pageSize)
	if err != nil {
		common.RespFail(ctx.Writer, "couldn't find any information about this user story", "couldn't find any information about this user story")
		return
	}
	common.RespOkList(ctx.Writer, storyList, "The user story has been found.", len(*storyList))
}

func GetUserShowStoryList(ctx *gin.Context) {
	getData, _ := ctx.GetRawData()
	var body map[string]int
	_ = json.Unmarshal(getData, &body)
	userId := body["userId"]
	storyList, count, err := dao.GetUserShowStoryList(uint(userId))
	if err != nil {
		common.RespFail(ctx.Writer, "couldn't find any information about this user story", "couldn't find any information about this user story")
		return
	}
	response := models.UserShowStoryListResponse{
		StoryList: storyList,
		Count:     count,
	}
	common.RespOk(ctx.Writer, response, "Successfully obtained user Story data.")
}

func AddStory(ctx *gin.Context) {
	getData, _ := ctx.GetRawData()
	var body map[string]string
	_ = json.Unmarshal(getData, &body)
	userId, _ := strconv.Atoi(body["userId"])
	content := body["content"]
	media := body["media"]
	storyType, _ := strconv.Atoi(body["type"])
	story := &models.UserStory{
		OwnerId: uint(userId),
		Content: content,
		Media:   media,
		Type:    storyType,
	}
	storyRes, err := dao.AddStory(story)
	if err != nil {
		common.RespFail(ctx.Writer, "failed to add the Story!", "failed to add the Story!")
		return
	}
	common.RespOk(ctx.Writer, storyRes, "successfully added the Story.")
}

func AddStoryLike(ctx *gin.Context) {
	getData, _ := ctx.GetRawData()
	var body map[string]int
	_ = json.Unmarshal(getData, &body)
	storyId, _ := body["storyId"]
	likeOwnerId, _ := body["ownerId"]
	userStoryLike := &models.UserStoryLike{
		UserStoryId: uint(storyId),
		LikeOwnerId: uint(likeOwnerId),
	}
	err := dao.AddOrRemoveStoryLike(userStoryLike)
	if err != nil {
		common.RespFail(ctx.Writer, "Failed to like the post!", "Failed to like the post!")
		return
	}
	common.RespOk(ctx.Writer, "Successfully!", "Successfully liked the post!")
}

func AddStoryComment(ctx *gin.Context) {
	getData, _ := ctx.GetRawData()
	var body map[string]string
	_ = json.Unmarshal(getData, &body)
	storyId, _ := strconv.Atoi(body["storyId"])
	commentOwnerId, _ := strconv.Atoi(body["ownerId"])
	commentContent := body["content"]
	commentType, _ := strconv.Atoi(body["type"])
	userStoryComment := &models.UserStoryComment{
		UserStoryId:    uint(storyId),
		CommentOwnerId: uint(commentOwnerId),
		CommentContent: commentContent,
		Type:           commentType,
	}
	err := dao.AddStoryComment(userStoryComment)
	if err != nil {
		common.RespFail(ctx.Writer, "Failed to comment the post!", "Failed to comment the post!")
		return
	}
	common.RespOk(ctx.Writer, "Successfully comment the post!", "Successfully comment the post!")
}
