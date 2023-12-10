package main

import (
	_ "GoChatCraft/docs"
	"GoChatCraft/initialize"
	router "GoChatCraft/router"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title			接囗文档
// @version		1.0
// @description	ChatCraft
// @termsofservice	https://github.com/taxze6
// @contact.name	Taxze
// @contact.email	taxze.xiaoyan@gmail.com
// @host			127.0.0.1:8889
func main() {
	//Initialize logging.
	initialize.InitLogger()
	//Initialize database.
	initialize.InitDB()
	//Initialize redis.
	initialize.InitRedis()
	r := router.Router()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	err := r.Run(":8889")
	if err != nil {
		return
	}
}
