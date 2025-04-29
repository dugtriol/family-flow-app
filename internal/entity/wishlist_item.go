package entity

import (
	`database/sql`
	`time`
)

type WishlistItem struct {
	ID          string         `json:"id" pgdb:"id"`
	Name        string         `json:"name" pgdb:"name"`
	Description string         `json:"description" pgdb:"description"`
	Link        string         `json:"link" pgdb:"link"`
	Status      string         `json:"status" pgdb:"status"`
	CreatedBy   string         `json:"created_by" pgdb:"created_by"`
	ReservedBy  sql.NullString `json:"reserved_by" pgdb:"reserved_by"`
	IsArchived  bool           `json:"is_archived" pgdb:"is_archived"`
	CreatedAt   time.Time      `json:"created_at" pgdb:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at" pgdb:"updated_at"`
}
