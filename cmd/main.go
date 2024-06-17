package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"lead_management/pkg/db"
	"lead_management/pkg/handlers"
)

const (
	address         = ":8080"
	shutdownTimeout = 5 * time.Second
	dbFilePath      = "lead_management.db"
)

// setupServer initializes the HTTP server and sets up the routes.
func setupServer(database *db.DB) *http.Server {
	mux := http.NewServeMux()
	handlers.SetupRoutes(mux, database)

	return &http.Server{
		Addr:    address,
		Handler: mux,
	}
}

// waitForShutdown waits for an interrupt signal and attempts a graceful shutdown.
func waitForShutdown(server *http.Server) {
	// Channel to listen for OS signals.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	// Create a context with a timeout for the graceful shutdown.
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// Attempt graceful shutdown.
	log.Println("Shutting down server...")
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}

func main() {
	database := db.InitDB(dbFilePath)
	server := setupServer(database)

	// Start the server in a goroutine.
	go func() {
		log.Printf("Server started on %s\n", address)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v\n", address, err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server.
	waitForShutdown(server)
}
