package main

import (
	"backend/db"
	"log"
	"net/http"

	"backend/routes"
)

func main() {

	db.ConnectDB()

	mux := routes.RegisterRoutes()
	routes.RegisterRoutes()

	port := "9080"

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../frontend/static"))))

	log.Printf("Starting server on %s...\n", port)
	err := http.ListenAndServe(":"+port, mux)

	if err != nil {
		log.Fatalf("Could not start server: %v", err)
	}

}
