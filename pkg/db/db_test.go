package db

import (
	"lead_management/pkg/models"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestOpenDB(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
}

// MockDB is a mock for the Database
type MockDB struct {
	mock.Mock
}

func (m *MockDB) CreateClient(c models.Client) error {
	args := m.Called(c)
	return args.Error(0)
}

func TestCreateClient(t *testing.T) {
	type testCase struct {
		name          string
		client        models.Client
		expectedError bool
		setupMock     func(*MockDB, models.Client)
	}

	tests := []testCase{
		{
			name: "Success",
			client: models.Client{
				ID: "test-id-success", Name: "Test Client", Priority: 1, LeadCapacity: 100, CurrentLeadCount: 0,
			},
			expectedError: false,
			setupMock: func(m *MockDB, c models.Client) {
				m.On("CreateClient", mock.AnythingOfType("models.Client")).Return(nil)
			},
		},
		{
			name: "SQL Error",
			client: models.Client{
				ID: "test-id-error", Name: "Test Client",
			},
			expectedError: true,
			setupMock: func(m *MockDB, c models.Client) {
				m.On("CreateClient", mock.AnythingOfType("models.Client")).Return(assert.AnError)
			},
		},
		{
			name: "Incomplete Data",
			client: models.Client{
				ID: "test-id-incomplete", // Name is missing
			},
			expectedError: true,
			setupMock: func(m *MockDB, c models.Client) {
				m.On("CreateClient", mock.AnythingOfType("models.Client")).Return(assert.AnError)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockDB := new(MockDB)
			tc.setupMock(mockDB, tc.client)

			err := mockDB.CreateClient(tc.client)

			mockDB.AssertExpectations(t)
			if tc.expectedError {
				require.Error(t, err, "Expected an error for test case: "+tc.name)
			} else {
				require.NoError(t, err, "Did not expect an error for test case: "+tc.name)
			}
		})
	}
}

func TestGetAllClients(t *testing.T) {
	// Setup in-memory DB
	database := InitDB(":memory:")
	defer database.Close()

	// Define test cases
	tests := []struct {
		name          string
		setupData     func(*DB)
		expectedData  []models.Client
		expectedError bool
	}{
		{
			name: "Empty Result",
			setupData: func(db *DB) {
				// No setup needed for empty result
			},
			expectedData:  []models.Client{},
			expectedError: false,
		},
		{
			name: "Non-Empty Result",
			setupData: func(db *DB) {
				setupEligibleClientsDatabase(db, []models.Client{
					{
						ID:               "1",
						Name:             "Test Client",
						Priority:         1,
						LeadCapacity:     100,
						CurrentLeadCount: 50,
						WorkingHours:     [2]time.Time{parseTime("09:00"), parseTime("17:00")},
					},
				})
			},
			expectedData: []models.Client{
				{
					ID:               "1",
					Name:             "Test Client",
					Priority:         1,
					LeadCapacity:     100,
					CurrentLeadCount: 50,
					WorkingHours:     [2]time.Time{parseTime("09:00"), parseTime("17:00")},
				},
			},
			expectedError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupData(database)

			clients, err := database.GetAllClients()

			if tc.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedData, clients, "Fetched clients do not match expected")
			}
		})
	}
}

func TestGetClientByID(t *testing.T) {
	// Setup in-memory DB
	database := InitDB(":memory:")
	defer database.Close()

	// Pre-populate the database with some data
	setupDatabase(database)

	// Define test cases
	tests := []struct {
		name         string
		clientID     string
		expectedData *models.Client
		expectedErr  bool
	}{
		{
			name:     "Successful retrieval",
			clientID: "1",
			expectedData: &models.Client{
				ID:               "1",
				Name:             "Test Client",
				Priority:         1,
				LeadCapacity:     100,
				CurrentLeadCount: 50,
				WorkingHours:     [2]time.Time{parseTime("09:00"), parseTime("17:00")},
			},
			expectedErr: false,
		},
		{
			name:         "Client not found",
			clientID:     "2",
			expectedData: nil,
			expectedErr:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client, err := database.GetClientByID(tc.clientID)

			if tc.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedData, client)
			}
		})
	}
}

func TestGetEligibleClient(t *testing.T) {
	// Setup in-memory DB
	database := InitDB(":memory:")
	defer database.Close()

	// Define test cases
	tests := []struct {
		name         string
		setupData    func(*DB)
		expectedData *models.Client
		expectedErr  bool
	}{
		{
			name: "Highest priority client available during working hours",
			setupData: func(db *DB) {
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
			expectedData: &models.Client{
				ID:               "1",
				Name:             "High Priority Client",
				Priority:         10,
				LeadCapacity:     100,
				CurrentLeadCount: 20,
				WorkingHours:     [2]time.Time{parseTime(time.Now().Add(-1 * time.Hour).Format("15:04")), parseTime(time.Now().Add(1 * time.Hour).Format("15:04"))},
			},
			expectedErr: false,
		},
		{
			name: "Clients with same priority but different lead counts",
			setupData: func(db *DB) {
				now := time.Now()
				start := now.Add(-1 * time.Hour).Format("15:04")
				end := now.Add(1 * time.Hour).Format("15:04")
				setupEligibleClientsDatabase(db, []models.Client{
					{
						ID:               "3",
						Name:             "Client with Fewer Leads",
						Priority:         10,
						LeadCapacity:     100,
						CurrentLeadCount: 5,
						WorkingHours:     [2]time.Time{parseTime(start), parseTime(end)},
					},
					{
						ID:               "4",
						Name:             "Client with More Leads",
						Priority:         10,
						LeadCapacity:     100,
						CurrentLeadCount: 15,
						WorkingHours:     [2]time.Time{parseTime(start), parseTime(end)},
					},
				})
			},
			expectedData: &models.Client{
				ID:               "3",
				Name:             "Client with Fewer Leads",
				Priority:         10,
				LeadCapacity:     100,
				CurrentLeadCount: 5,
				WorkingHours:     [2]time.Time{parseTime(time.Now().Add(-1 * time.Hour).Format("15:04")), parseTime(time.Now().Add(1 * time.Hour).Format("15:04"))},
			},
			expectedErr: false,
		},
		/*{
			name: "No eligible clients",
			setupData: func(db *DB) {
				setupEligibleClientsDatabase(db, []models.Client{
					{
						ID:               "5",
						Name:             "Unavailable Client",
						Priority:         10,
						LeadCapacity:     100,
						CurrentLeadCount: 50,
						WorkingHours:     [2]time.Time{parseTime(time.Now().Add(-2 * time.Hour).Format("15:04")), parseTime(time.Now().Add(-1 * time.Hour).Format("15:04"))},
					},
				})
			},
			expectedData: nil,
			expectedErr:  false,
		},*/
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup database with test-specific data
			tc.setupData(database)

			client, err := database.GetEligibleClient()

			if tc.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedData, client)
			}
		})
	}
}

func setupEligibleClientsDatabase(database *DB, clients []models.Client) {
	for _, client := range clients {
		err := database.CreateClient(client)
		if err != nil {
			panic("Failed to setup database: " + err.Error())
		}
	}
}

func setupDatabase(database *DB) {
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
	timeParsed, _ := time.Parse("15:04", t)
	return timeParsed
}
