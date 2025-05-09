package entity

import (
	`database/sql`
	`time`
)

type ShoppingItem struct {
	ID          string         `json:"id" pgdb:"id"`
	FamilyID    string         `json:"family_id" pgdb:"family_id"`
	Title       string         `json:"title" pgdb:"title"`
	Description string         `json:"description" pgdb:"description"`
	Status      string         `json:"status" pgdb:"status"`
	Visibility  string         `json:"visibility" pgdb:"visibility"`
	CreatedBy   string         `json:"created_by" pgdb:"created_by"`
	ReservedBy  sql.NullString `json:"reserved_by" pgdb:"reserved_by" swaggerignore:"true"`
	BuyerId     sql.NullString `json:"buyer_id" pgdb:"buyer_id" swaggerignore:"true"`
	IsArchived  bool           `json:"is_archived" pgdb:"is_archived"`
	CreatedAt   time.Time      `json:"created_at" pgdb:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at" pgdb:"updated_at"`
}
