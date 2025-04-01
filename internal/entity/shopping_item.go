package entity

import "time"

type ShoppingItem struct {
	ID          string    `json:"id" pgdb:"id"`
	FamilyID    string    `json:"family_id" pgdb:"family_id"`
	Title       string    `json:"title" pgdb:"title"`
	Description string    `json:"description" pgdb:"description"`
	Status      string    `json:"status" pgdb:"status"`
	Visibility  string    `json:"visibility" pgdb:"visibility"`
	CreatedBy   string    `json:"created_by" pgdb:"created_by"`
	CreatedAt   time.Time `json:"created_at" pgdb:"created_at"`
}
