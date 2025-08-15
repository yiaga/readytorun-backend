package handlers

import (
    "encoding/json"
    "log"
    "net/http"
    "readytorun-backend/internal/database"
    "readytorun-backend/internal/models"
)

func CreateContact(w http.ResponseWriter, r *http.Request) {
    var contact models.Contact
    err := json.NewDecoder(r.Body).Decode(&contact)
    if err != nil {
        http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
        return
    }
    
    // ✅ Validate subject before hitting the DB
    allowedSubjects := map[string]bool{
        "registration support": true,
        "volunteer":            true,
        "general inquiry":      true,
    }

    if !allowedSubjects[contact.Subject] {
        http.Error(w, "Invalid subject value", http.StatusBadRequest)
        return
    }
    stmt, err := database.DB.Prepare(`
        INSERT INTO contacts(fullname, email, phone, subject, message)
        VALUES (?, ?, ?, ?, ?)
    `)
    if err != nil {
        log.Println("Prepare error:", err)
        http.Error(w, "DB error", http.StatusInternalServerError)
        return
    }
    defer stmt.Close()

    _, err = stmt.Exec(contact.Fullname, contact.Email, contact.Phone, contact.Subject, contact.Message)
    if err != nil {
        log.Println("Exec error:", err)
        http.Error(w, "Insert failed", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
