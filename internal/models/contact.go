package models

import "time"

type Contact struct {
    ID        int64     `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    Message   string    `json:"message"`
    Subject   string    `json:"subject"`
    CreatedAt time.Time `json:"created_at"`
}