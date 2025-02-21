package routes

import (
	"net/http"
	"user-service/handlers"
)

func RegisterRoutes() {
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
}
