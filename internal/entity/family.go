package entity

import (
	`database/sql`
	`time`
)

type Family struct {
	Id        string         `pgdb:"id" json:"id"`
	Name      string         `pgdb:"name" json:"name"`
	CreatedAt time.Time      `pgdb:"created_at" json:"created_at"`
	Photo     sql.NullString `pgdb:"photo" json:"photo" swaggerignore:"true"`
}
