package main

import (
	"log"
	"net/http"
	"readytorun-backend/internal/database"
	"readytorun-backend/internal/handlers"
	"readytorun-backend/internal/middleware"
)

func main() {
	database.ConnectDB()
	mux := http.NewServeMux()

	// Public
	mux.HandleFunc("/api/contact", handlers.CreateContact)
	mux.HandleFunc("/api/registration", handlers.CreateRegistration)

	// Admin
	mux.HandleFunc("/api/admin/contacts", handlers.GetContacts)
	mux.HandleFunc("/api/admin/registrations", handlers.GetRegistrations)


	handler := middleware.CORSMiddleware(mux)


	log.Println("Server running on http://localhost:5001")
	if err := http.ListenAndServe(":5001", handler); err != nil {
        log.Fatal(err)
    }
}
