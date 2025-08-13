package models

import "time"

type Registration struct {
    ID                     int       `json:"id"`
    Fullname               string    `json:"fullname"`
    Dob                    string    `json:"dob"` // kept as string for simplicity (yyyy-mm-dd)
    Gender                 string    `json:"gender"`
    Email                  string    `json:"email"`
    Phone                  string    `json:"phone"`
    StateOfOrigin          string    `json:"state_of_origin"`
    StateOfResidence       string    `json:"state_of_residence"`
    Education              string    `json:"education"`
    PreviousOffice         string    `json:"previous_office"`
    CardCarryingMember     bool      `json:"card_carrying_member"`
    PartyMembershipDoc     string    `json:"party_membership_doc"`
    CV                     string    `json:"cv"`
    Motivation             string    `json:"motivation"`
    PoliticalUnderstanding string    `json:"political_understanding"`
    AssistanceNeeded       string    `json:"assistance_needed"`
    Availability           string    `json:"availability"`
    PreferredCommunication string    `json:"preferred_communication"`
    Consent                bool      `json:"consent"`
    CreatedAt              time.Time `json:"created_at"`
}
