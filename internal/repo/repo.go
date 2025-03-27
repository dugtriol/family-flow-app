package repo

import (
	"context"

	"family-flow-app/internal/entity"
	"family-flow-app/internal/repo/pgdb"
	"family-flow-app/pkg/postgres"
)

type User interface {
	Create(ctx context.Context, user entity.User) (string, error)
	GetByID(ctx context.Context, id string) (entity.User, error)
	GetByEmail(ctx context.Context, email string) (entity.User, error)
	UpdateFamilyID(ctx context.Context, userID, familyID string) error
	GetByFamilyID(ctx context.Context, familyID string) ([]entity.User, error)
	Update(ctx context.Context, user entity.User) error
}

type Task interface {
	Create(ctx context.Context, task entity.Task) (string, error)
	GetByID(ctx context.Context, id string) (entity.Task, error)
	GetByAssignedTo(ctx context.Context, assignedTo string) ([]entity.Task, error)
	GetByCreatedBy(ctx context.Context, createdBy string) ([]entity.Task, error)
	GetByStatus(ctx context.Context, status string) ([]entity.Task, error)
	GetOverdueTasks(ctx context.Context) ([]entity.Task, error)
	Update(ctx context.Context, task entity.Task) error
	Delete(ctx context.Context, id string) error
	Complete(ctx context.Context, id string, userId string) error
}

type Family interface {
	Create(ctx context.Context, family entity.Family) (string, error)
	GetByID(ctx context.Context, id string) (entity.Family, error)
}

type Repositories struct {
	User
	Task
	Family
}

func NewRepositories(db *postgres.Database) *Repositories {
	return &Repositories{
		User:   pgdb.NewUserRepo(db),
		Task:   pgdb.NewTaskRepo(db),
		Family: pgdb.NewFamilyRepo(db),
	}
}
