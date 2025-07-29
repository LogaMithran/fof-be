package main

import (
	"friends-of-friends-be/connectors"
	"friends-of-friends-be/routers"
	"log"
)

func main() {
	if _, err := connectors.ConnectToRedis(); err != nil {
		log.Fatalf(err.Error())
	}
	connectors.InitializeDbConnection()

	routers.Router()
}
