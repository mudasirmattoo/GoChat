package routes

import (
	"comment-service/handlers"

	"github.com/gorilla/mux"
)

func RegisterRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/create/comment", handlers.CreateComment).Methods("POST")
	r.HandleFunc("/get/comments", handlers.GetComments).Methods("GET")
	r.HandleFunc("/get/comments/{post_id}", handlers.GetCommentsByPostID).Methods("GET")
	r.HandleFunc("/delete/comments/{id}", handlers.DeleteComment).Methods("DELETE")
	r.HandleFunc("/edit/comments/{id}", handlers.EditComment).Methods("PUT")

	return r
}
