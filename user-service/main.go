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

	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		http.ServeFile(w, r, "frontend/static/"+r.URL.Path[len("/static/"):])
	})

	port := "9080"
	log.Printf("user service running on %s ...\n", port)

	err := http.ListenAndServe(":"+port, nil)

	if err != nil {
		log.Fatalf("Could not start server: %v", err)
	}

}
