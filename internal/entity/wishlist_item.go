package entity

import "time"

type WishlistItem struct {
	ID          string    `json:"id" pgdb:"id"`
	Name        string    `json:"name" pgdb:"name"`
	Description string    `json:"description" pgdb:"description"`
	Link        string    `json:"link" pgdb:"link"`
	Status      string    `json:"status" pgdb:"status"`
	IsReserved  bool      `json:"is_reserved" pgdb:"is_reserved"`
	CreatedBy   string    `json:"created_by" pgdb:"created_by"`
	CreatedAt   time.Time `json:"created_at" pgdb:"created_at"`
}
