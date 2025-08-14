package main

import (
	"log"
	"net/http"
	"readytorun-backend/internal/database"
	"readytorun-backend/internal/handlers"
)

func main() {
	database.ConnectDB()

	// Public
	http.HandleFunc("/api/contact", handlers.CreateContact)
	http.HandleFunc("/api/registration", handlers.CreateRegistration)

	// Admin
	http.HandleFunc("/api/admin/contacts", handlers.GetContacts)
	http.HandleFunc("/api/admin/registrations", handlers.GetRegistrations)

	log.Println("Server running on http://localhost:5001")
	log.Fatal(http.ListenAndServe(":5001", nil))
}
