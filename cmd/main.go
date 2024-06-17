package main

import (
	"log"
	"net/http"

	"lead_management/pkg/db"
	"lead_management/pkg/handlers"
)

func main() {
	database := db.InitDB("lead_management.db")

	mux := http.NewServeMux()
	handlers.SetupRoutes(mux, database)

	log.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
