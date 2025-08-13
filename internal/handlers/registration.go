package handlers

import (
    "encoding/json"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "readytorun-backend/internal/database"
    "readytorun-backend/internal/models"
)

func CreateRegistration(w http.ResponseWriter, r *http.Request) {
    err := r.ParseMultipartForm(10 << 20) // 10 MB
    if err != nil {
        http.Error(w, "Could not parse form", http.StatusBadRequest)
        return
    }

    var reg models.Registration
    reg.Fullname = r.FormValue("fullname")
    reg.Dob = r.FormValue("dob")
    reg.Gender = r.FormValue("gender")
    reg.Email = r.FormValue("email")
    reg.Phone = r.FormValue("phone")
    reg.StateOfOrigin = r.FormValue("state_of_origin")
    reg.StateOfResidence = r.FormValue("state_of_residence")
    reg.Education = r.FormValue("education")
    reg.PreviousOffice = r.FormValue("previous_office")
    reg.CardCarryingMember = (r.FormValue("card_carrying_member") == "true")
    reg.Motivation = r.FormValue("motivation")
    reg.PoliticalUnderstanding = r.FormValue("political_understanding")
    reg.AssistanceNeeded = r.FormValue("assistance_needed")
    reg.Availability = r.FormValue("availability")
    reg.PreferredCommunication = r.FormValue("preferred_communication")
    reg.Consent = (r.FormValue("consent") == "true")

    // Save CV
    if file, handler, err := r.FormFile("cv"); err == nil {
        defer file.Close()
        cvPath := filepath.Join("uploads", handler.Filename)
        out, _ := os.Create(cvPath)
        defer out.Close()
        _, _ = out.ReadFrom(file)
        reg.CV = cvPath
    }

    // Save Party Membership Doc
    if file, handler, err := r.FormFile("party_membership_doc"); err == nil {
        defer file.Close()
        docPath := filepath.Join("uploads", handler.Filename)
        out, _ := os.Create(docPath)
        defer out.Close()
        _, _ = out.ReadFrom(file)
        reg.PartyMembershipDoc = docPath
    }

    stmt, err := database.DB.Prepare(`
        INSERT INTO registrations(
            fullname, dob, gender, email, phone, state_of_origin, state_of_residence, education,
            previous_office, card_carrying_member, party_membership_doc, cv, motivation,
            political_understanding, assistance_needed, availability, preferred_communication, consent
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `)
    if err != nil {
        log.Println("Prepare error:", err)
        http.Error(w, "DB error", http.StatusInternalServerError)
        return
    }
    defer stmt.Close()

    _, err = stmt.Exec(
        reg.Fullname, reg.Dob, reg.Gender, reg.Email, reg.Phone, reg.StateOfOrigin, reg.StateOfResidence,
        reg.Education, reg.PreviousOffice, reg.CardCarryingMember, reg.PartyMembershipDoc, reg.CV,
        reg.Motivation, reg.PoliticalUnderstanding, reg.AssistanceNeeded, reg.Availability,
        reg.PreferredCommunication, reg.Consent,
    )
    if err != nil {
        log.Println("Exec error:", err)
        http.Error(w, "Insert failed", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
