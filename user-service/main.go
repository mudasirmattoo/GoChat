package main

import (
	"log"
	"net/http"
	"user-service/database"

	"user-service/routes"
)

func main() {

	database.ConnectDB()

	routes.RegisterRoutes()
	// Serve static files from frontend/ directory
	fs := http.FileServer(http.Dir("./frontend"))
	http.Handle("/", fs)
	port := "9080"
	log.Printf("user service running on %s ...\n", port)

	err := http.ListenAndServe(":"+port, nil)

	if err != nil {
		log.Fatalf("Could not start server: %v", err)
	}

}
