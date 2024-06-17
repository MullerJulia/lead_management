package handlers

import (
	"lead_management/pkg/db"
	"net/http"
)

// SetupRoutes sets up all the routes for the application.
func SetupRoutes(mux *http.ServeMux, database *db.DB) {
	// Create a new client
	mux.HandleFunc("/client/create", CreateClientHandler(database))

	// Retrieve all clients
	mux.HandleFunc("/client/all", GetAllClientsHandler(database))

	// Retrieve a specific client by their ID
	mux.HandleFunc("/client/", GetClientByIDHandler(database))

	// Endpoint for assigning a lead to a client
	mux.HandleFunc("/client/assign", AssignLeadHandler(database))
}
