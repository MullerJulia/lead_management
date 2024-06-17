package db

import (
	"database/sql"
	"lead_management/pkg/models"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// DB is a wrapper for the SQL database.
type DB struct {
	*sql.DB
}

// InitDB initializes and returns a database object.
func InitDB(dataSourceName string) *DB {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// Create the clients table if it does not already exist
	createTableSQL := `CREATE TABLE IF NOT EXISTS clients (
        id TEXT PRIMARY KEY,
        name TEXT NOT NULL,
        priority INTEGER NOT NULL,
        leadCapacity INTEGER NOT NULL,
        currentLeadCount INTEGER NOT NULL,
        workingHoursStart TEXT NOT NULL,
        workingHoursEnd TEXT NOT NULL
    );`
	if _, err = db.Exec(createTableSQL); err != nil {
		log.Fatalf("Error creating table: %v", err)
	} else {
		log.Println("Clients table created or already exists.")
	}

	return &DB{db}
}

// CreateClient inserts a new client into the database.
func (db *DB) CreateClient(c models.Client) error {
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error beginning transaction: %v", err)
		return err
	}

	stmt, err := tx.Prepare(`INSERT INTO clients (id, name, priority, leadCapacity, currentLeadCount, workingHoursStart, workingHoursEnd) VALUES (?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(c.ID, c.Name, c.Priority, c.LeadCapacity, c.CurrentLeadCount, c.WorkingHours[0].Format("15:04"), c.WorkingHours[1].Format("15:04"))
	if err != nil {
		log.Printf("Error executing statement: %v", err)
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		return err
	}

	log.Printf("Client created: %+v", c)
	return nil
}

// GetAllClients retrieves all clients from the database.
func (db *DB) GetAllClients() ([]models.Client, error) {
	log.Println("Attempting to fetch all clients")
	query := `SELECT id, name, priority, leadCapacity, currentLeadCount, workingHoursStart, workingHoursEnd FROM clients`
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Error querying clients: %v", err)
		return nil, err
	}
	defer rows.Close()

	var clients []models.Client
	for rows.Next() {
		var c models.Client
		var start, end string
		if err := rows.Scan(&c.ID, &c.Name, &c.Priority, &c.LeadCapacity, &c.CurrentLeadCount, &start, &end); err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}

		c.WorkingHours[0], err = time.Parse("15:04", start)
		if err != nil {
			log.Printf("Error parsing start time: %v", err)
			return nil, err
		}
		c.WorkingHours[1], err = time.Parse("15:04", end)
		if err != nil {
			log.Printf("Error parsing end time: %v", err)
			return nil, err
		}

		clients = append(clients, c)
	}
	if err = rows.Err(); err != nil {
		log.Printf("Error with rows: %v", err)
		return nil, err
	}
	if len(clients) == 0 {
		log.Println("No clients found")
	} else {
		log.Printf("Fetched clients: %v", clients)
	}
	return clients, nil
}

// GetClientByID retrieves a client by its ID from the database.
func (db *DB) GetClientByID(id string) (*models.Client, error) {
	query := `SELECT id, name, priority, leadCapacity, currentLeadCount, workingHoursStart, workingHoursEnd FROM clients WHERE id = ?`
	row := db.QueryRow(query, id)

	var c models.Client
	var start, end string
	if err := row.Scan(&c.ID, &c.Name, &c.Priority, &c.LeadCapacity, &c.CurrentLeadCount, &start, &end); err != nil {
		if err == sql.ErrNoRows {
			log.Println("No rows found")
			return nil, nil
		}
		log.Printf("Error scanning row: %v", err)
		return nil, err
	}

	var err error

	c.WorkingHours[0], err = time.Parse("15:04", start)
	if err != nil {
		log.Printf("Error parsing start time: %v", err)
		return nil, err
	}
	c.WorkingHours[1], err = time.Parse("15:04", end)
	if err != nil {
		log.Printf("Error parsing end time: %v", err)
		return nil, err
	}

	log.Printf("Fetched client: %+v", c)
	return &c, nil
}

// GetEligibleClient finds the most eligible client based on priority, current lead count, and working hours.
func (db *DB) GetEligibleClient() (*models.Client, error) {
	log.Println("Attempting to find eligible client for lead")
	currentTime := time.Now().Format("15:04")
	log.Printf("Current Time: %s", currentTime)

	query := `
        SELECT id, name, priority, leadCapacity, currentLeadCount, workingHoursStart, workingHoursEnd 
        FROM clients 
        WHERE 
            (workingHoursStart < workingHoursEnd AND ? BETWEEN workingHoursStart AND workingHoursEnd)
            OR
            (workingHoursStart > workingHoursEnd AND (? >= workingHoursStart OR ? <= workingHoursEnd))
        AND currentLeadCount < leadCapacity 
        ORDER BY priority DESC, currentLeadCount ASC 
        LIMIT 1
    `
	log.Printf("Running query: %s with currentTime: %s", query, currentTime)

	row := db.QueryRow(query, currentTime, currentTime, currentTime)

	var c models.Client
	var start, end string
	err := row.Scan(&c.ID, &c.Name, &c.Priority, &c.LeadCapacity, &c.CurrentLeadCount, &start, &end)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No eligible clients found")
			return nil, nil
		}
		log.Printf("Error querying eligible client: %v", err)
		return nil, err
	}

	c.WorkingHours[0], err = time.Parse("15:04", start)
	if err != nil {
		log.Printf("Error parsing start time: %v", err)
		return nil, err
	}
	c.WorkingHours[1], err = time.Parse("15:04", end)
	if err != nil {
		log.Printf("Error parsing end time: %v", err)
		return nil, err
	}

	log.Printf("Found client: %+v", c)
	return &c, nil
}
