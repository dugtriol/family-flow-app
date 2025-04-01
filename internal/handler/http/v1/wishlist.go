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
	wishlistString = "/wishlist"
)

type WishlistRoutes struct {
	wishlistService service.WishlistItem
}

func NewWishlistRoutes(
	ctx context.Context, log *slog.Logger, route chi.Router, wishlistService service.WishlistItem,
) {
	u := WishlistRoutes{wishlistService: wishlistService}
	route.Route(
		wishlistString, func(r chi.Router) {
			r.Post("/", u.create(ctx, log))
			r.Put("/{id}", u.update(ctx, log))
			r.Delete("/{id}", u.delete(ctx, log))
			r.Get("/", u.getByUserID(ctx, log))
		},
	)
}

type inputWishlistCreate struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	Link        string `json:"link" validate:"required"`
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
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	Link        string `json:"link" validate:"required"`
	Status      string `json:"status" validate:"required"`
	IsReserved  bool   `json:"is_reserved" validate:"required"`
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
				Name:        input.Name,
				Description: input.Description,
				Link:        input.Link,
				Status:      input.Status,
				IsReserved:  input.IsReserved,
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

		items, err := u.wishlistService.GetByID(ctx, log, user.Id)
		if err != nil {
			response.NewError(w, r, log, err, http.StatusInternalServerError, "Failed to get wishlist items")
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, items)
	}
}
