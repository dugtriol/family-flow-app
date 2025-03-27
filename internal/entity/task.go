package entity

import "time"

type Task struct {
	Id          string    `pgdb:"id" json:"id"`
	Title       string    `pgdb:"title" json:"title"`
	Description string    `pgdb:"description" json:"description"`
	Status      string    `pgdb:"status" json:"status"`
	Deadline    time.Time `pgdb:"deadline" json:"deadline"`
	AssignedTo  string    `pgdb:"assigned_to" json:"assigned_to"`
	CreatedBy   string    `pgdb:"created_by" json:"created_by"`
	Reward      int       `pgdb:"reward" json:"reward"`
}
