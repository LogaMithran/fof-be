package controllers

import (
	"encoding/json"
	"friends-of-friends-be/connectors"
	"github.com/redis/go-redis/v9"
	"net/http"
	"strconv"
	"sync"
)

func GetPeopleWithInXKm(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Access-Control-Allow-Origin", "*")
	q := request.URL.Query()

	kms, _ := strconv.ParseFloat(q.Get("km"), 64)
	lat, _ := strconv.ParseFloat(q.Get("lat"), 64)
	lng, _ := strconv.ParseFloat(q.Get("lng"), 64)

	location := connectors.GeoSearch("user-location", lat, lng, kms)

	filteredLocation := CheckIfActive(location)
	response.WriteHeader(http.StatusOK)
	response.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(response).Encode(filteredLocation); err != nil {
		http.Error(response, "Error in writing the response", http.StatusInternalServerError)
	}
}

func seedIn(location []redis.GeoLocation) <-chan redis.GeoLocation {
	seedChan := make(chan redis.GeoLocation)

	go func() {
		defer close(seedChan)

		for _, geoLocation := range location {
			seedChan <- geoLocation
		}
	}()

	return seedChan
}

func CheckIfActive(location []redis.GeoLocation) []redis.GeoLocation {
	var filteredLocation []redis.GeoLocation

	done := make(chan interface{})
	defer close(done)

	seedChan := seedIn(location)
	resChan := make([]<-chan redis.GeoLocation, len(location))

	for i := 0; i < len(location); i++ {
		resChan[i] = checkIfUserActiveWithRedis(done, seedChan)
	}

	multiplexedStread := FanIn(done, resChan...)

	for value := range multiplexedStread {

		if value.Name != "" {
			filteredLocation = append(filteredLocation, value)
		}
	}

	return filteredLocation
}

func checkIfUserActiveWithRedis(done <-chan interface{}, seedChan <-chan redis.GeoLocation) <-chan redis.GeoLocation {
	statusChan := make(chan redis.GeoLocation)

	go func() {
		defer close(statusChan)
		for val := range seedChan {
			response := connectors.Get(val.Name)

			select {
			case <-done:
				return
			default:
				if response != "" {
					statusChan <- val
				}
			}
		}
	}()

	return statusChan
}

func FanIn(done <-chan interface{}, chanList ...<-chan redis.GeoLocation) <-chan redis.GeoLocation {
	var wg sync.WaitGroup

	multiplexedStream := make(chan redis.GeoLocation)

	multiplex := func(c <-chan redis.GeoLocation) {
		defer wg.Done()

		select {
		case <-done:
			return
		case multiplexedStream <- <-c:

		}
	}

	wg.Add(len(chanList))
	for _, c := range chanList {
		go multiplex(c)
	}

	go func() {
		wg.Wait()
		close(multiplexedStream)
	}()

	return multiplexedStream
}

func GetFriends(response http.ResponseWriter, request *http.Request) {
	//vars := mux.Vars(request)
	//var result []entities.User
	//
	//connectors.Db.Preload()
}
