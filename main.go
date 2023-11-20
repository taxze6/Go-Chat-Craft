package main

import "GoChatCraft/initialize"

func main() {
	//Initialize logging.
	initialize.InitLogger()
	//Initialize database.
	initialize.InitDB()
}
