package entity

import (
	"time"
)

type RewardRedemption struct {
	ID         string    `json:"id" pgdb:"id"`
	UserID     string    `json:"user_id" pgdb:"user_id"`
	RewardID   string    `json:"reward_id" pgdb:"reward_id"`
	RedeemedAt time.Time `json:"redeemed_at" pgdb:"redeemed_at"`
	Reward     Reward    `json:"reward" pgdb:"reward"`
}
