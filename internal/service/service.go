package service

import (
	"context"
	"log/slog"

	"family-flow-app/internal/entity"
	"family-flow-app/internal/repo"
)

type UserCreateInput struct {
	Name     string
	Email    string
	Password string
	Role     string
}

type UserGetByIdInput struct {
	Id string
}

type UserGetByEmailInput struct {
	Email string
}

type AuthInput struct {
	Email    string
	Password string
}

type User interface {
	Create(ctx context.Context, log *slog.Logger, input UserCreateInput) (string, error)
	Login(ctx context.Context, log *slog.Logger, input AuthInput) (string, error)
	GetById(ctx context.Context, log *slog.Logger, input UserGetByIdInput) (entity.User, error)
	GetByEmail(ctx context.Context, log *slog.Logger, input UserGetByEmailInput) (
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
