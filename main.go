package main

import (
	_ "GoChatCraft/docs"
	"GoChatCraft/global"
	"GoChatCraft/initialize"
	"GoChatCraft/models"
	router "GoChatCraft/router"
	"fmt"
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

	fmt.Println(global.ServiceConfig.Port)
	//Initialize database.
	initialize.InitDB()
	//Initialize redis.
	initialize.InitRedis()
	r := router.Router()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	err := r.Run(fmt.Sprintf(":%d", global.ServiceConfig.Port))
	if err != nil {
		return
	}
}

func init() {
	//Initialize the configuration.
	initialize.InitConfig()

	//UDP
	//go UdpSendProc()
	//go UdpRecProc()

	//rabbitMQ
	models.RabbitmqCreateExchange()
	go models.RabbitmqRecProc()
	go models.RabbitmqSendProc()
}
