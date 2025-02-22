package main

import (
	"like-service/database"
	"like-service/routes"
	"log"
	"net/http"
)

func main() {
	database.ConnectDB()

	r := routes.RegisterRoutes()

	port := "9083"
	log.Printf("Like service running on %s ...\n", port)

	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
