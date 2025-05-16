package v1

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"family-flow-app/internal/service"
	"family-flow-app/pkg/response"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

const (
	shoppingString = "/shopping"
)

type ShoppingRoutes struct {
	shoppingService     service.ShoppingItem
	notificationService service.Notification
	familyService       service.Family
}

func NewShoppingRoutes(
	ctx context.Context, log *slog.Logger, route chi.Router, shoppingService service.ShoppingItem,
	notificationService service.Notification, familyService service.Family,
) {
	u := ShoppingRoutes{
		shoppingService: shoppingService, notificationService: notificationService, familyService: familyService,
	}
	route.Route(
		shoppingString, func(r chi.Router) {
			r.Post("/", u.create(ctx, log))
			r.Put("/{id}", u.update(ctx, log))
			r.Delete("/{id}", u.delete(ctx, log))
			r.Get("/public", u.getPublicByFamilyID(ctx, log))
			r.Get("/private", u.getPrivateByCreatedBy(ctx, log))
			r.Put("/reserved/{id}", u.updateReservedBy(ctx, log))
			r.Put("/buyer/{id}", u.updateBuyerId(ctx, log))
			r.Get("/archived", u.getArchivedByUserID(ctx, log))
			r.Put("/cancel_reserved/{id}", u.cancelUpdateReservedBy(ctx, log))
		},
	)
}

type inputShoppingCreate struct {
	FamilyId    string `json:"family_id" validate:"required"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	Visibility  string `json:"visibility" validate:"required"`
}

// @Summary Create shopping item
// @Description Create shopping item
// @Tags shopping
// @Accept json
// @Produce json
// @Param family_id body string true "Family ID"
// @Param title body string true "Title"
// @Param description body string true "Description"
// @Param status body string true "Status"
// @Param visibility body string true "Visibility"
// @Param created_by body string true "Created by"
// @Success 201 {string} string "Shopping item created"
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /shopping [post]
func (u *ShoppingRoutes) create(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input inputShoppingCreate
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

		_, err = u.shoppingService.Create(
			ctx, log, service.ShoppingCreateInput{
				FamilyID:    input.FamilyId,
				Title:       input.Title,
				Description: input.Description,
				Visibility:  input.Visibility,
				CreatedBy:   user.Id,
			},
		)

		if err != nil {
			response.NewError(w, r, log, err, http.StatusInternalServerError, "Failed to create shopping item")
			return
		}

		if input.Visibility == "Public" {
			family, err := u.familyService.GetByFamilyID(ctx, log, input.FamilyId)
			if err != nil {
				response.NewError(w, r, log, err, http.StatusInternalServerError, "Failed to get family members")
				return
			}
			for _, familyMember := range family {
				if familyMember.Id == user.Id {
					continue
				}
				err = u.notificationService.SendNotification(
					ctx, log, service.NotificationCreateInput{
						UserID: familyMember.Id, // Отправляем всем членам семьи
						Title:  "Новый элемент в списке покупок",
						Body:   fmt.Sprintf("Пользователь %s добавил новый элемент: '%s'", user.Name, input.Title),
					},
				)
				if err != nil {
					log.Error("Failed to send notification: %v", err)
				}
			}
		}

		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, "Shopping item created")
	}
}

// @Summary Delete shopping item
// @Description Delete shopping item
// @Tags shopping
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success 200 {string} string "Shopping item deleted"
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /shopping/{id} [delete]
func (u *ShoppingRoutes) delete(ctx context.Context, log *slog.Logger) http.HandlerFunc {
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

		err = u.shoppingService.Delete(ctx, log, id)
		if err != nil {
			response.NewError(w, r, log, err, http.StatusInternalServerError, "Failed to delete shopping item")
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, "Shopping item deleted")
	}
}

type inputShoppingUpdate struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Visibility  string `json:"visibility"`
	IsArchived  bool   `json:"is_archived"`
}

// @Summary Update shopping item
// @Description Update shopping item
// @Tags shopping
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Param title body string true "Title"
// @Param description body string true "Description"
// @Param status body string true "Status"
// @Param visibility body string true "Visibility"
// @Success 200 {string} string "Shopping item updated"
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /shopping/{id} [put]
func (u *ShoppingRoutes) update(ctx context.Context, log *slog.Logger) http.HandlerFunc {
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

		var input inputShoppingUpdate

		if err = render.DecodeJSON(r.Body, &input); err != nil {
			response.NewError(w, r, log, err, http.StatusBadRequest, "Failed to parse request")
			return
		}

		if err = validator.New().Struct(input); err != nil {
			response.NewValidateError(w, r, log, http.StatusBadRequest, "Invalid request", err)
			return
		}

		log.Info(
			"ShoppingRoutes - Update - Input fields",
			"Title", input.Title,
			"Description", input.Description,
			"Status", input.Status,
			"Visibility", input.Visibility,
			"IsArchived", input.IsArchived,
		)

		err = u.shoppingService.Update(
			ctx, log, service.ShoppingUpdateInput{
				ID:          id,
				Title:       input.Title,
				Description: input.Description,
				Status:      input.Status,
				Visibility:  input.Visibility,
				IsArchived:  input.IsArchived,
			},
		)
		if err != nil {
			response.NewError(w, r, log, err, http.StatusInternalServerError, "Failed to update shopping item")
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, "Shopping item updated")
	}
}

// @Summary Get public shopping items by family ID
// @Description Get public shopping items by family ID
// @Tags shopping
// @Accept json
// @Produce json
// @Success 200 {object} []entity.ShoppingItem
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /shopping/public [get]
func (u *ShoppingRoutes) getPublicByFamilyID(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := GetCurrentUserFromContext(r.Context())
		if err != nil {
			response.NewError(w, r, log, err, http.StatusUnauthorized, ErrNoUserInContextMsg)
			return
		}

		familyID := r.URL.Query().Get("family_id")
		if familyID == "" {
			response.NewError(w, r, log, nil, http.StatusBadRequest, "Invalid request")
			return
		}

		items, err := u.shoppingService.GetPublicByFamilyID(ctx, log, familyID)
		if err != nil {
			response.NewError(w, r, log, err, http.StatusInternalServerError, "Failed to get public shopping items")
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, items)
	}
}

// @Summary Get private shopping items by created by
// @Description Get private shopping items by created by
// @Tags shopping
// @Accept json
// @Produce json
// @Success 200 {object} []entity.ShoppingItem
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /shopping/private [get]
func (u *ShoppingRoutes) getPrivateByCreatedBy(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := GetCurrentUserFromContext(r.Context())
		if err != nil {
			response.NewError(w, r, log, err, http.StatusUnauthorized, ErrNoUserInContextMsg)
			return
		}

		items, err := u.shoppingService.GetPrivateByCreatedBy(ctx, log, user.Id)
		if err != nil {
			response.NewError(w, r, log, err, http.StatusInternalServerError, "Failed to get private shopping items")
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, items)
	}
}

// @Summary Update reserved by
// @Description Update reserved by
// @Tags shopping
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Param reserved_by body string true "Reserved by"
// @Success 200 {string} string "Shopping item updated"
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /shopping/reserved/{id} [put]
func (u *ShoppingRoutes) updateReservedBy(ctx context.Context, log *slog.Logger) http.HandlerFunc {
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

		err = u.shoppingService.UpdateReservedBy(
			ctx, log, service.ShoppingUpdateReservedByInput{
				Id:         id,
				ReservedBy: user.Id,
			},
		)
		if err != nil {
			response.NewError(w, r, log, err, http.StatusInternalServerError, "Failed to update reserved by")
			return
		}

		item, err := u.shoppingService.GetByID(ctx, log, id)
		if err != nil {
			response.NewError(w, r, log, err, http.StatusInternalServerError, "Failed to get shopping item")
			return
		}
		// Отправляем уведомление создателю элемента
		err = u.notificationService.SendNotification(
			ctx, log, service.NotificationCreateInput{
				UserID: user.Id, // ID создателя элемента
				Title:  "Элемент зарезервирован",
				Body:   fmt.Sprintf("Элемент '%s' был зарезервирован пользователем %s", item.Title, user.Name),
			},
		)
		if err != nil {
			log.Error("Failed to send notification: %v", err)
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, "Shopping item updated")
	}
}

// @Summary Update buyer ID
// @Description Update buyer ID
// @Tags shopping
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Param buyer_id body string true "Buyer ID"
// @Success 200 {string} string "Shopping item updated"
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /shopping/buyer/{id} [put]
func (u *ShoppingRoutes) updateBuyerId(ctx context.Context, log *slog.Logger) http.HandlerFunc {
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

		err = u.shoppingService.UpdateBuyerId(
			ctx, log, service.ShoppingUpdateBuyerIdInput{
				Id:      id,
				BuyerId: user.Id,
			},
		)
		if err != nil {
			response.NewError(w, r, log, err, http.StatusInternalServerError, "Failed to update buyer ID")
			return
		}

		item, err := u.shoppingService.GetByID(ctx, log, id)
		if err != nil {
			response.NewError(w, r, log, err, http.StatusInternalServerError, "Failed to get shopping item")
			return
		}

		// Отправляем уведомление создателю элемента
		err = u.notificationService.SendNotification(
			ctx, log, service.NotificationCreateInput{
				UserID: user.Id, // ID создателя элемента
				Title:  "Элемент куплен",
				Body:   fmt.Sprintf("Элемент '%s' был куплен пользователем %s", item.Title, user.Name),
			},
		)
		if err != nil {
			log.Error("Failed to send notification: %v", err)
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, "Shopping item updated")
	}
}

// @Summary Get archived shopping items by user ID
// @Description Get archived shopping items by user ID
// @Tags shopping
// @Accept json
// @Produce json
// @Success 200 {object} []entity.ShoppingItem
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /shopping/archived [get]
func (u *ShoppingRoutes) getArchivedByUserID(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := GetCurrentUserFromContext(r.Context())
		if err != nil {
			response.NewError(w, r, log, err, http.StatusUnauthorized, ErrNoUserInContextMsg)
			return
		}

		items, err := u.shoppingService.GetArchivedByUserID(ctx, log, user.Id)
		if err != nil {
			response.NewError(w, r, log, err, http.StatusInternalServerError, "Failed to get archived shopping items")
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, items)
	}
}

// @Summary Cancel update reserved by
// @Description Cancel update reserved by
// @Tags shopping
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success 200 {string} string "Shopping item updated"
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /shopping/cancel_reserved/{id} [put]
func (u *ShoppingRoutes) cancelUpdateReservedBy(ctx context.Context, log *slog.Logger) http.HandlerFunc {
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

		err = u.shoppingService.CancelUpdateReservedBy(
			ctx, log, service.ShoppingCancelUpdateReservedByInput{
				Id: id,
			},
		)
		if err != nil {
			response.NewError(w, r, log, err, http.StatusInternalServerError, "Failed to cancel update reserved by")
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, "Shopping item updated")
	}
}
