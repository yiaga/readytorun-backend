package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"readytorun-backend/internal/database"
	"readytorun-backend/internal/handlers"
	"syscall"
	"time"

    "github.com/joho/godotenv"
)

func main() {

	_ = godotenv.Load()
	// Initialize database connection
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Setup routes
	mux := http.NewServeMux()

	// API v1 routes
	mux.HandleFunc("/api/registrations", handlers.RegistrationHandler(db))
	mux.HandleFunc("/api/contacts", handlers.ContactHandler(db))
	mux.HandleFunc("/api/volunteers", handlers.VolunteerHandler(db))
	mux.HandleFunc("/api/registration", handlers.GetRegistration(db))
	mux.HandleFunc("/api/contact", handlers.GetContact(db))
	mux.HandleFunc("/api/volunteer", handlers.GetVolunteer(db))

	// Health check route
	mux.HandleFunc("/health/", func(w http.ResponseWriter, r *http.Request) {
		if err := db.Ping(); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{
				"status":  "down",
				"message": "database connection failed",
			})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "up"})
	})

	// Wrap mux with logging middleware
	handler := loggingMiddleware(mux)

	// Setup HTTP server
	srv := &http.Server{
		Addr:         ":" + getPort(),
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("🚀 Server starting on http://localhost%s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Server failed: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server exited properly")
}

// ------------------ Helpers ------------------

// getPort reads port from env or defaults to 8080
func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}

// loggingMiddleware logs all incoming requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("%s %s %s", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
		log.Printf("Completed in %v", time.Since(start))
	})
}

// writeJSON writes JSON responses with proper headers
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("❌ Failed to write JSON response: %v", err)
	}
}
