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

type EmailRequestBody struct {
	ToAddr  string `json:"to_addr"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

const (
	emailString = "/email"
)

type EmailRoutes struct {
	emailService service.Email
}

func NewEmailRoutes(ctx context.Context, log *slog.Logger, route chi.Router, emailService service.Email) {
	u := EmailRoutes{emailService: emailService}
	route.Route(
		emailString, func(r chi.Router) {
			r.Post("/send", u.sendCode(ctx, log))
			r.Post("/compare", u.compareCode(ctx, log))
			//r.Post("/invite", u.sendInvite(ctx, log))
		},
	)
}

type inputSendCode struct {
	Email string `json:"email" validate:"required,email"`
}

// @Summary Send verification code
// @Description Send verification code
// @Tags email
// @Accept json
// @Produce json
// @Param email body string true "Email"
// @Success 200 {string} string "Verification code sent"
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /email/send [post]
func (u *EmailRoutes) sendCode(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input inputSendCode
		var err error

		if err = render.DecodeJSON(r.Body, &input); err != nil {
			response.NewError(w, r, log, err, http.StatusBadRequest, MsgFailedParsing)
			return
		}
		if err = validator.New().Struct(input); err != nil {
			response.NewValidateError(w, r, log, http.StatusBadRequest, MsgInvalidReq, err)
			return
		}

		err = u.emailService.SendCode(ctx, []string{input.Email})
		if err != nil {
			response.NewError(w, r, log, err, http.StatusInternalServerError, "Failed to send verification code")
			return
		}

		keys, err := u.emailService.GetAllKeys(ctx)
		if err != nil {
			log.Error("Failed to get all keys", "error", err)
		} else {
			log.Info("All keys retrieved", "keys", keys)
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, "Verification code sent")
	}
}

type inputCompareCode struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required"`
}

// @Summary Compare verification code
// @Description Compare verification code
// @Tags email
// @Accept json
// @Produce json
// @Param email body string true "Email"
// @Param code body string true "Code"
// @Success 200 {string} string "Verification code compared"
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /email/compare [post]
func (u *EmailRoutes) compareCode(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input inputCompareCode
		var err error
		var status bool

		if err = render.DecodeJSON(r.Body, &input); err != nil {
			response.NewError(w, r, log, err, http.StatusBadRequest, MsgFailedParsing)
			return
		}
		if err = validator.New().Struct(input); err != nil {
			response.NewValidateError(w, r, log, http.StatusBadRequest, MsgInvalidReq, err)
			return
		}

		keys, err := u.emailService.GetAllKeys(ctx)
		if err != nil {
			log.Error("Failed to get all keys", "error", err)
		} else {
			log.Info("All keys retrieved - before compare", "keys", keys)
		}

		if status, err = u.emailService.CompareCode(ctx, input.Email, input.Code); err != nil {
			response.NewError(w, r, log, err, http.StatusInternalServerError, "Failed to compare verification code")
			return
		}

		keys1, err := u.emailService.GetAllKeys(ctx)
		if err != nil {
			log.Error("Failed to get all keys", "error", err)
		} else {
			log.Info("All keys retrieved - after compare", "keys", keys1)
		}

		if !status {
			response.NewError(w, r, log, nil, http.StatusBadRequest, "Verification code is invalid")
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, "Verification code compared")
	}
}
