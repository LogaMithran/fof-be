package controllers

import (
	"encoding/json"
	"friends-of-friends-be/connectors"
	"net/http"
	"strconv"
)

func GetPeopleWithInXKm(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Access-Control-Allow-Origin", "*")
	q := request.URL.Query()

	kms, _ := strconv.ParseFloat(q.Get("km"), 64)
	lat, _ := strconv.ParseFloat(q.Get("lat"), 64)
	lng, _ := strconv.ParseFloat(q.Get("lng"), 64)

	location := connectors.GeoSearch("user-location", lat, lng, kms)

	response.WriteHeader(http.StatusOK)
	response.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(response).Encode(location); err != nil {
		http.Error(response, "Error in writing the response", http.StatusInternalServerError)
	}
}
