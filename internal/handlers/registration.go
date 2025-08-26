package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"github.com/lib/pq"

	"readytorun-backend/internal/models"
)

// RegistrationHandler handles incoming registration requests
func RegistrationHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
			case http.MethodPost:
				var reg models.Registration
				if err := json.NewDecoder(r.Body).Decode(&reg); err != nil {
					http.Error(w, fmt.Sprintf("Invalid request payload: %v", err), http.StatusBadRequest)
					return
				}

				// Basic validation
				if reg.Fullname == "" || reg.Email == "" {
					http.Error(w, "fullname and email are required", http.StatusBadRequest)
					return
				}

				reg.CreatedAt = time.Now()

				// Insert into DB
				query := `
					INSERT INTO registrations (
						fullname, dob, gender, email, phone,
						state_of_origin, state_of_residence, education, previous_office, interested_office,
						previous_contest, card_carrying_member, party_membership_doc_link, motivation,
						political_understanding, assistance_needed, other_support,
						preferred_communication, consent, created_at
					) VALUES (
						$1, $2, $3, $4, $5,
						$6, $7, $8, $9, $10,
						$11, $12, $13, $14, $15,
						$16, $17, $18, $19,
						$20
					) RETURNING id
				`

				err := db.QueryRow(
					query,
					reg.Fullname,
					reg.Dob,
					reg.Gender,
					reg.Email,
					reg.Phone,
					reg.StateOfOrigin,
					reg.StateOfResidence,
					reg.Education,
					reg.PreviousOffice,
					reg.InterestedOffice,
					reg.PreviousContest,
					reg.CardCarryingMember,
					reg.PartyMembershipDocLink,
					reg.Motivation,
					reg.PoliticalUnderstanding,
					pq.Array(reg.AssistanceNeeded),
					reg.OtherSupport,
					reg.PreferredCommunication,
					reg.Consent,
					reg.CreatedAt,
				).Scan(&reg.ID)

				if err != nil {
					http.Error(w, fmt.Sprintf("Failed to insert record: %v", err), http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(reg)
				return

			case http.MethodGet:
				query := `
					SELECT id, fullname, dob, gender, email, phone,
						state_of_origin, state_of_residence, education,
						previous_office, interested_office, previous_contest,
						card_carrying_member, party_membership_doc_link, motivation,
						political_understanding, assistance_needed, other_support,
						preferred_communication, consent, created_at
					FROM registrations
					ORDER BY created_at DESC
				`

				rows, err := db.Query(query)
				if err != nil {
					http.Error(w, "failed to fetch: "+err.Error(), http.StatusInternalServerError)
					return
				}
				defer rows.Close()

				var registrations []models.Registration

				for rows.Next() {
					var reg models.Registration
					var assistance []string

					if err := rows.Scan(
						&reg.ID,
						&reg.Fullname,
						&reg.Dob,
						&reg.Gender,
						&reg.Email,
						&reg.Phone,
						&reg.StateOfOrigin,
						&reg.StateOfResidence,
						&reg.Education,
						&reg.PreviousOffice,
						&reg.InterestedOffice,
						&reg.PreviousContest,
						&reg.CardCarryingMember,
						&reg.PartyMembershipDocLink,
						&reg.Motivation,
						&reg.PoliticalUnderstanding,
						pq.Array(&assistance),
						&reg.OtherSupport,
						&reg.PreferredCommunication,
						&reg.Consent,
						&reg.CreatedAt,
					); err != nil {
						http.Error(w, "failed to scan: "+err.Error(), http.StatusInternalServerError)
						return
					}

					reg.AssistanceNeeded = assistance
					registrations = append(registrations, reg)
				}

				if err := rows.Err(); err != nil {
					http.Error(w, "error iterating rows: "+err.Error(), http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(registrations)
				return

			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
		}

	}
}

// GetRegistration fetches a single registration by ID
func GetRegistration(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		var reg models.Registration
		var assistance []string

		query := `
			SELECT id, fullname, dob, gender, email, phone,
			       state_of_origin, state_of_residence, education,
			       previous_office, interested_office, previous_contest,
			       card_carrying_member, party_membership_doc_link, motivation,
			       political_understanding, assistance_needed, other_support,
			       preferred_communication, consent, created_at
			FROM registrations
			WHERE id = $1
		`

		err = db.QueryRow(query, id).Scan(
			&reg.ID,
			&reg.Fullname,
			&reg.Dob,
			&reg.Gender,
			&reg.Email,
			&reg.Phone,
			&reg.StateOfOrigin,
			&reg.StateOfResidence,
			&reg.Education,
			&reg.PreviousOffice,
			&reg.InterestedOffice,
			&reg.PreviousContest,
			&reg.CardCarryingMember,
			&reg.PartyMembershipDocLink,
			&reg.Motivation,
			&reg.PoliticalUnderstanding,
			pq.Array(&assistance),
			&reg.OtherSupport,
			&reg.PreferredCommunication,
			&reg.Consent,
			&reg.CreatedAt,
		)

		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "registration not found", http.StatusNotFound)
				return
			}
			http.Error(w, "failed to fetch: "+err.Error(), http.StatusInternalServerError)
			return
		}
		reg.AssistanceNeeded = assistance

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(reg)
	}
}
