package models

import "time"

type Contact struct {
    ID        int       `json:"id"`
    Fullname  string    `json:"fullname"`
    Email     string    `json:"email"`
    Phone     string    `json:"phone"`
    Subject   string    `json:"subject"`
    Message   string    `json:"message"`
    CreatedAt time.Time `json:"created_at"`
}
