package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
	"strconv"

	"github.com/lib/pq"
	"readytorun-backend/internal/models"
)

func VolunteerHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		switch r.Method{
			case http.MethodPost:
				var vol models.Volunteer
				if err := json.NewDecoder(r.Body).Decode(&vol); err != nil {
					http.Error(w, "invalid request payload", http.StatusBadRequest)
					return
				}

				// Basic validation
				if vol.FullName == "" || vol.Email == "" {
					http.Error(w, "fullname, and email are required", http.StatusBadRequest)
					return
				}

				now := time.Now()
				vol.CreatedAt = now
				vol.UpdatedAt = now

				query := `
					INSERT INTO volunteers (
						full_name, email, phone, location, skills, created_at, updated_at
					) VALUES ($1,$2,$3,$4,$5,$6,$7)
					RETURNING id
				`

				if err := db.QueryRow(
					query,
					vol.FullName,
					vol.Email,
					vol.Phone,
					vol.Location,
					pq.Array(vol.Skills),
					vol.CreatedAt,
					vol.UpdatedAt,
				).Scan(&vol.ID); err != nil {
					http.Error(w, "failed to insert: "+err.Error(), http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(vol)
				return

			case http.MethodGet:
				query := `
					SELECT id, full_name, email, phone, location, skills, created_at, updated_at
					FROM volunteers ORDER BY created_at DESC
				`

				rows, err := db.Query(query)
				if err != nil {
					http.Error(w, "failed to fetch: "+err.Error(), http.StatusInternalServerError)
					return
				}
				defer rows.Close()

				var volunteers []models.Volunteer

				for rows.Next() {
					var vol models.Volunteer
					var skills []string

					if err := rows.Scan(
						&vol.ID,
						&vol.FullName,
						&vol.Email,
						&vol.Phone,
						&vol.Location,
						pq.Array(&skills),
						&vol.CreatedAt,
						&vol.UpdatedAt,
					); err != nil {
						http.Error(w, "failed to scan: "+err.Error(), http.StatusInternalServerError)
						return
					}
					vol.Skills = skills
					volunteers = append(volunteers, vol)
				}

				if err := rows.Err(); err != nil {
					http.Error(w, "error iterating rows: "+err.Error(), http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(volunteers)
				return
			
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
		}

	}
}

// GetVolunteer fetches a volunteer by ID
func GetVolunteer(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Expecting /volunteers/{id}, so extract ID
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

		var vol models.Volunteer
		query := `
			SELECT id, full_name, email, phone, location, skills, created_at, updated_at
			FROM volunteers WHERE id = $1
		`

		row := db.QueryRow(query, id)
		var skills []string
		if err := row.Scan(
			&vol.ID,
			&vol.FullName,
			&vol.Email,
			&vol.Phone,
			&vol.Location,
			pq.Array(&skills),
			&vol.CreatedAt,
			&vol.UpdatedAt,
		); err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "volunteer not found", http.StatusNotFound)
				return
			}
			http.Error(w, "failed to fetch: "+err.Error(), http.StatusInternalServerError)
			return
		}
		vol.Skills = skills

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(vol)
	}
}
