package models

import (
	"time"

	"github.com/lib/pq"
)

type Volunteer struct {
	ID               int            `json:"id"`
	FullName        string         `json:"full_name"`
	Email            string         `json:"email"`
	Phone            *string        `json:"phone,omitempty"`
	Location         *string        `json:"location,omitempty"readytorun_user"`
	Skills           pq.StringArray `json:"skills" gorm:"type:text[]"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}
