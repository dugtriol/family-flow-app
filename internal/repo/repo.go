package repo

import (
	"context"
	"log/slog"

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

type Family interface {
	Create(ctx context.Context, family entity.Family) (string, error)
	GetByID(ctx context.Context, id string) (entity.Family, error)
}

type ShoppingItem interface {
	Create(ctx context.Context, log *slog.Logger, item entity.ShoppingItem) (string, error)
	Delete(ctx context.Context, log *slog.Logger, id string) error
	Update(ctx context.Context, log *slog.Logger, item entity.ShoppingItem) error
	GetPublicByFamilyID(
		ctx context.Context, log *slog.Logger, familyID string,
	) ([]entity.ShoppingItem, error)
	GetPrivateByCreatedBy(
		ctx context.Context, log *slog.Logger, createdBy string,
	) ([]entity.ShoppingItem, error)
}

type TodosItem interface {
	Create(ctx context.Context, log *slog.Logger, item entity.TodoItem) (string, error)
	Delete(ctx context.Context, log *slog.Logger, id string) error
	Update(ctx context.Context, log *slog.Logger, item entity.TodoItem) error
	GetByAssignedTo(ctx context.Context, log *slog.Logger, assignedTo string) (
		[]entity.TodoItem, error,
	)
	GetByCreatedBy(ctx context.Context, log *slog.Logger, createdBy string) ([]entity.TodoItem, error)
}

type WishlistItem interface {
	Create(ctx context.Context, log *slog.Logger, item entity.WishlistItem) (string, error)
	Delete(ctx context.Context, log *slog.Logger, id string) error
	Update(ctx context.Context, log *slog.Logger, item entity.WishlistItem) error
	GetByUserID(ctx context.Context, log *slog.Logger, userID string) ([]entity.WishlistItem, error)
}

type Repositories struct {
	User
	Family
	ShoppingItem
	TodosItem
	WishlistItem
}

func NewRepositories(db *postgres.Database) *Repositories {
	return &Repositories{
		User:         pgdb.NewUserRepo(db),
		Family:       pgdb.NewFamilyRepo(db),
		ShoppingItem: pgdb.NewShoppingRepo(db),
		TodosItem:    pgdb.NewTodoRepo(db),
		WishlistItem: pgdb.NewWishlistRepo(db),
	}
}
