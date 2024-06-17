package models

import (
	"time"
)

// Client represents a client in the system.
type Client struct {
	ID               string       `json:"id"`
	Name             string       `json:"name"`
	Priority         int          `json:"priority"`
	LeadCapacity     int          `json:"leadCapacity"`
	CurrentLeadCount int          `json:"currentLeadCount"`
	WorkingHours     [2]time.Time `json:"workingHours"` // Client opening and closing times
}
