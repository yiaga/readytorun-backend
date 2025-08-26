package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
	"strconv"

	"readytorun-backend/internal/models"
)

func ContactHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
			case http.MethodPost:
				var contact models.Contact
				if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
					http.Error(w, "invalid JSON", http.StatusBadRequest)
					return
				}

				contact.CreatedAt = time.Now()

				query := `INSERT INTO contacts (name, email, message, subject, created_at) VALUES ($1,$2,$3,$4,$5) RETURNING id`
				if err := db.QueryRow(query, contact.Name, contact.Email, contact.Message, contact.Subject, contact.CreatedAt).Scan(&contact.ID); err != nil {
					http.Error(w, "failed to insert", http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(contact)
				return
			
			case http.MethodGet:

				rows, err := db.Query("SELECT id, name, email, message, subject, created_at FROM contacts")
				if err != nil {
					http.Error(w, "failed to fetch", http.StatusInternalServerError)
					return
				}
				defer rows.Close()

				var contacts []models.Contact
				for rows.Next() {
					var c models.Contact
					if err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.Message, &c.Subject, &c.CreatedAt); err != nil {
						http.Error(w, "scan error", http.StatusInternalServerError)
						return
					}
					contacts = append(contacts, c)
				}

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(contacts)
				return

			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
		}
		
	}
}

// Get a single contact by ID
func GetContact(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get ID from query param like /contact?id=1
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, "id is required", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		var contact models.Contact
		query := `SELECT id, name, email, subject, message, created_at FROM contacts WHERE id=$1`
		err = db.QueryRow(query, id).Scan(&contact.ID, &contact.Name, &contact.Email, &contact.Subject, &contact.Message, &contact.CreatedAt)
		if err == sql.ErrNoRows {
			http.Error(w, "contact not found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, "failed to fetch contact", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(contact)
		
	}
}
