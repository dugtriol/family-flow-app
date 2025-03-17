package entity

import "time"

type Family struct {
	Id        string    `pgdb:"id"`
	Name      string    `pgdb:"name"`
	CreatedAt time.Time `pgdb:"created_at"`
}
