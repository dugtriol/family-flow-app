package repo

import (
	"context"
	"log/slog"
	"time"

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
	UpdatePassword(ctx context.Context, email, password string) error
	ResetFamilyID(ctx context.Context, id string) error
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	UpdateRole(ctx context.Context, email, role string) error
	UpdateLocation(ctx context.Context, userID string, latitude, longitude float64) error
	UpdatePoint(ctx context.Context, userID string, point int) error
}

type Family interface {
	Create(ctx context.Context, family entity.Family) (string, error)
	GetByID(ctx context.Context, id string) (entity.Family, error)
	UpdatePhoto(ctx context.Context, familyId, photoURL string) error
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
	UpdateReservedBy(
		ctx context.Context, log *slog.Logger, id string, reservedBy string, updatedAt time.Time,
	) error
	UpdateBuyerId(
		ctx context.Context, log *slog.Logger, id string, buyerId string, updatedAt time.Time,
	) error
	GetArchivedByUserID(
		ctx context.Context, log *slog.Logger, userID string,
	) ([]entity.ShoppingItem, error)
	CancelUpdateReservedBy(
		ctx context.Context, log *slog.Logger, id string, updatedAt time.Time,
	) error
	GetByID(
		ctx context.Context, log *slog.Logger, id string,
	) (entity.ShoppingItem, error)
}

type TodosItem interface {
	Create(ctx context.Context, log *slog.Logger, item entity.TodoItem) (string, error)
	Delete(ctx context.Context, log *slog.Logger, id string) error
	Update(ctx context.Context, log *slog.Logger, item entity.TodoItem) error
	GetByAssignedTo(ctx context.Context, log *slog.Logger, assignedTo string) (
		[]entity.TodoItem, error,
	)
	GetByCreatedBy(ctx context.Context, log *slog.Logger, createdBy string) ([]entity.TodoItem, error)
	GetByID(ctx context.Context, log *slog.Logger, id string) (entity.TodoItem, error)
}

type WishlistItem interface {
	Create(ctx context.Context, log *slog.Logger, item entity.WishlistItem) (string, error)
	Delete(ctx context.Context, log *slog.Logger, id string) error
	Update(ctx context.Context, log *slog.Logger, item entity.WishlistItem) error
	GetByUserID(ctx context.Context, log *slog.Logger, userID string) ([]entity.WishlistItem, error)
	UpdateReservedBy(ctx context.Context, log *slog.Logger, id, reservedBy string) error
	GetArchivedByUserID(ctx context.Context, log *slog.Logger, userID string) (
		[]entity.WishlistItem, error,
	)
	CancelUpdateReservedBy(ctx context.Context, log *slog.Logger, id string) error
	GetByID(ctx context.Context, log *slog.Logger, id string) (entity.WishlistItem, error)
}

type Notification interface {
	Create(ctx context.Context, notification entity.Notification) (string, error)
	GetByUserID(ctx context.Context, userID string) ([]entity.Notification, error)
}

type Chat interface {
	Create(ctx context.Context, chat entity.Chat) (string, error)
	GetByID(ctx context.Context, id string) (entity.Chat, error)
	GetAll(ctx context.Context) ([]entity.Chat, error)
	AddParticipant(ctx context.Context, chatID, userID string) error
	GetParticipants(ctx context.Context, chatID string) ([]string, error)
	GetChatsWithParticipants(ctx context.Context, userID string) ([]entity.Chat, error)
}

type Message interface {
	Create(ctx context.Context, message entity.Message) (entity.Message, error)
	GetByChatID(ctx context.Context, chatID string) ([]entity.Message, error)
	GetLastMessageByChatID(ctx context.Context, chatID string) (entity.Message, error)
}

type Rewards interface {
	Create(ctx context.Context, reward entity.Reward) (string, error)
	GetByFamilyID(ctx context.Context, familyID string) ([]entity.Reward, error)
	AddPoints(ctx context.Context, userID string, points int) error
	SubtractPoints(ctx context.Context, userID string, points int) error
	GetPoints(ctx context.Context, userID string) (int, error)
	Redeem(ctx context.Context, userID, rewardID string) error
	GetRedemptionsByUserID(ctx context.Context, userID string) ([]entity.RewardRedemption, error)
	GetByID(ctx context.Context, id string) (entity.Reward, error)
	Update(ctx context.Context, reward entity.Reward) error
}

type Diary interface {
	Create(ctx context.Context, log *slog.Logger, item entity.DiaryItem) (string, error)
	GetByUserID(ctx context.Context, log *slog.Logger, userID string) ([]entity.DiaryItem, error)
	Update(ctx context.Context, log *slog.Logger, item entity.DiaryItem) error
	Delete(ctx context.Context, log *slog.Logger, id string) error
}

type NotificationToken interface {
	SaveOrUpdate(ctx context.Context, log *slog.Logger, token entity.NotificationToken) error
	GetByUserID(ctx context.Context, log *slog.Logger, userID string) (*entity.NotificationToken, error)
}

type Repositories struct {
	User
	Family
	ShoppingItem
	TodosItem
	WishlistItem
	Notification
	Chat
	Message
	Rewards
	Diary
	NotificationToken
}

func NewRepositories(db *postgres.Database) *Repositories {
	return &Repositories{
		User:              pgdb.NewUserRepo(db),
		Family:            pgdb.NewFamilyRepo(db),
		ShoppingItem:      pgdb.NewShoppingRepo(db),
		TodosItem:         pgdb.NewTodoRepo(db),
		WishlistItem:      pgdb.NewWishlistRepo(db),
		Notification:      pgdb.NewNotificationsRepo(db),
		Chat:              pgdb.NewChatsRepo(db),
		Message:           pgdb.NewMessagesRepo(db),
		Rewards:           pgdb.NewRewardsRepo(db),
		Diary:             pgdb.NewDiaryRepo(db),
		NotificationToken: pgdb.NewNotificationTokenRepo(db),
	}
}
