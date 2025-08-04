package controllers

import (
	"encoding/json"
	"fmt"
	"friends-of-friends-be/connectors"
	"friends-of-friends-be/services"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

type WsMessage struct {
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
	UserName string  `json:"userName"`
}

type wsConnection struct {
	userName     string
	wsConnection *websocket.Conn
}

var (
	activeConnections = make(map[string]*wsConnection)
	mutex             = &sync.Mutex{}
)

func WsController(response http.ResponseWriter, request *http.Request) {
	var wsMessage WsMessage
	connection, err := WsUpgrade.Upgrade(response, request, nil)

	if err != nil {
		fmt.Println("Error in upgrading the connection")
		return
	}

	defer func() {
		connection.Close()
		delete(activeConnections, wsMessage.UserName)
	}()

	for {
		_, message, readErr := connection.ReadMessage()

		if readErr != nil {
			fmt.Println("Error in reading the message")
			return
		}

		if parseErr := json.Unmarshal(message, &wsMessage); parseErr != nil {
			fmt.Println("Error in parsing the message", parseErr)
			return
		}

		mutex.Lock()
		wsConn := &wsConnection{
			userName:     wsMessage.UserName,
			wsConnection: connection,
		}
		activeConnections[wsConn.userName] = wsConn
		mutex.Unlock()

		go AddUserToGeo(wsMessage)
	}
}

func AddUserToGeo(wsMessage WsMessage) {
	connectors.GeolocationAdd("user-location", services.Location{
		Longitude: wsMessage.Lng,
		Latitude:  wsMessage.Lat,
		TimeStamp: time.Now().Unix(),
	}, wsMessage.UserName)

	value, _ := json.Marshal(wsMessage)
	connectors.Set(wsMessage.UserName, string(value), 24*time.Hour)

	go getSurroundingUsers(wsMessage, activeConnections)
}

func getSurroundingUsers(wsMessage WsMessage, activeConnections map[string]*wsConnection) {
	defer func() {
		if r := recover(); r != nil {
			println("Error in writing back to sockets")
		}
	}()
	location := connectors.GeoSearch("user-location", wsMessage.Lat, wsMessage.Lng, 100)
	location = CheckIfActive(location)

	for _, geoLocation := range location {
		message := fmt.Sprintf("Spotted %s", wsMessage.UserName)
		if activeConnections[geoLocation.Name].userName != "" && wsMessage.UserName != activeConnections[geoLocation.Name].userName {
			if err := activeConnections[geoLocation.Name].wsConnection.WriteMessage(1, []byte(message)); err != nil {
				println("Error in sending the message to client", err)
			}
		}
	}
}
