package entity

import `time`

type Notification struct {
	ID        string    `db:"id" json:"id"`
	UserID    string    `db:"user_id" json:"user_id"`
	Title     string    `db:"title" json:"title"`
	Body      string    `db:"body" json:"body"`
	Data      string    `db:"data" json:"data"`
	IsRead    bool      `db:"is_read" json:"is_read"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
