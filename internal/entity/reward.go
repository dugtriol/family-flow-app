package entity

import (
	"time"
)

type Reward struct {
	ID          string    `json:"id" pgdb:"id"`
	FamilyID    string    `json:"family_id" pgdb:"family_id"`
	Title       string    `json:"title" pgdb:"title"`
	Description string    `json:"description" pgdb:"description"`
	Cost        int       `json:"cost" pgdb:"cost"`
	CreatedAt   time.Time `json:"created_at" pgdb:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" pgdb:"updated_at"`
}
