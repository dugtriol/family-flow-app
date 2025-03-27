package service

import (
	"context"
	"log/slog"
	"time"

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

type TaskCreateInput struct {
	Title       string
	Description string
	Status      string
	Deadline    time.Time
	AssignedTo  string
	CreatedBy   string
	Reward      int
}

type TaskGetByIdInput struct {
	Id string
}

type TaskGetByAssignedToInput struct {
	AssignedTo string
}

type TaskGetByCreatedByInput struct {
	CreatedBy string
}

type TaskUpdateInput struct {
	Id          string
	Title       string
	Description string
	Status      string
	Deadline    time.Time
	AssignedTo  string
	Reward      int
}

type TaskDeleteInput struct {
	Id string
}

type TaskGetByStatusInput struct {
	Status string
}

type TaskCompleteInput struct {
	Id     string
	UserId string
}

type Task interface {
	Create(ctx context.Context, log *slog.Logger, input TaskCreateInput) (string, error)
	GetById(ctx context.Context, log *slog.Logger, input TaskGetByIdInput) (entity.Task, error)
	GetByAssignedTo(ctx context.Context, log *slog.Logger, input TaskGetByAssignedToInput) ([]entity.Task, error)
	GetByCreatedBy(ctx context.Context, log *slog.Logger, input TaskGetByCreatedByInput) ([]entity.Task, error)
	GetOverdueTasks(ctx context.Context, log *slog.Logger) ([]entity.Task, error)
	GetByStatus(ctx context.Context, log *slog.Logger, input TaskGetByStatusInput) (
		[]entity.Task, error,
	)
	Update(ctx context.Context, log *slog.Logger, input TaskUpdateInput) error
	Delete(ctx context.Context, log *slog.Logger, input TaskDeleteInput) error
	Complete(ctx context.Context, log *slog.Logger, input TaskCompleteInput) error
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
	FamilyId string
	UserId   string
}

type Family interface {
	CreateFamily(ctx context.Context, log *slog.Logger, input FamilyCreateInput) (string, error)
	GetFamilyByID(ctx context.Context, log *slog.Logger, id string) (entity.Family, error)
	AddMemberToFamily(ctx context.Context, log *slog.Logger, input AddMemberToFamilyInput) error
}

type Services struct {
	User   User
	Task   Task
	Email  Email
	Family Family
}

type ServicesDependencies struct {
	Rds    *redis.Redis
	Repos  *repo.Repositories
	Config *config.Config
}

func NewServices(dep ServicesDependencies) *Services {
	return &Services{
		User:   NewUserService(dep.Repos.User),
		Task:   NewTaskService(dep.Repos.Task),
		Email:  NewEmailService(dep.Rds, dep.Config.Email),
		Family: NewFamilyService(dep.Repos.Family, dep.Repos.User),
	}
}
