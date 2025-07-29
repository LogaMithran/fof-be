package controllers

import (
	"encoding/json"
	"friends-of-friends-be/connectors"
	"friends-of-friends-be/entities"
	"net/http"
	"strconv"
)

func GetUsersController(response http.ResponseWriter, request *http.Request) {

	var users []entities.User

	connectors.Db.Scopes(connectors.Paginate(request)).Find(&users)

	response.WriteHeader(http.StatusOK)
	response.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(response).Encode(users); err != nil {
		http.Error(response, "Error in encoding the json", http.StatusInternalServerError)
	}
}

func AddFriendController(response http.ResponseWriter, request *http.Request) {
	q := request.URL.Query()

	userId, _ := strconv.ParseUint(q.Get("user_id"), 10, 64)
	friendId, _ := strconv.ParseUint(q.Get("friend_id"), 10, 64)

	connectors.Db.Create(entities.UserFriend{
		UserId:   userId,
		FriendId: friendId,
	})

	response.WriteHeader(http.StatusNoContent)
}

func CreateUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Access-Control-Allow-Headers", "content-type")
	response.Header().Set("Access-Control-Allow-Origin", "*")
	response.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

	if request.Method == http.MethodOptions {
		return
	}

	var user entities.User

	defer func() {
		if r := recover(); r != nil {
			println(r)
		}
	}()
	if err := json.NewDecoder(request.Body).Decode(&user); err != nil {
		http.Error(response, "Error in decoding the body", http.StatusInternalServerError)
	}

	if user.Email == "" {
		http.Error(response, "Email is necessary", http.StatusBadRequest)
		return
	}

	if err := connectors.Db.Create(&user).Error; err != nil {
		http.Error(response, "Error in creating the record", http.StatusInternalServerError)
	}

	response.WriteHeader(http.StatusCreated)
}
