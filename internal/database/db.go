package database

import (
    "database/sql"
    "log"
    _ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func ConnectDB() {
    var err error
    DB, err = sql.Open("sqlite3", "./readytorun.db")
    if err != nil {
        log.Fatal("Failed to connect to SQLite:", err)
    }

    createTables()
}

func createTables() {
    contactTable := `
    CREATE TABLE IF NOT EXISTS contacts (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        fullname TEXT NOT NULL,
        email TEXT NOT NULL,
        phone TEXT,
        subject TEXT CHECK (subject IN ('registration support', 'volunteer', 'general inquiry')),
        message TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );`

    registrationTable := `
    CREATE TABLE IF NOT EXISTS registrations (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        fullname TEXT NOT NULL,
        dob TEXT NOT NULL,
        gender TEXT,
        email TEXT NOT NULL,
        phone TEXT,
        state_of_origin TEXT,
        state_of_residence TEXT,
        education TEXT,
        previous_office TEXT,
        card_carrying_member BOOLEAN,
        party_membership_doc_link TEXT,
        cv_link TEXT,
        motivation TEXT,
        political_understanding TEXT CHECK (political_understanding IN ('beginner', 'intermediate', 'advance', 'expert')),
        assistance_needed TEXT,
        availability TEXT,
        preferred_communication TEXT,
        consent BOOLEAN DEFAULT 0,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );`

    _, err := DB.Exec(contactTable)
    if err != nil {
        log.Fatal("Failed to create contacts table:", err)
    }

    _, err = DB.Exec(registrationTable)
    if err != nil {
        log.Fatal("Failed to create registrations table:", err)
    }
}
