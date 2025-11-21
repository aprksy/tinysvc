package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpdelivery "github.com/aprksy/tinysvc/internal/delivery/http"
	"github.com/aprksy/tinysvc/internal/infrastructure/config"
	"github.com/aprksy/tinysvc/internal/infrastructure/persistence/sqlite"
	"github.com/aprksy/tinysvc/internal/usecase"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize database
	db, err := sqlite.InitDB(cfg.Database.Path)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer db.Close()

	// Initialize repositories
	pasteRepo := sqlite.NewPasteRepository(db)

	// Initialize services
	pasteService := usecase.NewPasteService(pasteRepo)
	ipService := usecase.NewIPService()

	// Setup cleanup job for expired pastes
	go runCleanupJob(pasteService)

	// Initialize HTTP router
	router := httpdelivery.NewRouter(pasteService, ipService)
	handler := router.SetupRoutes()

	// Setup HTTP server
	srv := &http.Server{
		Addr:         cfg.ServerAddress(),
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("Starting server on %s", cfg.ServerAddress())
		serverErrors <- srv.ListenAndServe()
	}()

	// Wait for interrupt signal
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case sig := <-shutdown:
		log.Printf("Received signal %v, starting graceful shutdown", sig)

		// Graceful shutdown with 30 second timeout
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			srv.Close()
			return fmt.Errorf("failed to gracefully shutdown server: %w", err)
		}

		log.Println("Server stopped")
	}

	return nil
}

// runCleanupJob runs periodic cleanup of expired pastes
func runCleanupJob(pasteService usecase.PasteService) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		count, err := pasteService.CleanupExpired(ctx)
		cancel()

		if err != nil {
			log.Printf("Failed to cleanup expired pastes: %v", err)
			continue
		}

		if count > 0 {
			log.Printf("Cleaned up %d expired pastes", count)
		}
	}
}
