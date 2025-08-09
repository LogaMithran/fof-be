package routers

import (
	"friends-of-friends-be/controllers"
	"github.com/gorilla/mux"
	"net/http"
)

func Initialize(router *mux.Router) {
	router.HandleFunc("/health", controllers.HealthController).Methods(http.MethodGet)
	router.HandleFunc("/scores", controllers.GetScores).Methods(http.MethodGet)
	router.HandleFunc("/scores", controllers.UpdateScores).Methods(http.MethodPost)
	router.HandleFunc("/scores/{user_id}", controllers.GetUserScore).Methods(http.MethodGet)
}
