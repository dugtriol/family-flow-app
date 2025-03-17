package v1

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"family-flow-app/internal/entity"
	"family-flow-app/internal/service"
	"family-flow-app/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

const (
	taskString = "/task"
)

type TaskRoutes struct {
	taskService service.Task
}

func NewTaskRoutes(ctx context.Context, log *slog.Logger, route chi.Router, taskService service.Task) {
	t := TaskRoutes{taskService: taskService}
	route.Route(
		taskString, func(r chi.Router) {
			r.Get("/", t.list(ctx, log))                   // GET /tasks
			r.Post("/", t.create(ctx, log))                // POST /tasks
			r.Get("/{id}", t.get(ctx, log))                // GET /tasks/{id}
			r.Put("/{id}", t.update(ctx, log))             // PUT /tasks/{id}
			r.Delete("/{id}", t.delete(ctx, log))          // DELETE /tasks/{id}
			r.Post("/{id}/complete", t.complete(ctx, log)) // POST /tasks/{id}/complete
		},
	)
}

type inputTaskList struct {
	Status string `json:"status" validate:"omitempty,oneof=active completed overdue"`
}

func (t *TaskRoutes) list(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input inputTaskList
		var err error

		// Получаем текущего пользователя из контекста
		user, err := GetCurrentUserFromContext(r.Context())
		if err != nil {
			response.NewError(w, r, log, err, http.StatusUnauthorized, "Failed to get current user")
			return
		}

		// Парсим параметры запроса
		if err := render.DecodeJSON(r.Body, &input); err != nil {
			response.NewError(w, r, log, err, http.StatusBadRequest, MsgFailedParsing)
			return
		}

		// Валидация входных данных
		if err := validator.New().Struct(input); err != nil {
			response.NewValidateError(w, r, log, http.StatusBadRequest, MsgInvalidReq, err)
			return
		}

		// Получаем задачи в зависимости от статуса
		var tasks []entity.Task
		switch input.Status {
		case "active":
			tasks, err = t.taskService.GetByAssignedTo(ctx, log, service.TaskGetByAssignedToInput{AssignedTo: user.Id})
		case "completed":
			tasks, err = t.taskService.GetByStatus(ctx, log, service.TaskGetByStatusInput{Status: "completed"})
		case "overdue":
			tasks, err = t.taskService.GetOverdueTasks(ctx, log)
		default:
			tasks, err = t.taskService.GetByCreatedBy(ctx, log, service.TaskGetByCreatedByInput{CreatedBy: user.Id})
		}

		if err != nil {
			response.NewError(w, r, log, err, http.StatusInternalServerError, MsgInternalServerErr)
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, tasks)
	}
}

type inputTaskCreate struct {
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Deadline    time.Time `json:"deadline" validate:"required"`
	AssignedTo  string    `json:"assigned_to" validate:"required,uuid"`
	Reward      int       `json:"reward" validate:"required"`
}

func (t *TaskRoutes) create(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input inputTaskCreate
		var err error

		// Получаем текущего пользователя из контекста
		user, err := GetCurrentUserFromContext(r.Context())
		if err != nil {
			response.NewError(w, r, log, err, http.StatusUnauthorized, "Failed to get current user")
			return
		}

		// Парсим входные данные
		if err := render.DecodeJSON(r.Body, &input); err != nil {
			response.NewError(w, r, log, err, http.StatusBadRequest, MsgFailedParsing)
			return
		}

		// Валидация входных данных
		if err := validator.New().Struct(input); err != nil {
			response.NewValidateError(w, r, log, http.StatusBadRequest, MsgInvalidReq, err)
			return
		}

		// Создаем задачу
		id, err := t.taskService.Create(
			ctx, log, service.TaskCreateInput{
				Title:       input.Title,
				Description: input.Description,
				Status:      "active", // По умолчанию задача активна
				Deadline:    input.Deadline,
				AssignedTo:  input.AssignedTo,
				CreatedBy:   user.Id,
				Reward:      input.Reward,
			},
		)
		if err != nil {
			response.NewError(w, r, log, err, http.StatusInternalServerError, MsgInternalServerErr)
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, map[string]string{"id": id})
	}
}

type inputTaskGet struct {
	Id string `validate:"uuid"`
}

func (t *TaskRoutes) get(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var id string
		var err error

		// Получаем ID задачи из URL
		id = chi.URLParam(r, "id")
		if err := validator.New().Struct(inputTaskGet{Id: id}); err != nil {
			response.NewValidateError(w, r, log, http.StatusBadRequest, MsgInvalidReq, err)
			return
		}

		// Получаем задачу по ID
		task, err := t.taskService.GetById(ctx, log, service.TaskGetByIdInput{Id: id})
		if err != nil {
			if errors.Is(err, service.ErrTaskNotFound) {
				response.NewError(w, r, log, err, http.StatusNotFound, MsgTaskNotFound)
				return
			}
			response.NewError(w, r, log, err, http.StatusInternalServerError, MsgInternalServerErr)
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, task)
	}
}

type inputTaskUpdate struct {
	Id          string    `json:"id" validate:"required,uuid"`
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Deadline    time.Time `json:"deadline" validate:"required"`
	AssignedTo  string    `json:"assigned_to" validate:"required,uuid"`
	Reward      int       `json:"reward" validate:"required"`
}

func (t *TaskRoutes) update(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input inputTaskUpdate
		var err error

		// Парсим входные данные
		if err := render.DecodeJSON(r.Body, &input); err != nil {
			response.NewError(w, r, log, err, http.StatusBadRequest, MsgFailedParsing)
			return
		}

		// Валидация входных данных
		if err := validator.New().Struct(input); err != nil {
			response.NewValidateError(w, r, log, http.StatusBadRequest, MsgInvalidReq, err)
			return
		}

		// Обновляем задачу
		err = t.taskService.Update(
			ctx, log, service.TaskUpdateInput{
				Id:          input.Id,
				Title:       input.Title,
				Description: input.Description,
				Deadline:    input.Deadline,
				AssignedTo:  input.AssignedTo,
				Reward:      input.Reward,
			},
		)
		if err != nil {
			response.NewError(w, r, log, err, http.StatusInternalServerError, MsgInternalServerErr)
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, map[string]string{"message": "Task updated successfully"})
	}
}

type inputTaskDelete struct {
	Id string `validate:"uuid"`
}

func (t *TaskRoutes) delete(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var id string
		var err error

		// Получаем ID задачи из URL
		id = chi.URLParam(r, "id")
		if err := validator.New().Struct(inputTaskDelete{Id: id}); err != nil {
			response.NewValidateError(w, r, log, http.StatusBadRequest, MsgInvalidReq, err)
			return
		}

		// Удаляем задачу
		err = t.taskService.Delete(ctx, log, service.TaskDeleteInput{Id: id})
		if err != nil {
			response.NewError(w, r, log, err, http.StatusInternalServerError, MsgInternalServerErr)
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, map[string]string{"message": "Task deleted successfully"})
	}
}

type inputTaskComplete struct {
	Id string `validate:"uuid"`
}

func (t *TaskRoutes) complete(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var id string
		var err error

		// Получаем ID задачи из URL
		id = chi.URLParam(r, "id")
		if err := validator.New().Struct(inputTaskComplete{Id: id}); err != nil {
			response.NewValidateError(w, r, log, http.StatusBadRequest, MsgInvalidReq, err)
			return
		}

		// Получаем текущего пользователя из контекста
		user, err := GetCurrentUserFromContext(r.Context())
		if err != nil {
			response.NewError(w, r, log, err, http.StatusUnauthorized, "Failed to get current user")
			return
		}

		// Обновляем статус задачи
		err = t.taskService.Complete(
			ctx, log, service.TaskCompleteInput{
				Id:     id,
				UserId: user.Id,
			},
		)
		if err != nil {
			response.NewError(w, r, log, err, http.StatusInternalServerError, MsgInternalServerErr)
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, map[string]string{"message": "Task completed successfully"})
	}
}
