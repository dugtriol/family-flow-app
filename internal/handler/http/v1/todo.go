package v1

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"family-flow-app/internal/service"
	"family-flow-app/pkg/response"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

const (
	todoString = "/todo"
)

type TodoRoutes struct {
	todoService         service.TodoItem
	notificationService service.Notification
}

func NewTodoRoutes(
	ctx context.Context, log *slog.Logger, route chi.Router, todoService service.TodoItem,
	notificationService service.Notification,
) {
	u := TodoRoutes{todoService: todoService, notificationService: notificationService}
	route.Route(
		todoString, func(r chi.Router) {
			r.Post("/", u.create(ctx, log))
			r.Put("/{id}", u.update(ctx, log))
			r.Delete("/{id}", u.delete(ctx, log))
			r.Get("/assigned_to", u.getByAssignedTo(ctx, log))
			r.Get("/created_by", u.getByCreatedBy(ctx, log))
		},
	)
}

type inputTodoCreate struct {
	FamilyId    string    `json:"family_id" validate:"required"`
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Deadline    time.Time `json:"deadline" validate:"required"`
	AssignedTo  string    `json:"assigned_to" validate:"required"`
	Point       int       `json:"point"`
}

// @Summary Create todo
// @Description Create todo
// @Tags todo
// @Accept json
// @Produce json
// @Param family_id body string true "Family ID"
// @Param title body string true "Title"
// @Param description body string true "Description"
// @Param status body string true "Status"
// @Param deadline body string true "Deadline"
// @Param assigned_to body string true "Assigned to"
// @Success 201 {string} string "Todo created"
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /todo [post]
func (u *TodoRoutes) create(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input inputTodoCreate
		var err error

		user, err := GetCurrentUserFromContext(r.Context())
		if err != nil {
			response.NewError(w, r, log, err, http.StatusUnauthorized, ErrNoUserInContextMsg)
			return
		}

		if err = render.DecodeJSON(r.Body, &input); err != nil {
			response.NewError(w, r, log, err, http.StatusBadRequest, MsgFailedParsing)
			return
		}
		if err = validator.New().Struct(input); err != nil {
			response.NewValidateError(w, r, log, http.StatusBadRequest, MsgInvalidReq, err)
			return
		}

		_, err = u.todoService.Create(
			ctx, log, service.TodoCreateInput{
				FamilyId:    input.FamilyId,
				Title:       input.Title,
				Description: input.Description,
				Deadline:    input.Deadline,
				AssignedTo:  input.AssignedTo,
				CreatedBy:   user.Id,
				Point:       input.Point,
			},
		)

		if err != nil {
			response.NewError(w, r, log, err, http.StatusInternalServerError, "Failed to create family")
			return
		}

		// Отправляем уведомление пользователю, которому назначено задание
		err = u.notificationService.SendNotification(
			ctx, log, service.NotificationCreateInput{
				UserID: input.AssignedTo,
				Title:  "Новое задание",
				Body:   fmt.Sprintf("Вам назначено новое задание: '%s'", input.Title),
			},
		)
		if err != nil {
			log.Error("Failed to send notification: %v", err)
		}

		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, "Todo created")
	}
}

// @Summary Delete todo
// @Description Delete todo
// @Tags todo
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success 200 {string} string "Todo deleted"
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /todo/{id} [delete]
func (u *TodoRoutes) delete(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		_, err := GetCurrentUserFromContext(r.Context())
		if err != nil {
			response.NewError(w, r, log, err, http.StatusUnauthorized, ErrNoUserInContextMsg)
			return
		}

		id := chi.URLParam(r, "id")
		if id == "" {
			response.NewError(w, r, log, nil, http.StatusBadRequest, MsgInvalidReq)
			return
		}

		err = u.todoService.Delete(ctx, log, id)
		if err != nil {
			response.NewError(w, r, log, err, http.StatusInternalServerError, "Failed to delete todo")
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, "Todo deleted")
	}
}

type inputTodoUpdate struct {
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Status      string    `json:"status" validate:"required"`
	Deadline    time.Time `json:"deadline" validate:"required"`
	AssignedTo  string    `json:"assigned_to" validate:"required"`
	Point       int       `json:"point"`
}

// @Summary Update todo
// @Description Update todo
// @Tags todo
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Param title body string true "Title"
// @Param description body string true "Description"
// @Param status body string true "Status"
// @Param deadline body string true "Deadline"
// @Param assigned_to body string true "Assigned to"
// @Success 200 {string} string "Todo updated"
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /todo/{id} [put]
func (u *TodoRoutes) update(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := GetCurrentUserFromContext(r.Context())
		if err != nil {
			response.NewError(w, r, log, err, http.StatusUnauthorized, ErrNoUserInContextMsg)
			return
		}

		id := chi.URLParam(r, "id")
		log.Info("id: %v", id)
		if id == "" {
			response.NewError(w, r, log, nil, http.StatusBadRequest, MsgInvalidReq)
			return
		}

		var input inputTodoUpdate

		if err = render.DecodeJSON(r.Body, &input); err != nil {
			response.NewError(w, r, log, err, http.StatusBadRequest, MsgFailedParsing)
			return
		}
		if err = validator.New().Struct(input); err != nil {
			response.NewValidateError(w, r, log, http.StatusBadRequest, MsgInvalidReq, err)
			return
		}

		err = u.todoService.Update(
			ctx, log, service.TodoUpdateInput{
				ID:          id,
				Title:       input.Title,
				Description: input.Description,
				Status:      input.Status,
				Deadline:    input.Deadline,
				AssignedTo:  input.AssignedTo,
				Point:       input.Point,
			},
		)
		if err != nil {
			response.NewError(w, r, log, err, http.StatusInternalServerError, "Failed to update todo")
			return
		}

		// Если задание завершено, отправляем уведомление создателю
		if input.Status == "Completed" {
			// Получаем задание, чтобы узнать, кто его создал
			todo, err := u.todoService.GetByID(ctx, log, id)
			if err != nil {
				response.NewError(w, r, log, err, http.StatusInternalServerError, "Failed to fetch todo")
				return
			}

			// Отправляем уведомление создателю задания
			err = u.notificationService.SendNotification(
				ctx, log, service.NotificationCreateInput{
					UserID: todo.AssignedTo,
					Title:  "Задание выполнено",
					Body:   fmt.Sprintf("Задание '%s' было выполнено пользователем %s", todo.Title, user.Id),
				},
			)
			if err != nil {
				log.Error("Failed to send notification: %v", err)
			}
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, "Todo updated")
	}
}

// @Summary Get todo by assigned to
// @Description Get todo by assigned to
// @Tags todo
// @Accept json
// @Produce json
// @Success 200 {object} []entity.TodoItem
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /todo/assigned_to [get]
func (u *TodoRoutes) getByAssignedTo(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := GetCurrentUserFromContext(r.Context())
		if err != nil {
			response.NewError(w, r, log, err, http.StatusUnauthorized, ErrNoUserInContextMsg)
			return
		}

		items, err := u.todoService.GetByAssignedTo(ctx, log, user.Id)
		if err != nil {
			response.NewError(w, r, log, err, http.StatusInternalServerError, "Failed to get todo by assigned to")
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, items)
	}
}

// @Summary Get todo by created by
// @Description Get todo by created by
// @Tags todo
// @Accept json
// @Produce json
// @Success 200 {object} []entity.TodoItem
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /todo/created_by [get]
func (u *TodoRoutes) getByCreatedBy(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := GetCurrentUserFromContext(r.Context())
		if err != nil {
			response.NewError(w, r, log, err, http.StatusUnauthorized, ErrNoUserInContextMsg)
			return
		}

		items, err := u.todoService.GetByCreatedBy(ctx, log, user.Id)
		if err != nil {
			response.NewError(w, r, log, err, http.StatusInternalServerError, "Failed to get todo by created by")
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, items)
	}
}
