package service

import (
	"GoChatCraft/dao"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetUserList(ctx *gin.Context) {
	list, err := dao.GetUserList()
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    -1,
			"message": "Failed to retrieve the user list.",
		})
		return
	}
	ctx.JSON(http.StatusOK, list)
}
