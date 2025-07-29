package routers

import (
	"fmt"
	"friends-of-friends-be/controllers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func Router() {
	router := mux.NewRouter()
	router.Use(mux.CORSMethodMiddleware(router))
	router.HandleFunc("/ws", controllers.WsController)
	router.HandleFunc("/friends", controllers.GetPeopleWithInXKm).Methods("GET")

	router.HandleFunc("/users", controllers.GetUsersController).Methods("GET")
	router.HandleFunc("/users", controllers.CreateUser).Methods(http.MethodGet, http.MethodPut, http.MethodPost, http.MethodPatch, http.MethodOptions)
	router.HandleFunc("/addFriend", controllers.AddFriendController).Methods("PUT")

	fmt.Println("Server started on 8080")
	http.Handle("/", router)
	if httpErr := http.ListenAndServe(":8080", nil); httpErr != nil {
		log.Fatalf("Error in creating the server")
	}
}
