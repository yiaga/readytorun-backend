package models

import (
    "time"
)

// Registration represents a user registration.
type Registration struct {
    ID                     int64          `json:"id"`
    Fullname               string         `json:"fullname"`
    Dob                    *string `json:"dob,omitempty"`
    Gender                 *string `json:"gender,omitempty"`
    Email                  string         `json:"email"`
    Phone                  *string `json:"phone,omitempty"`
    StateOfOrigin          *string `json:"stateOfOrigin,omitempty"`
    StateOfResidence       *string `json:"stateOfResidence,omitempty"`
    Education              *string `json:"education,omitempty"`
    PreviousOffice         *string `json:"previousOffice,omitempty"`
    InterestedOffice       *string `json:"interestedOffice,omitempty"`
    PreviousContest        *string `json:"previousContest,omitempty"`
    CardCarryingMember     bool           `json:"partyMember"`
    PartyMembershipDocLink string         `json:"partyMembershipDocLink,omitempty"`
    Motivation             *string `json:"motivation,omitempty"`
    PoliticalUnderstanding *string `json:"politicalUnderstanding,omitempty"`
    AssistanceNeeded       []string       `json:"assistanceNeeded,omitempty"` 
    OtherSupport           *string `json:"otherSupport,omitempty"`
    PreferredCommunication *string `json:"preferred_communication,omitempty"`
    Consent                bool           `json:"consent"`
    CreatedAt              time.Time      `json:"createdAt"`
}