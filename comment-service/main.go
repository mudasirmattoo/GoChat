package main

import (
	"comment-service/database"
	"comment-service/routes"
	"log"
	"net/http"
)

func main() {
	database.ConnectDB()

	r := routes.RegisterRoutes()

	port := "9082"
	log.Printf("Comment service running on %s ...\n", port)

	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
