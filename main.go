package main

import (
	"GoChatCraft/initialize"
	router "GoChatCraft/router"
)

func main() {
	//Initialize logging.
	initialize.InitLogger()
	//Initialize database.
	initialize.InitDB()
	//Initialize redis.
	initialize.InitRedis()
	r := router.Router()
	err := r.Run(":8889")
	if err != nil {
		return
	}
}
