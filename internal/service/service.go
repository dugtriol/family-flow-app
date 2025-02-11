package service

import (
	"context"
	"log/slog"

	"family-flow-app/internal/entity"
	"family-flow-app/internal/repo"
)

type UserCreateInput struct {
	Username  string
	FirstName string
	LastName  string
}

type UserGetByIdInput struct {
	Id string
}

type UserGetByUsernameInput struct {
	Username string
}

type User interface {
	Create(ctx context.Context, log *slog.Logger, input UserCreateInput) (string, error)
	GetById(ctx context.Context, log *slog.Logger, input UserGetByIdInput) (entity.User, error)
	GetByUsername(ctx context.Context, log *slog.Logger, input UserGetByUsernameInput) (
		entity.User, error,
	)
}

type Services struct {
	User User
}

type ServicesDependencies struct {
	Repos *repo.Repositories
}

func NewServices(dep ServicesDependencies) *Services {
	return &Services{
		User: NewUserService(dep.Repos.User),
	}
}
