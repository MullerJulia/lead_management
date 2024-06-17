package handlers

import (
	"bytes"
	"encoding/json"
	"lead_management/pkg/db"
	"lead_management/pkg/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateClientHandler(t *testing.T) {
	// Test case struct to hold the parameters for each test
	type testCase struct {
		name         string
		method       string
		body         []byte
		expectedCode int
	}

	// Test data
	clientReq := CreateClientRequest{
		Name:              "New Client",
		Priority:          5,
		LeadCapacity:      10,
		CurrentLeadCount:  0,
		WorkingHoursStart: "09:00",
		WorkingHoursEnd:   "17:00",
	}
	clientBody, _ := json.Marshal(clientReq)
	invalidBody := []byte("{invalid json}")

	tests := []testCase{
		{
			name:         "Successful creation",
			method:       "POST",
			body:         clientBody,
			expectedCode: http.StatusCreated,
		},
		{
			name:         "Incorrect HTTP method",
			method:       "GET",
			body:         nil,
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "Invalid JSON data",
			method:       "POST",
			body:         invalidBody,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			database := db.InitDB(":memory:")
			defer database.Close()

			handler := CreateClientHandler(database)

			req, err := http.NewRequest(tc.method, "/client/create", bytes.NewReader(tc.body))
			if err != nil {
				t.Fatalf("Could not create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedCode, rr.Code, "Handler returned wrong status code")

			if tc.expectedCode == http.StatusCreated {
				var responseClient models.Client
				err := json.NewDecoder(rr.Body).Decode(&responseClient)
				if err != nil {
					t.Fatalf("Could not parse response: %v", err)
				}
				assert.Equal(t, clientReq.Name, responseClient.Name)
				assert.NotEmpty(t, responseClient.ID)
			} else {
				// Debugging output for non-201 responses
				t.Logf("Response body: %s", rr.Body.String())
			}
		})
	}
}

func TestGetAllClientsHandler(t *testing.T) {
	// Setup in-memory DB
	database := db.InitDB(":memory:")
	defer database.Close()

	// Pre-populate the database with some data
	setupDatabase(database)

	// Define test cases
	tests := []struct {
		name         string
		method       string
		expectedCode int
		expectedData []models.Client
	}{
		{
			name:         "Successful retrieval",
			method:       "GET",
			expectedCode: http.StatusOK,
			expectedData: []models.Client{
				{ID: "1", Name: "Test Client", Priority: 1, LeadCapacity: 100, CurrentLeadCount: 50, WorkingHours: [2]time.Time{parseTime("09:00"), parseTime("17:00")}},
			},
		},
		{
			name:         "Incorrect HTTP method",
			method:       "POST",
			expectedCode: http.StatusMethodNotAllowed,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handler := GetAllClientsHandler(database)

			req, _ := http.NewRequest(tc.method, "/clients", nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedCode, rr.Code)

			if tc.expectedCode == http.StatusOK {
				var clients []models.Client
				err := json.Unmarshal(rr.Body.Bytes(), &clients)
				require.NoError(t, err)
				assert.Equal(t, tc.expectedData, clients)
			}
		})
	}
}

func TestGetClientByIDHandler(t *testing.T) {
	// Setup in-memory DB
	database := db.InitDB(":memory:")
	defer database.Close()

	// Pre-populate the database with some data
	setupDatabase(database)

	// Define test cases
	tests := []struct {
		name         string
		method       string
		url          string
		expectedCode int
		expectedData *models.Client
	}{
		{
			name:         "Successful retrieval",
			method:       "GET",
			url:          "/client/1",
			expectedCode: http.StatusOK,
			expectedData: &models.Client{ID: "1", Name: "Test Client", Priority: 1, LeadCapacity: 100, CurrentLeadCount: 50, WorkingHours: [2]time.Time{parseTime("09:00"), parseTime("17:00")}},
		},
		{
			name:         "Client not found",
			method:       "GET",
			url:          "/client/2",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "Incorrect HTTP method",
			method:       "POST",
			url:          "/client/1",
			expectedCode: http.StatusMethodNotAllowed,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handler := GetClientByIDHandler(database)

			req, _ := http.NewRequest(tc.method, tc.url, nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedCode, rr.Code)

			if tc.expectedCode == http.StatusOK {
				var client models.Client
				err := json.Unmarshal(rr.Body.Bytes(), &client)
				require.NoError(t, err)
				assert.Equal(t, tc.expectedData, &client)
			}
		})
	}
}

func TestAssignLeadHandler(t *testing.T) {
	// Setup in-memory DB
	database := db.InitDB(":memory:")
	defer database.Close()

	// Define test cases
	tests := []struct {
		name         string
		method       string
		setupData    func(*db.DB)
		expectedCode int
		expectedData *models.Client
	}{
		{
			name:   "Successful assignment",
			method: "GET",
			setupData: func(db *db.DB) {
				now := time.Now()
				start := now.Add(-1 * time.Hour).Format("15:04")
				end := now.Add(1 * time.Hour).Format("15:04")
				setupEligibleClientsDatabase(db, []models.Client{
					{
						ID:               "1",
						Name:             "High Priority Client",
						Priority:         10,
						LeadCapacity:     100,
						CurrentLeadCount: 20,
						WorkingHours:     [2]time.Time{parseTime(start), parseTime(end)},
					},
					{
						ID:               "2",
						Name:             "Low Priority Client",
						Priority:         5,
						LeadCapacity:     100,
						CurrentLeadCount: 10,
						WorkingHours:     [2]time.Time{parseTime(start), parseTime(end)},
					},
				})
			},
			expectedCode: http.StatusOK,
			expectedData: &models.Client{
				ID:               "1",
				Name:             "High Priority Client",
				Priority:         10,
				LeadCapacity:     100,
				CurrentLeadCount: 20,
				WorkingHours:     [2]time.Time{parseTime(time.Now().Add(-1 * time.Hour).Format("15:04")), parseTime(time.Now().Add(1 * time.Hour).Format("15:04"))},
			},
		},
		/*{
			name:   "No eligible clients",
			method: "GET",
			setupData: func(db *db.DB) {
				setupEligibleClientsDatabase(db, []models.Client{
					{
						ID:               "3",
						Name:             "Unavailable Client",
						Priority:         10,
						LeadCapacity:     100,
						CurrentLeadCount: 50,
						WorkingHours:     [2]time.Time{parseTime("00:00"), parseTime("01:00")}, // Time range that is always outside of current time
					},
				})
			},
			expectedCode: http.StatusNotFound,
			expectedData: nil,
		},*/
		{
			name:         "Incorrect HTTP method",
			method:       "POST",
			setupData:    func(db *db.DB) {}, // No setup needed for this test
			expectedCode: http.StatusMethodNotAllowed,
			expectedData: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupData(database)

			handler := AssignLeadHandler(database)

			req, _ := http.NewRequest(tc.method, "/client/assign", nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedCode, rr.Code)

			if tc.expectedCode == http.StatusOK {
				var client models.Client
				err := json.Unmarshal(rr.Body.Bytes(), &client)
				require.NoError(t, err)
				assert.Equal(t, tc.expectedData, &client)
			}
		})
	}
}

// Helper function to set up the eligible clients database
func setupEligibleClientsDatabase(database *db.DB, clients []models.Client) {
	for _, client := range clients {
		err := database.CreateClient(client)
		if err != nil {
			panic("Failed to setup database: " + err.Error())
		}
	}
}

func setupDatabase(database *db.DB) {
	client := models.Client{
		ID:               "1",
		Name:             "Test Client",
		Priority:         1,
		LeadCapacity:     100,
		CurrentLeadCount: 50,
		WorkingHours:     [2]time.Time{parseTime("09:00"), parseTime("17:00")},
	}
	err := database.CreateClient(client)
	if err != nil {
		panic("Failed to setup database: " + err.Error())
	}
}

func parseTime(t string) time.Time {
	parsedTime, _ := time.Parse("15:04", t)
	return parsedTime
}
