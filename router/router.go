package router

import (
	"GoChatCraft/middlewear"
	"GoChatCraft/service"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	router := gin.Default()
	v1 := router.Group("v1")
	//assets
	router.Static("/assets", "assets/")

	user := v1.Group("user")
	{
		user.GET("/user_list", middlewear.JWY(), service.GetUserList)
		user.POST("/login", service.LoginByNameAndPassWord)
		user.POST("/register", service.NewUser)
		user.POST("/register_email_code_check", service.CheckRegisterEmailCode)
		user.POST("/email_login", service.EmailLogin)
		user.POST("/email_login_code_check", service.CheckLoginEmailCode)
		user.POST("/find_user_with_name", middlewear.JWY(), service.FindUserWithUserName)
	}
	relation := v1.Group("relation").Use(middlewear.JWY())
	{
		relation.POST("/list", service.FriendList)
		relation.POST("/add_username", service.AddFriendByName)
		relation.POST("/add_userid", service.AddFriendByUserId)
	}
	group := v1.Group("group").Use(middlewear.JWY())
	{
		group.POST("/new_group", service.NewGroup)
	}
	message := v1.Group("message").Use(middlewear.JWY())
	{
		message.GET("/send_user_msg", service.SendUserMsg)
		message.POST("/get_redis_msg", service.GetRedisMsg)
	}

	upload := v1.Group("upload").Use(middlewear.JWY())
	{
		upload.POST("/image", service.Image)
	}
	return router
}
