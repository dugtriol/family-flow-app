package entity

import `time`

type ChatParticipant struct {
	ID       string    `json:"id" db:"id"`
	ChatId   string    `json:"chat_id" db:"chat_id"`
	UserId   string    `json:"user_id" db:"user_id"`
	JoinedAt time.Time `json:"joined_at" db:"joined_at"`
}
