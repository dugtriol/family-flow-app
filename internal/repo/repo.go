package repo

import (
	"context"

	"family-flow-app/internal/entity"
	"family-flow-app/internal/repo/pgdb"
	"family-flow-app/pkg/postgres"
)

type User interface {
	Create(ctx context.Context, input entity.User) (string, error)
	GetById(ctx context.Context, id string) (entity.User, error)
	GetByUsername(ctx context.Context, username string) (entity.User, error)
}

type Repositories struct {
	User
}

func NewRepositories(db *postgres.Database) *Repositories {
	return &Repositories{
		User: pgdb.NewUserRepo(db),
	}
}
