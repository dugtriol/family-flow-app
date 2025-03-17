package v1

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"family-flow-app/internal/service"
	"family-flow-app/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

const (
	authString = "/auth"
)

type AuthRoutes struct {
	userService service.User
}

func NewAuthRoutes(ctx context.Context, log *slog.Logger, route chi.Router, userService service.User) {
	u := AuthRoutes{userService: userService}
	route.Route(
		authString, func(r chi.Router) {
			r.Post("/register", u.create(ctx, log))
			r.Post("/login", u.login(ctx, log))
		},
	)
}

type inputUserCreate struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Role     string `json:"role" validate:"required,oneof=Parent Child"`
}

func (u *AuthRoutes) create(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input inputUserCreate
		var err error

		if err = render.DecodeJSON(r.Body, &input); err != nil {
			response.NewError(w, r, log, err, http.StatusBadRequest, MsgFailedParsing)
			return
		}
		if err = validator.New().Struct(input); err != nil {
			response.NewValidateError(w, r, log, http.StatusBadRequest, MsgInvalidReq, err)
			return
		}
		tokenString, err := u.userService.Create(
			ctx, log, service.UserCreateInput{
				Name:     input.Name,
				Email:    input.Email,
				Password: input.Password,
				Role:     input.Role,
			},
		)
		if err != nil {
			if errors.Is(err, service.ErrUserAlreadyExists) {
				response.NewError(w, r, log, err, http.StatusBadRequest, MsgUserAlreadyExists)
				return
			}
			response.NewError(w, r, log, err, http.StatusInternalServerError, MsgInternalServerErr)
			return
		}

		type resp struct {
			Token string `json:"token"`
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, resp{Token: tokenString})
	}
}

type inputUserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (u *AuthRoutes) login(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input inputUserLogin
		var err error
		var tokenString string

		if err = render.DecodeJSON(r.Body, &input); err != nil {
			response.NewError(w, r, log, err, http.StatusBadRequest, MsgFailedParsing)
			return
		}
		if err = validator.New().Struct(input); err != nil {
			response.NewValidateError(w, r, log, http.StatusBadRequest, MsgInvalidReq, err)
			return
		}
		if tokenString, err = u.userService.Login(
			ctx, log, service.AuthInput{
				Email:    input.Email,
				Password: input.Password,
			},
		); err != nil {
			if errors.Is(err, service.ErrInvalidPassword) {
				response.NewError(w, r, log, err, http.StatusBadRequest, MsgInvalidPasswordErr)
				return
			}
			log.Error(err.Error())
			response.NewError(w, r, log, err, http.StatusInternalServerError, MsgInternalServerErr)
			return
		}

		type resp struct {
			Token string `json:"token"`
		}
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, resp{Token: tokenString})
	}
}
