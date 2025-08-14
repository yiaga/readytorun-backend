package models

type Registration struct {
    ID                     int    `json:"id"`
    Fullname               string `json:"fullname"`
    Dob                    string `json:"dateOfBirth"`
    Gender                 string `json:"gender"`
    Email                  string `json:"email"`
    Phone                  string `json:"phone"`
    StateOfOrigin          string `json:"stateOfOrigin"`
    StateOfResidence       string `json:"stateOfResidence"`
    Education              string `json:"education"`
    PreviousOffice         string `json:"previousOffice"`
    InterestedOffice       string `json:"interestedOffice"`
    PreviousContest        string `json:"previousContest"`
    CardCarryingMember     bool   `json:"partyMember"`
    PartyMembershipDocLink string `json:"partyMembershipDocLink"` // Changed to store Google Drive link
    CVLink                 string `json:"cvLink"`                 // Changed to store Google Drive link
    Motivation             string `json:"motivation"`
    PoliticalUnderstanding string `json:"politicalUnderstanding"`
    AssistanceNeeded       string `json:"assistanceNeeded"` // Stored as a JSON string
    OtherSupport           string `json:"otherSupport"`
    Availability           string `json:"availability"`     // Stored as a JSON string
    PreferredCommunication string `json:"communication"`
    Consent                bool   `json:"consent"`
    CreatedAt              string `json:"createdAt"` // Add this new field
}