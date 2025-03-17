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
}

type Repositories struct {
	User
}

func NewRepositories(db *postgres.Database) *Repositories {
	return &Repositories{
		User: pgdb.NewUserRepo(db),
	}
}
