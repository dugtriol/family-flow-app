package v1

import (
	"context"
	"log/slog"
	"net/http"

	"family-flow-app/internal/service"
	"family-flow-app/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

const (
	diaryPath = "/diary"
)

type DiaryRoutes struct {
	diaryService service.Diary
}

func NewDiaryRoutes(ctx context.Context, log *slog.Logger, route chi.Router, diaryService service.Diary) {
	d := DiaryRoutes{diaryService: diaryService}
	route.Route(diaryPath, func(r chi.Router) {
		r.Post("/", d.create(ctx, log))
		r.Get("/", d.getByUserID(ctx, log))
		r.Put("/{id}", d.update(ctx, log))
		r.Delete("/{id}", d.delete(ctx, log))
	})
}

type inputDiaryCreate struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	Emoji       string `json:"emoji"`
}

// @Summary Create diary entry
// @Description Create a new diary entry
// @Tags diary
// @Accept json
// @Produce json
// @Param input body inputDiaryCreate true "Diary entry data"
// @Success 201 {string} string "Diary entry created"
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /diary [post]
func (d *DiaryRoutes) create(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input inputDiaryCreate
		user, err := GetCurrentUserFromContext(r.Context())
		if err != nil {
			response.NewError(w, r, log, err, http.StatusUnauthorized, "Unauthorized")
			return
		}

		if err = render.DecodeJSON(r.Body, &input); err != nil {
			response.NewError(w, r, log, err, http.StatusBadRequest, "Invalid input")
			return
		}
		if err = validator.New().Struct(input); err != nil {
			response.NewValidateError(w, r, log, http.StatusBadRequest, "Validation error", err)
			return
		}

		id, err := d.diaryService.Create(ctx, log, service.DiaryCreateInput{
			Title:       input.Title,
			Description: input.Description,
			Emoji:       input.Emoji,
			CreatedBy:   user.Id,
		})
		if err != nil {
			response.NewError(w, r, log, err, http.StatusInternalServerError, "Failed to create diary entry")
			return
		}

		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, map[string]string{"id": id})
	}
}

// @Summary Get diary entries by user ID
// @Description Get all diary entries for a specific user
// @Tags diary
// @Accept json
// @Produce json
// @Success 200 {array} entity.DiaryItem
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /diary [get]
func (d *DiaryRoutes) getByUserID(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := GetCurrentUserFromContext(r.Context())
		if err != nil {
			response.NewError(w, r, log, err, http.StatusUnauthorized, "Unauthorized")
			return
		}

		items, err := d.diaryService.GetByUserID(ctx, log, user.Id)
		if err != nil {
			response.NewError(w, r, log, err, http.StatusInternalServerError, "Failed to get diary entries")
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, items)
	}
}

type inputDiaryUpdate struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	Emoji       string `json:"emoji"`
}

// @Summary Update diary entry
// @Description Update an existing diary entry
// @Tags diary
// @Accept json
// @Produce json
// @Param id path string true "Diary entry ID"
// @Param input body inputDiaryUpdate true "Diary entry data"
// @Success 200 {string} string "Diary entry updated"
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /diary/{id} [put]
func (d *DiaryRoutes) update(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			response.NewError(w, r, log, nil, http.StatusBadRequest, "Diary entry ID is required")
			return
		}

		var input inputDiaryUpdate
		if err := render.DecodeJSON(r.Body, &input); err != nil {
			response.NewError(w, r, log, err, http.StatusBadRequest, "Invalid input")
			return
		}
		if err := validator.New().Struct(input); err != nil {
			response.NewValidateError(w, r, log, http.StatusBadRequest, "Validation error", err)
			return
		}

		err := d.diaryService.Update(ctx, log, service.DiaryUpdateInput{
			ID:          id,
			Title:       input.Title,
			Description: input.Description,
			Emoji:       input.Emoji,
		})
		if err != nil {
			response.NewError(w, r, log, err, http.StatusInternalServerError, "Failed to update diary entry")
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, "Diary entry updated")
	}
}

// @Summary Delete diary entry
// @Description Delete a diary entry by ID
// @Tags diary
// @Accept json
// @Produce json
// @Param id path string true "Diary entry ID"
// @Success 200 {string} string "Diary entry deleted"
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /diary/{id} [delete]
func (d *DiaryRoutes) delete(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			response.NewError(w, r, log, nil, http.StatusBadRequest, "Diary entry ID is required")
			return
		}

		err := d.diaryService.Delete(ctx, log, id)
		if err != nil {
			response.NewError(w, r, log, err, http.StatusInternalServerError, "Failed to delete diary entry")
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, "Diary entry deleted")
	}
}
