package controllers

import (
	"encoding/json"
	"fmt"
	"friends-of-friends-be/connectors"
	"friends-of-friends-be/services"
	"net/http"
	"time"
)

type WsMessage struct {
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
	UserName string  `json:"userName"`
}

func WsController(response http.ResponseWriter, request *http.Request) {
	var wsMessage WsMessage
	connection, err := WsUpgrade.Upgrade(response, request, nil)

	if err != nil {
		fmt.Println("Error in upgrading the connection")
		return
	}

	defer connection.Close()

	for {
		_, message, readErr := connection.ReadMessage()

		if readErr != nil {
			fmt.Println("Error in reading the message")
		}

		if parseErr := json.Unmarshal(message, &wsMessage); parseErr != nil {
			fmt.Println("Error in parsing the message", parseErr)
			return
		}

		connectors.GeolocationAdd("user-location", services.Location{
			Longitude: wsMessage.Lng,
			Latitude:  wsMessage.Lat,
			TimeStamp: time.Now().Unix(),
		}, wsMessage.UserName)
	}
}
