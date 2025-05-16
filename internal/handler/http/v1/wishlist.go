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
	wishlistString = "/wishlist"
)

type WishlistRoutes struct {
	wishlistService     service.WishlistItem
	notificationService service.Notification
	familyService       service.Family
}

func NewWishlistRoutes(
	ctx context.Context, log *slog.Logger, route chi.Router, wishlistService service.WishlistItem, notificationService service.Notification, familyService service.Family,
) {
	u := WishlistRoutes{wishlistService: wishlistService, notificationService: notificationService, familyService: familyService}
	route.Route(
		wishlistString, func(r chi.Router) {
			r.Post("/", u.create(ctx, log))
			r.Put("/{id}", u.update(ctx, log))
			r.Delete("/{id}", u.delete(ctx, log))
			r.Get("/", u.getByUserID(ctx, log))
			r.Get("/{id}", u.getByFamilyUserID(ctx, log))
			r.Put("/{id}/reserved_by", u.updateReservedBy(ctx, log))
			r.Put("/{id}/cancel_reserved_by", u.cancelUpdateReservedBy(ctx, log))
			r.Get("/archived", u.getArchivedByUserID(ctx, log))
		},
	)
}

type inputWishlistCreate struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Link        string `json:"link"`
}

// @Summary Create wishlist item
// @Description Create wishlist item
// @Tags wishlist
// @Accept json
// @Produce json
// @Param name body string true "Name"
// @Param description body string true "Description"
// @Param link body string true "Link"
// @Param status body string true "Status"
// @Param is_reserved body bool true "Is Reserved"
// @Success 201 {string} string "Wishlist item created"
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /wishlist [post]
func (u *WishlistRoutes) create(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input inputWishlistCreate
		var err error

		user, err := GetCurrentUserFromContext(r.Context())
		if err != nil {
			response.NewError(w, r, log, err, http.StatusUnauthorized, ErrNoUserInContextMsg)
			return
		}

		if err = render.DecodeJSON(r.Body, &input); err != nil {
			response.NewError(w, r, log, err, http.StatusBadRequest, "Failed to parse request")
			return
		}
		if err = validator.New().Struct(input); err != nil {
			response.NewValidateError(w, r, log, http.StatusBadRequest, "Invalid request", err)
			return
		}

		_, err = u.wishlistService.Create(
			ctx, log, service.WishlistCreateInput{
				Name:        input.Name,
				Description: input.Description,
				Link:        input.Link,
				CreatedBy:   user.Id,
			},
		)

		if err != nil {
			response.NewError(w, r, log, err, http.StatusInternalServerError, "Failed to create wishlist item")
			return
		}

		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, "Wishlist item created")
	}
}

// @Summary Delete wishlist item
// @Description Delete wishlist item
// @Tags wishlist
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success 200 {string} string "Wishlist item deleted"
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /wishlist/{id} [delete]
func (u *WishlistRoutes) delete(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := GetCurrentUserFromContext(r.Context())
		if err != nil {
			response.NewError(w, r, log, err, http.StatusUnauthorized, ErrNoUserInContextMsg)
			return
		}

		id := chi.URLParam(r, "id")
		if id == "" {
			response.NewError(w, r, log, nil, http.StatusBadRequest, "Invalid request")
			return
		}

		err = u.wishlistService.Delete(ctx, log, id)
		if err != nil {
			response.NewError(w, r, log, err, http.StatusInternalServerError, "Failed to delete wishlist item")
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, "Wishlist item deleted")
	}
}

type inputWishlistUpdate struct {
	//Id          string `json:"id" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Link        string `json:"link"`
	Status      string `json:"status" validate:"required"`
	IsArchived  bool   `json:"is_archived"`
}

// @Summary Update wishlist item
// @Description Update wishlist item
// @Tags wishlist
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Param name body string true "Name"
// @Param description body string true "Description"
// @Param link body string true "Link"
// @Param status body string true "Status"
// @Param is_reserved body bool true "Is Reserved"
// @Success 200 {string} string "Wishlist item updated"
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /wishlist/{id} [put]
func (u *WishlistRoutes) update(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := GetCurrentUserFromContext(r.Context())
		if err != nil {
			response.NewError(w, r, log, err, http.StatusUnauthorized, ErrNoUserInContextMsg)
			return
		}

		id := chi.URLParam(r, "id")
		if id == "" {
			response.NewError(w, r, log, nil, http.StatusBadRequest, "Invalid request")
			return
		}

		var input inputWishlistUpdate

		if err = render.DecodeJSON(r.Body, &input); err != nil {
			response.NewError(w, r, log, err, http.StatusBadRequest, "Failed to parse request")
			return
		}
		if err = validator.New().Struct(input); err != nil {
			response.NewValidateError(w, r, log, http.StatusBadRequest, "Invalid request", err)
			return
		}

		err = u.wishlistService.Update(
			ctx, log, service.WishlistUpdateInput{
				ID:          id,
				Name:        input.Name,
				Description: input.Description,
				Link:        input.Link,
				Status:      input.Status,
				IsArchived:  input.IsArchived,
				CreatedBy:   user.Id,
				UpdatedAt:   time.Now().Add(time.Hour * 3),
			},
		)
		if err != nil {
			response.NewError(w, r, log, err, http.StatusInternalServerError, "Failed to update wishlist item")
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, "Wishlist item updated")
	}
}

// @Summary Get wishlist items by user ID
// @Description Get wishlist items by user ID
// @Tags wishlist
// @Accept json
// @Produce json
// @Success 200 {object} []entity.WishlistItem
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /wishlist [get]
func (u *WishlistRoutes) getByUserID(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := GetCurrentUserFromContext(r.Context())
		if err != nil {
			response.NewError(w, r, log, err, http.StatusUnauthorized, ErrNoUserInContextMsg)
			return
		}

		items, err := u.wishlistService.GetByIDs(ctx, log, user.Id)
		if err != nil {
			response.NewError(w, r, log, err, http.StatusInternalServerError, "Failed to get wishlist items")
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, items)
	}
}

// @Summary Get wishlist items by family user ID
// @Description Get wishlist items by user ID
// @Tags wishlist
// @Accept json
// @Produce json
// @Success 200 {object} []entity.WishlistItem
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /wishlist/{id} [get]
func (u *WishlistRoutes) getByFamilyUserID(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := GetCurrentUserFromContext(r.Context())
		if err != nil {
			response.NewError(w, r, log, err, http.StatusUnauthorized, ErrNoUserInContextMsg)
			return
		}

		id := chi.URLParam(r, "id")
		if id == "" {
			response.NewError(w, r, log, nil, http.StatusBadRequest, "Invalid request")
			return
		}

		items, err := u.wishlistService.GetByIDs(ctx, log, id)
		if err != nil {
			response.NewError(w, r, log, err, http.StatusInternalServerError, "Failed to get wishlist items")
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, items)
	}
}

// UpdateReservedBy
// @Summary Update wishlist item reserved by
// @Description Update wishlist item reserved by
// @Tags wishlist
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Param reserved_by body string true "Reserved By"
// @Success 200 {string} string "Wishlist item updated"
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /wishlist/{id}/reserved_by [put]
func (u *WishlistRoutes) updateReservedBy(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := GetCurrentUserFromContext(r.Context())
		if err != nil {
			response.NewError(w, r, log, err, http.StatusUnauthorized, ErrNoUserInContextMsg)
			return
		}

		id := chi.URLParam(r, "id")
		if id == "" {
			response.NewError(w, r, log, nil, http.StatusBadRequest, "Invalid request")
			return
		}

		var input struct {
			ReservedBy string `json:"reserved_by" validate:"required"`
		}

		if err = render.DecodeJSON(r.Body, &input); err != nil {
			response.NewError(w, r, log, err, http.StatusBadRequest, "Failed to parse request")
			return
		}
		if err = validator.New().Struct(input); err != nil {
			response.NewValidateError(w, r, log, http.StatusBadRequest, "Invalid request", err)
			return
		}

		err = u.wishlistService.UpdateReservedBy(
			ctx, log, service.WishlistUpdateReservedByInput{
				ID:         id,
				ReservedBy: user.Id,
			},
		)

		if err != nil {
			response.NewError(
				w, r, log, err, http.StatusInternalServerError,
				"Failed to update wishlist item reserved by",
			)
			return
		}
		// Получаем информацию о создателе элемента
		wishlistItem, err := u.wishlistService.GetByID(ctx, log, id)
		if err != nil {
			response.NewError(w, r, log, err, http.StatusInternalServerError, "Failed to get wishlist item")
			return
		}

		// Отправляем уведомление создателю элемента
		err = u.notificationService.SendNotification(
			ctx, log, service.NotificationCreateInput{
				UserID: wishlistItem.CreatedBy, // ID создателя элемента
				Title:  "Ваш элемент желаний зарезервирован",
				Body:   fmt.Sprintf("Ваш элемент '%s' был кем то зарезервирован!", wishlistItem.Name),
			},
		)
		if err != nil {
			log.Error("Failed to send notification: %v", err)
		}
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, "Wishlist item updated")
	}
}

// GetArchivedByUserID
// @Summary Get archived wishlist items by user ID
// @Description Get archived wishlist items by user ID
// @Tags wishlist
// @Accept json
// @Produce json
// @Success 200 {object} []entity.WishlistItem
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /wishlist/archived [get]
func (u *WishlistRoutes) getArchivedByUserID(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := GetCurrentUserFromContext(r.Context())
		if err != nil {
			response.NewError(w, r, log, err, http.StatusUnauthorized, ErrNoUserInContextMsg)
			return
		}

		items, err := u.wishlistService.GetArchivedByUserID(ctx, log, user.Id)
		if err != nil {
			response.NewError(w, r, log, err, http.StatusInternalServerError, "Failed to get archived wishlist items")
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, items)
	}
}

// cancel update reserved by
// @Summary Cancel update wishlist item reserved by
// @Description Cancel update wishlist item reserved by
// @Tags wishlist
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success 200 {string} string "Wishlist item updated"
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /wishlist/{id}/cancel_reserved_by [put]
func (u *WishlistRoutes) cancelUpdateReservedBy(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := GetCurrentUserFromContext(r.Context())
		if err != nil {
			response.NewError(w, r, log, err, http.StatusUnauthorized, ErrNoUserInContextMsg)
			return
		}

		id := chi.URLParam(r, "id")
		if id == "" {
			response.NewError(w, r, log, nil, http.StatusBadRequest, "Invalid request")
			return
		}

		err = u.wishlistService.CancelUpdateReservedBy(
			ctx, log, service.WishlistCancelUpdateReservedByInput{
				ID: id,
			},
		)

		if err != nil {
			response.NewError(
				w, r, log, err, http.StatusInternalServerError,
				"Failed to cancel update wishlist item reserved by",
			)
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, "Wishlist item updated")
	}
}
