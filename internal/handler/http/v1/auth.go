package v1

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	_ "family-flow-app/docs"
	"family-flow-app/pkg/response"

	"family-flow-app/internal/service"
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
			r.Get("/exists", u.existsByEmail(ctx, log))
			r.Put("/password", u.updatePassword(ctx, log))
		},
	)
}

type inputUserCreate struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Role     string `json:"role" validate:"required,oneof=Parent Child"`
}

type tokenResponse struct {
	Token string `json:"token"`
}

// @Summary Register new user
// @Description Register new user
// @Tags auth
// @Accept json
// @Produce json
// @Param name body string true "Name"
// @Param email body string true "Email"
// @Param password body string true "Password"
// @Param role body string true "Role"
// @Success 200 {object} tokenResponse
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /auth/register [post]
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

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, tokenResponse{Token: tokenString})
	}
}

type inputUserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// @Summary Login
// @Description Login
// @Tags auth
// @Accept json
// @Produce json
// @Param email body string true "Email"
// @Param password body string true "Password"
// @Success 200 {object} tokenResponse
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /auth/login [post]
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
			} else if errors.Is(err, service.ErrUserNotFound) {
				response.NewError(w, r, log, err, http.StatusBadRequest, MsgUserNotFound)
				return
			}
			log.Error(err.Error())
			response.NewError(w, r, log, err, http.StatusInternalServerError, MsgInternalServerErr)
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, tokenResponse{Token: tokenString})
	}
}

// @Summary Check if user exists by email
// @Description Check if user exists by email
// @Tags user
// @Accept json
// @Produce json
// @Param email query string true "Email"
// @Success 200 {object} map[string]bool
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /auth/exists [get]
func (u *AuthRoutes) existsByEmail(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.URL.Query().Get("email")
		if email == "" {
			response.NewError(w, r, log, nil, http.StatusBadRequest, "Email is required")
			return
		}

		exists, err := u.userService.ExistsByEmail(ctx, log, email)
		if err != nil {
			response.NewError(w, r, log, err, http.StatusInternalServerError, "Failed to check user existence")
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, map[string]bool{"exists": exists})
	}
}

type UpdatePasswordInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// @Summary Update user password
// @Description Update user password
// @Tags user
// @Accept json
// @Produce json
// @Param input body UpdatePasswordInput true "Update user password"
// @Success 200 {object} string
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /auth/password [put]
func (u *AuthRoutes) updatePassword(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		//user, err := GetCurrentUserFromContext(r.Context())
		//if err != nil {
		//	response.NewError(w, r, log, err, http.StatusUnauthorized, "Failed to get current user")
		//	return
		//}

		var input UpdatePasswordInput
		if err = render.DecodeJSON(r.Body, &input); err != nil {
			response.NewError(w, r, log, err, http.StatusBadRequest, MsgFailedParsing)
			return
		}
		if err = validator.New().Struct(input); err != nil {
			response.NewValidateError(w, r, log, http.StatusBadRequest, MsgInvalidReq, err)
			return
		}

		err = u.userService.UpdatePassword(ctx, log, input.Email, input.Password)
		if err != nil {
			response.NewError(
				w, r, log, err, http.StatusInternalServerError,
				"Failed to update user password",
			)
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, "Password updated successfully")
	}
}
