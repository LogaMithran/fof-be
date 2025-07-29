package controllers

import (
	"github.com/gorilla/websocket"
	"net/http"
)

var WsUpgrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
