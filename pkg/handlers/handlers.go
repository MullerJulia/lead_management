package handlers

import (
	"encoding/json"
	"lead_management/pkg/db"
	"lead_management/pkg/models"
	"lead_management/pkg/utils"
	"log"
	"net/http"
	"strings"
	"time"
)

// CreateClientRequest is used to decode the JSON request payload.
type CreateClientRequest struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	Priority          int    `json:"priority"`
	LeadCapacity      int    `json:"leadCapacity"`
	CurrentLeadCount  int    `json:"currentLeadCount"`
	WorkingHoursStart string `json:"workingHoursStart"`
	WorkingHoursEnd   string `json:"workingHoursEnd"`
}

// CreateClientHandler handles the creation of a new client.
func CreateClientHandler(db *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Unsupported HTTP method", http.StatusMethodNotAllowed)
			return
		}

		var req CreateClientRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		start, err := time.Parse("15:04", req.WorkingHoursStart)
		if err != nil {
			http.Error(w, "Invalid working hours start time format", http.StatusBadRequest)
			return
		}
		end, err := time.Parse("15:04", req.WorkingHoursEnd)
		if err != nil {
			http.Error(w, "Invalid working hours end time format", http.StatusBadRequest)
			return
		}

		clientID := req.ID
		if clientID == "" {
			clientID = utils.GenerateUUID()
		}

		client := models.Client{
			ID:               clientID,
			Name:             req.Name,
			Priority:         req.Priority,
			LeadCapacity:     req.LeadCapacity,
			CurrentLeadCount: req.CurrentLeadCount,
			WorkingHours:     [2]time.Time{start, end},
		}

		if err := db.CreateClient(client); err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				http.Error(w, "Client ID already exists", http.StatusConflict)
				return
			}
			http.Error(w, "Failed to create client", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(client)
	}
}

// GetAllClientsHandler retrieves all client records from the database.
func GetAllClientsHandler(db *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Unsupported HTTP method", http.StatusMethodNotAllowed)
			return
		}
		clients, err := db.GetAllClients()
		if err != nil {
			log.Printf("Error fetching clients: %v", err)
			http.Error(w, "Failed to fetch clients", http.StatusInternalServerError)
			return
		}
		if len(clients) == 0 {
			log.Println("No clients found in handler")
			http.Error(w, "No clients found", http.StatusNotFound)
			return
		}
		log.Printf("Clients fetched in handler: %v", clients)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(clients)
	}
}

// GetClientByIDHandler retrieves a specific client by ID from the database.
func GetClientByIDHandler(db *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Unsupported HTTP method", http.StatusMethodNotAllowed)
			return
		}
		id := strings.TrimPrefix(r.URL.Path, "/client/")
		client, err := db.GetClientByID(id)
		if err != nil {
			http.Error(w, "Failed to fetch client", http.StatusInternalServerError)
			return
		}
		if client == nil {
			http.Error(w, "Client not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(client)
	}
}

// AssignLeadHandler determines the appropriate client for a lead.
func AssignLeadHandler(db *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Unsupported HTTP method", http.StatusMethodNotAllowed)
			return
		}
		client, err := db.GetEligibleClient()
		if err != nil {
			http.Error(w, "No eligible client found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(client)
	}
}
