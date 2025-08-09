package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"friends-of-friends-be/connectors"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
	"net/http"
	"time"
)

type Scores struct {
	Score int
}

type HttpError struct {
	Status  int
	cause   error
	Message string
}

func (h *HttpError) Error() string {
	return "HTTP::ERROR Failed due to " + h.cause.Error()
}

func (h *HttpError) ResponseHeaders() {

}

type RedisSortedResult struct {
	Result []redis.Z
	err    error
}

func GetScores(response http.ResponseWriter, request *http.Request) {
	valueChan := make(chan RedisSortedResult)

	ctx, cancel := context.WithTimeout(request.Context(), 100*time.Second)
	defer cancel()

	go func() {
		value, err := connectors.SortedSetGet(ctx, "user-scores", 0, 20)
		valueChan <- RedisSortedResult{
			Result: value,
			err:    err,
		}
	}()

	select {
	case <-ctx.Done():
		http.Error(response, "Request timed out", http.StatusGatewayTimeout)
	case redisResponse := <-valueChan:
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusOK)
		json.NewEncoder(response).Encode(redisResponse.Result)
	}
}

func UpdateScores(response http.ResponseWriter, request *http.Request) {

}

func GetUserScore(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	user_id, ok, err := validateUserVars(vars)
	if !ok {
		http.Error(response, err.Error(), err.Status)
	}

	result, getErr := connectors.GetRank(request.Context(), "user-scores", user_id)

	if getErr != nil {
		http.Error(response, getErr.Error(), http.StatusInternalServerError)
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(result)
}

func validateUserVars(vars map[string]string) (string, bool, *HttpError) {
	if vars["user_id"] != "" {
		return vars["user_id"], true, nil
	}
	return "", false, &HttpError{
		Status:  http.StatusBadRequest,
		cause:   fmt.Errorf("%s is not present", "user_id"),
		Message: "",
	}
}
