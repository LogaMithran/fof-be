package controllers

import (
	"encoding/json"
	"net/http"
)

type HealthResponse struct {
	Status  int
	Message string
}

func HealthController(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(&HealthResponse{
		Status:  http.StatusOK,
		Message: "PONG",
	})

	return
}
