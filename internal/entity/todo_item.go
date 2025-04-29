package entity

import "time"

type TodoItem struct {
	ID          string    `json:"id" pgdb:"id"`
	FamilyID    string    `json:"family_id" pgdb:"family_id"`
	Title       string    `json:"title" pgdb:"title"`
	Description string    `json:"description" pgdb:"description"`
	Status      string    `json:"status" pgdb:"status"`
	Deadline    time.Time `json:"deadline" pgdb:"deadline"`
	AssignedTo  string    `json:"assigned_to" pgdb:"assigned_to"`
	CreatedBy   string    `json:"created_by" pgdb:"created_by"`
	IsArchived  bool      `json:"is_archived" pgdb:"is_archived"`
	CreatedAt   time.Time `json:"created_at" pgdb:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" pgdb:"updated_at"`
}
