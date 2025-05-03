package service

import (
	"context"
	"log/slog"

	"family-flow-app/config"
	"family-flow-app/internal/entity"
	"family-flow-app/internal/repo"
	"family-flow-app/pkg/redis"
	firebase "firebase.google.com/go/v4"
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
	Update(ctx context.Context, log *slog.Logger, input UpdateUserInput) error
	UpdatePassword(ctx context.Context, log *slog.Logger, email, password string) error
	ResetFamilyID(ctx context.Context, log *slog.Logger, id string) error
	ExistsByEmail(ctx context.Context, log *slog.Logger, email string) (bool, error)
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
	Role      string
}

type Family interface {
	Create(ctx context.Context, log *slog.Logger, input FamilyCreateInput) (string, error)
	GetFamilyByUserID(ctx context.Context, log *slog.Logger, id string) (entity.Family, error)
	AddMember(ctx context.Context, log *slog.Logger, input AddMemberToFamilyInput) error
	GetByFamilyID(ctx context.Context, log *slog.Logger, familyId string) ([]entity.User, error)
	GetByID(ctx context.Context, log *slog.Logger, id string) (entity.Family, error)
}

type WishlistItem interface {
	Create(ctx context.Context, log *slog.Logger, input WishlistCreateInput) (string, error)
	Delete(ctx context.Context, log *slog.Logger, id string) error
	Update(ctx context.Context, log *slog.Logger, input WishlistUpdateInput) error
	GetByID(ctx context.Context, log *slog.Logger, id string) ([]entity.WishlistItem, error)
	UpdateReservedBy(ctx context.Context, log *slog.Logger, input WishlistUpdateReservedByInput) error
	GetArchivedByUserID(
		ctx context.Context, log *slog.Logger, userID string,
	) ([]entity.WishlistItem, error)
	CancelUpdateReservedBy(
		ctx context.Context, log *slog.Logger, input WishlistCancelUpdateReservedByInput,
	) error
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
	UpdateReservedBy(
		ctx context.Context, log *slog.Logger, input ShoppingUpdateReservedByInput,
	) error
	UpdateBuyerId(
		ctx context.Context, log *slog.Logger, input ShoppingUpdateBuyerIdInput,
	) error
	GetArchivedByUserID(
		ctx context.Context, log *slog.Logger, userID string,
	) ([]entity.ShoppingItem, error)
	CancelUpdateReservedBy(
		ctx context.Context, log *slog.Logger, input ShoppingCancelUpdateReservedByInput,
	) error
}

type TodoItem interface {
	Create(ctx context.Context, log *slog.Logger, input TodoCreateInput) (string, error)
	Delete(ctx context.Context, log *slog.Logger, id string) error
	Update(ctx context.Context, log *slog.Logger, input TodoUpdateInput) error
	GetByAssignedTo(ctx context.Context, log *slog.Logger, assignedTo string) ([]entity.TodoItem, error)
	GetByCreatedBy(ctx context.Context, log *slog.Logger, createdBy string) ([]entity.TodoItem, error)
}

type Notification interface {
	SendNotification(ctx context.Context, log *slog.Logger, input NotificationCreateInput) error
	SaveToken(ctx context.Context, userID, token string) error
	GetNotificationsByUserID(
		ctx context.Context, log *slog.Logger, userID string,
	) ([]entity.Notification, error)
}

type Chats interface {
	CreateChat(ctx context.Context, log *slog.Logger, input CreateChatInput) (string, error)
	AddParticipant(ctx context.Context, log *slog.Logger, input AddParticipantInput) error
	CreateChatWithParticipants(ctx context.Context, log *slog.Logger, input CreateChatWithParticipantsInput) (
		string, error,
	)
	CreateMessage(ctx context.Context, log *slog.Logger, input CreateMessageInput) (string, error)
	GetParticipants(ctx context.Context, log *slog.Logger, chatID string) ([]string, error)
	GetMessagesByChatID(ctx context.Context, log *slog.Logger, chatID string) ([]entity.Message, error)
}

type Services struct {
	User         User
	Email        Email
	Family       Family
	WishlistItem WishlistItem
	ShoppingItem ShoppingItem
	TodoItem     TodoItem
	Notification Notification
	Chats        Chats
}

type ServicesDependencies struct {
	Rds    *redis.Redis
	Repos  *repo.Repositories
	Config *config.Config
	App    *firebase.App
}

func NewServices(ctx context.Context, dep ServicesDependencies) *Services {
	return &Services{
		User:         NewUserService(dep.Repos.User),
		Email:        NewEmailService(dep.Rds, dep.Config.Email),
		Family:       NewFamilyService(dep.Repos.Family, dep.Repos.User),
		WishlistItem: NewWishlistService(dep.Repos.WishlistItem),
		ShoppingItem: NewShoppingService(dep.Repos.ShoppingItem),
		TodoItem:     NewTodoService(dep.Repos.TodosItem),
		Notification: NewNotificationService(ctx, dep.Rds, dep.App, dep.Repos.Notification),
		Chats:        NewChatMessageService(dep.Repos.Chat, dep.Repos.Message),
	}
}
