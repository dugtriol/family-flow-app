package entity

import (
	"time"
)

type Chat struct {
	ID           string            `json:"id" db:"id"`
	Name         string            `json:"name" db:"name"`
	CreatedAt    time.Time         `json:"created_at" db:"created_at"`
	Participants []ChatParticipant `json:"participants" db:"participants"`
}
