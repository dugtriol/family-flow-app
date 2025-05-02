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
	notificationString = "/notification"
)

type NotificationRoutes struct {
	notificationService service.Notification
}

func NewNotificationRoutes(
	ctx context.Context, log *slog.Logger, route chi.Router, notificationService service.Notification,
) {
	u := NotificationRoutes{notificationService: notificationService}
	route.Route(
		notificationString, func(r chi.Router) {
			r.Post("/save-fcm-token", u.saveFcmToken(ctx, log))
			r.Post("/send", u.sendNotification(ctx, log))
			r.Get("/", u.getNotificationsByUserID(ctx, log))
		},
	)
}

type inputSaveFcmToken struct {
	//UserID   string `json:"userId" validate:"required,uuid"`
	FcmToken string `json:"fcmToken" validate:"required"`
}

// @Summary Save FCM token
// @Description Save FCM token for a user
// @Tags notification
// @Accept json
// @Produce json
// @Param userId body string true "User ID"
// @Param fcmToken body string true "FCM Token"
// @Success 200 {string} string "FCM token saved successfully"
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /notification/save-fcm-token [post]
func (u *NotificationRoutes) saveFcmToken(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input inputSaveFcmToken
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

		err = u.notificationService.SaveToken(ctx, user.Id, input.FcmToken)
		if err != nil {
			response.NewError(w, r, log, err, http.StatusInternalServerError, "Failed to save FCM token")
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, "FCM token saved successfully")
	}
}

type inputSendNotification struct {
	//UserID string `json:"userId" validate:"required,uuid"`
	Title string `json:"title" validate:"required"`
	Body  string `json:"body" validate:"required"`
}

// @Summary Send notification
// @Description Send notification to a user
// @Tags notification
// @Accept json
// @Produce json
// @Param title body string true "Title"
// @Param body body string true "Body"
// @Success 200 {string} string "Notification sent successfully"
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /notification/send [post]
func (u *NotificationRoutes) sendNotification(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input inputSendNotification
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

		err = u.notificationService.SendNotification(
			ctx, log, service.NotificationCreateInput{
				UserID: user.Id,
				Title:  input.Title,
				Body:   input.Body,
			},
		)
		if err != nil {
			response.NewError(w, r, log, err, http.StatusInternalServerError, "Failed to send notification")
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, "Notification sent successfully")
	}
}

// @Summary Get notifications by user ID
// @Description Get notifications for a specific user
// @Tags notification
// @Accept json
// @Produce json
// @Param userID path string true "User ID"
// @Success 200 {object} []entity.Notification
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /notification [get]
func (u *NotificationRoutes) getNotificationsByUserID(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user, err := GetCurrentUserFromContext(r.Context())
		if err != nil {
			response.NewError(w, r, log, err, http.StatusUnauthorized, ErrNoUserInContextMsg)
			return
		}

		notifications, err := u.notificationService.GetNotificationsByUserID(ctx, log, user.Id)
		if err != nil {
			response.NewError(w, r, log, err, http.StatusInternalServerError, "Failed to get notifications")
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, notifications)
	}
}
