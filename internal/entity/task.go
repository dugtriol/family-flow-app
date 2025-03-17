package entity

import "time"

type Task struct {
	Id          string    `pgdb:"id"`
	Title       string    `pgdb:"title"`
	Description string    `pgdb:"description"`
	Status      string    `pgdb:"status"`
	Deadline    time.Time `pgdb:"deadline"`
	AssignedTo  string    `pgdb:"assigned_to"`
	CreatedBy   string    `pgdb:"created_by"`
	Reward      int       `pgdb:"reward"`
}
