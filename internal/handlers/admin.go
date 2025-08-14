package handlers

import (
    "encoding/json"
    "log"
    "net/http"
    "readytorun-backend/internal/database"
    "readytorun-backend/internal/models"
)

// GET /api/admin/contacts
func GetContacts(w http.ResponseWriter, r *http.Request) {
    rows, err := database.DB.Query(`
        SELECT id, fullname, email, phone, subject, message, created_at
        FROM contacts
        ORDER BY created_at DESC
    `)
    if err != nil {
        log.Println("Query error:", err)
        http.Error(w, "DB error", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var contacts []models.Contact
    for rows.Next() {
        var c models.Contact
        err := rows.Scan(&c.ID, &c.Fullname, &c.Email, &c.Phone, &c.Subject, &c.Message, &c.CreatedAt)
        if err != nil {
            log.Println("Scan error:", err)
            continue
        }
        contacts = append(contacts, c)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(contacts)
}

// GET /api/admin/registrations
func GetRegistrations(w http.ResponseWriter, r *http.Request) {
    rows, err := database.DB.Query(`
        SELECT id, fullname, dob, gender, email, phone, state_of_origin, state_of_residence,
               education, previous_office, card_carrying_member, party_membership_doc_link, cv_link, motivation,
               political_understanding, assistance_needed, availability, preferred_communication,
               consent, created_at
        FROM registrations
        ORDER BY created_at DESC
    `)
    if err != nil {
        log.Println("Query error:", err)
        http.Error(w, "DB error", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var regs []models.Registration
    for rows.Next() {
        var r models.Registration
        err := rows.Scan(
            &r.ID, &r.Fullname, &r.Dob, &r.Gender, &r.Email, &r.Phone, &r.StateOfOrigin, &r.StateOfResidence,
            &r.Education, &r.PreviousOffice, &r.CardCarryingMember, &r.PartyMembershipDocLink, &r.CVLink,
            &r.Motivation, &r.PoliticalUnderstanding, &r.AssistanceNeeded, &r.Availability,
            &r.PreferredCommunication, &r.Consent, &r.CreatedAt,
        )
        if err != nil {
            log.Println("Scan error:", err)
            continue
        }
        regs = append(regs, r)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(regs)
}
