package service

import (
	"context"
	"log/slog"
	"time"

	"family-flow-app/internal/entity"
	"family-flow-app/internal/repo"
)

type TodoService struct {
	todoRepo repo.TodosItem
	userRepo repo.User
}

func NewTodoService(todoRepo repo.TodosItem, userRepo repo.User) *TodoService {
	return &TodoService{todoRepo: todoRepo, userRepo: userRepo}
}

type TodoCreateInput struct {
	FamilyId    string
	Title       string
	Description string
	Deadline    time.Time
	AssignedTo  string
	CreatedBy   string
	Point       int
}

func (t *TodoService) Create(ctx context.Context, log *slog.Logger, input TodoCreateInput) (string, error) {
	log.Info("Service - TodoService - Create")

	item := entity.TodoItem{
		FamilyID:    input.FamilyId,
		Title:       input.Title,
		Description: input.Description,
		Deadline:    input.Deadline,
		AssignedTo:  input.AssignedTo,
		CreatedBy:   input.CreatedBy,
		Point:       input.Point,
	}

	id, err := t.todoRepo.Create(ctx, log, item)
	if err != nil {
		log.Error("Service - TodoService - Create: %v", err)
		return "", err
	}

	return id, nil
}

func (t *TodoService) Delete(ctx context.Context, log *slog.Logger, id string) error {
	log.Info("Service - TodoService - Delete")

	err := t.todoRepo.Delete(ctx, log, id)
	if err != nil {
		log.Error("Service - TodoService - Delete: %v", err)
		return err
	}

	return nil
}

type TodoUpdateInput struct {
	ID          string
	Title       string
	Description string
	Status      string
	Deadline    time.Time
	AssignedTo  string
	Point       int
}

func (t *TodoService) Update(ctx context.Context, log *slog.Logger, input TodoUpdateInput) error {
	log.Info("Service - TodoService - Update")

	err := t.todoRepo.Update(
		ctx, log, entity.TodoItem{
			ID:          input.ID,
			Title:       input.Title,
			Description: input.Description,
			Status:      input.Status,
			Deadline:    input.Deadline,
			AssignedTo:  input.AssignedTo,
			Point:       input.Point,
		},
	)

	if err != nil {
		log.Error("Service - TodoService - Update: %v", err)
		return err
	}

	err = t.userRepo.UpdatePoint(ctx, input.AssignedTo, input.Point)
	if err != nil {
		log.Error("Service - TodoService - UpdatePoint: %v", err)
		return err
	}

	return nil
}

// get by assigned to
func (t *TodoService) GetByAssignedTo(ctx context.Context, log *slog.Logger, assignedTo string) (
	[]entity.TodoItem, error,
) {
	log.Info("Service - TodoService - GetByAssignedTo")

	items, err := t.todoRepo.GetByAssignedTo(ctx, log, assignedTo)
	if err != nil {
		log.Error("Service - TodoService - GetByAssignedTo: %v", err)
		return nil, err
	}

	return items, nil
}

// get by created by
func (t *TodoService) GetByCreatedBy(ctx context.Context, log *slog.Logger, createdBy string) (
	[]entity.TodoItem, error,
) {
	log.Info("Service - TodoService - GetByCreatedBy")

	items, err := t.todoRepo.GetByCreatedBy(ctx, log, createdBy)
	if err != nil {
		log.Error("Service - TodoService - GetByCreatedBy: %v", err)
		return nil, err
	}

	return items, nil
}

// get by id
func (t *TodoService) GetByID(ctx context.Context, log *slog.Logger, id string) (entity.TodoItem, error) {
	log.Info("Service - TodoService - GetByID")

	item, err := t.todoRepo.GetByID(ctx, log, id)
	if err != nil {
		log.Error("Service - TodoService - GetByID: %v", err)
		return entity.TodoItem{}, err
	}

	return item, nil
}
