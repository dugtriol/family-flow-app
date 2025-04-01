package service

import (
	"context"
	"log/slog"

	"family-flow-app/config"
	"family-flow-app/internal/entity"
	"family-flow-app/internal/repo"
	"family-flow-app/pkg/redis"
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
	GetById(ctx context.Context, log *slog.Logger, id string) (entity.User, error)
	GetByEmail(ctx context.Context, log *slog.Logger, input UserGetByEmailInput) (
		entity.User, error,
	)
}

type InputSendInvite struct {
	To         []string
	From       string
	FromName   string
	FamilyName string
}

type Email interface {
	SendCode(ctx context.Context, to []string) error
	CompareCode(ctx context.Context, email, code string) (bool, error)
	GetAllKeys(ctx context.Context) ([]string, error)
	SendInvite(ctx context.Context, invite InputSendInvite) error
}

type FamilyCreateInput struct {
	Name          string
	CreatorUserId string
}

type AddMemberToFamilyInput struct {
	FamilyId  string
	UserEmail string
}

type Family interface {
	Create(ctx context.Context, log *slog.Logger, input FamilyCreateInput) (string, error)
	GetFamilyByUserID(ctx context.Context, log *slog.Logger, id string) (entity.Family, error)
	AddMember(ctx context.Context, log *slog.Logger, input AddMemberToFamilyInput) error
	GetByFamilyID(ctx context.Context, log *slog.Logger, familyId string) ([]entity.User, error)
}

type WishlistItem interface {
	Create(ctx context.Context, log *slog.Logger, input WishlistCreateInput) (string, error)
	Delete(ctx context.Context, log *slog.Logger, id string) error
	Update(ctx context.Context, log *slog.Logger, input WishlistUpdateInput) error
	GetByID(ctx context.Context, log *slog.Logger, id string) ([]entity.WishlistItem, error)
}

type ShoppingItem interface {
	Create(ctx context.Context, log *slog.Logger, input ShoppingCreateInput) (string, error)
	Delete(ctx context.Context, log *slog.Logger, id string) error
	Update(ctx context.Context, log *slog.Logger, input ShoppingUpdateInput) error
	GetPublicByFamilyID(
		ctx context.Context, log *slog.Logger, familyID string,
	) ([]entity.ShoppingItem, error)
	GetPrivateByCreatedBy(
		ctx context.Context, log *slog.Logger, createdBy string,
	) ([]entity.ShoppingItem, error)
}

type TodoItem interface {
	Create(ctx context.Context, log *slog.Logger, input TodoCreateInput) (string, error)
	Delete(ctx context.Context, log *slog.Logger, id string) error
	Update(ctx context.Context, log *slog.Logger, input TodoUpdateInput) error
	GetByAssignedTo(ctx context.Context, log *slog.Logger, assignedTo string) ([]entity.TodoItem, error)
	GetByCreatedBy(ctx context.Context, log *slog.Logger, createdBy string) ([]entity.TodoItem, error)
}

type Services struct {
	User         User
	Email        Email
	Family       Family
	WishlistItem WishlistItem
	ShoppingItem ShoppingItem
	TodoItem     TodoItem
}

type ServicesDependencies struct {
	Rds    *redis.Redis
	Repos  *repo.Repositories
	Config *config.Config
}

func NewServices(dep ServicesDependencies) *Services {
	return &Services{
		User:         NewUserService(dep.Repos.User),
		Email:        NewEmailService(dep.Rds, dep.Config.Email),
		Family:       NewFamilyService(dep.Repos.Family, dep.Repos.User),
		WishlistItem: NewWishlistService(dep.Repos.WishlistItem),
		ShoppingItem: NewShoppingService(dep.Repos.ShoppingItem),
		TodoItem:     NewTodoService(dep.Repos.TodosItem),
	}
}
