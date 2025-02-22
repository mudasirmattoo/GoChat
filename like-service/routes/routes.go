package routes

import (
	"like-service/handlers"

	"github.com/gorilla/mux"
)

func RegisterRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/like/add", handlers.AddLike).Methods("POST")
	r.HandleFunc("/like/{post_id}", handlers.GetLikesByPostID).Methods("GET")
	r.HandleFunc("/like/remove/{id}", handlers.RemoveLike).Methods("DELETE")

	return r
}
