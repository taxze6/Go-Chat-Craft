package router

import (
	"GoChatCraft/service"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	router := gin.Default()
	v1 := router.Group("v1")
	user := v1.Group("user")
	{
		user.GET("/user_list", service.GetUserList)
	}
	return router
}
