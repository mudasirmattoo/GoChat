package routes

import (
	"post/handlers"

	"github.com/gorilla/mux"
)

func RegisterRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", handlers.ServeHome).Methods("GET")
	r.HandleFunc("/create/post", handlers.CreatePost).Methods("POST")
	r.HandleFunc("/get/posts", handlers.GetPost).Methods("GET")
	r.HandleFunc("/get/posts/{id}", handlers.GetPostByID).Methods("GET")
	r.HandleFunc("/update/posts/{id}", handlers.UpdatePost).Methods("PUT")
	r.HandleFunc("/delete/posts/{id}", handlers.DeletePost).Methods("DELETE")

	return r
}
