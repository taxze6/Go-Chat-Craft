package router

import (
	"GoChatCraft/middlewear"
	"GoChatCraft/service"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	router := gin.Default()
	v1 := router.Group("v1")
	user := v1.Group("user")
	{
		user.GET("/user_list", middlewear.JWY(), service.GetUserList)
		user.POST("/login", service.LoginByNameAndPassWord)
		user.POST("/register", service.NewUser)
		user.POST("/register_email_code_check", service.CheckRegisterEmailCode)
		user.POST("/email_login", service.EmailLogin)
		user.POST("/email_login_code_check", service.CheckLoginEmailCode)
	}
	return router
}
