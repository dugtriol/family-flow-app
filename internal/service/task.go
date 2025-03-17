package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"family-flow-app/internal/entity"
	"family-flow-app/internal/repo"
	"family-flow-app/internal/repo/repoerrs"
)

type TaskService struct {
	taskRepo repo.Task
}

func NewTaskService(taskRepo repo.Task) *TaskService {
	return &TaskService{taskRepo: taskRepo}
}

func (t *TaskService) Create(ctx context.Context, log *slog.Logger, input TaskCreateInput) (string, error) {
	log.Info("Service - TaskService - Create")

	task := entity.Task{
		Title:       input.Title,
		Description: input.Description,
		Status:      input.Status,
		Deadline:    input.Deadline,
		AssignedTo:  input.AssignedTo,
		CreatedBy:   input.CreatedBy,
		Reward:      input.Reward,
	}

	id, err := t.taskRepo.Create(ctx, task)
	if err != nil {
		if errors.Is(err, repoerrs.ErrAlreadyExists) {
			return "", ErrTaskAlreadyExists
		}
		log.Error(fmt.Sprintf("Service - TaskService - Create: %v", err))
		return "", ErrCannotCreateTask
	}

	log.Info(fmt.Sprintf("Service - TaskService - taskRepo.Create - id: %s", id))
	return id, nil
}

func (t *TaskService) GetById(ctx context.Context, log *slog.Logger, input TaskGetByIdInput) (entity.Task, error) {
	log.Info("Service - TaskService - GetById")

	task, err := t.taskRepo.GetByID(ctx, input.Id)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return entity.Task{}, ErrTaskNotFound
		}
		log.Error(fmt.Sprintf("Service - TaskService - GetById: %v", err))
		return entity.Task{}, ErrCannotGetTask
	}

	return task, nil
}

func (t *TaskService) GetByAssignedTo(
	ctx context.Context, log *slog.Logger, input TaskGetByAssignedToInput,
) ([]entity.Task, error) {
	log.Info("Service - TaskService - GetByAssignedTo")

	tasks, err := t.taskRepo.GetByAssignedTo(ctx, input.AssignedTo)
	if err != nil {
		log.Error(fmt.Sprintf("Service - TaskService - GetByAssignedTo: %v", err))
		return nil, ErrCannotGetTasks
	}

	return tasks, nil
}

func (t *TaskService) GetByCreatedBy(
	ctx context.Context, log *slog.Logger, input TaskGetByCreatedByInput,
) ([]entity.Task, error) {
	log.Info("Service - TaskService - GetByCreatedBy")

	tasks, err := t.taskRepo.GetByCreatedBy(ctx, input.CreatedBy)
	if err != nil {
		log.Error(fmt.Sprintf("Service - TaskService - GetByCreatedBy: %v", err))
		return nil, ErrCannotGetTasks
	}

	return tasks, nil
}

func (t *TaskService) Update(ctx context.Context, log *slog.Logger, input TaskUpdateInput) error {
	log.Info("Service - TaskService - Update")

	task := entity.Task{
		Id:          input.Id,
		Title:       input.Title,
		Description: input.Description,
		Status:      input.Status,
		Deadline:    input.Deadline,
		AssignedTo:  input.AssignedTo,
		Reward:      input.Reward,
	}

	err := t.taskRepo.Update(ctx, task)
	if err != nil {
		log.Error(fmt.Sprintf("Service - TaskService - Update: %v", err))
		return ErrCannotUpdateTask
	}

	return nil
}

func (t *TaskService) Delete(ctx context.Context, log *slog.Logger, input TaskDeleteInput) error {
	log.Info("Service - TaskService - Delete")

	err := t.taskRepo.Delete(ctx, input.Id)
	if err != nil {
		log.Error(fmt.Sprintf("Service - TaskService - Delete: %v", err))
		return ErrCannotDeleteTask
	}

	return nil
}

func (t *TaskService) GetByStatus(ctx context.Context, log *slog.Logger, input TaskGetByStatusInput) (
	[]entity.Task, error,
) {
	log.Info("Service - TaskService - GetByStatus")

	tasks, err := t.taskRepo.GetByStatus(ctx, input.Status)
	if err != nil {
		log.Error(fmt.Sprintf("Service - TaskService - GetByStatus: %v", err))
		return nil, ErrCannotGetTasks
	}

	return tasks, nil
}

func (t *TaskService) GetOverdueTasks(ctx context.Context, log *slog.Logger) ([]entity.Task, error) {
	log.Info("Service - TaskService - GetOverdueTasks")

	tasks, err := t.taskRepo.GetOverdueTasks(ctx)
	if err != nil {
		log.Error(fmt.Sprintf("Service - TaskService - GetOverdueTasks: %v", err))
		return nil, ErrCannotGetTasks
	}

	return tasks, nil
}

func (t *TaskService) Complete(ctx context.Context, log *slog.Logger, input TaskCompleteInput) error {
	log.Info("Service - TaskService - Complete")

	// Проверяем, что задача существует
	task, err := t.taskRepo.GetByID(ctx, input.Id)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return ErrTaskNotFound
		}
		log.Error(fmt.Sprintf("Service - TaskService - Complete - GetByID: %v", err))
		return ErrCannotGetTask
	}

	// Проверяем, что задача назначена на текущего пользователя
	if task.AssignedTo != input.UserId {
		return ErrForbidden // Задача не назначена на текущего пользователя
	}

	// Обновляем статус задачи
	err = t.taskRepo.Complete(ctx, input.Id, input.UserId)
	if err != nil {
		log.Error(fmt.Sprintf("Service - TaskService - Complete: %v", err))
		return ErrCannotCompleteTask
	}

	return nil
}
