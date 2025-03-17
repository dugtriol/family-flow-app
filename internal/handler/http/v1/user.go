package v1

import (
	"context"
	"errors"
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
	userString = "/user"
)

type UserRoutes struct {
	userService service.User
}

func NewUserRoutes(ctx context.Context, log *slog.Logger, route chi.Router, userService service.User) {
	u := UserRoutes{userService: userService}
	route.Route(
		userString, func(r chi.Router) {
			r.Get("/{id}", u.get(ctx, log))
		},
	)
}

type inputUserGet struct {
	Id string `validate:"uuid"`
}

func (u *UserRoutes) get(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var id string
		var err error

		id = chi.URLParam(r, "id")
		if err = validator.New().Struct(inputUserGet{Id: id}); err != nil {
			response.NewValidateError(w, r, log, http.StatusBadRequest, MsgInvalidReq, err)
			return
		}
		log.Info(fmt.Sprintf("Handler - User - Create - validate is ok"))
		user, err := u.userService.GetById(ctx, log, id)
		if err != nil {
			if errors.Is(err, service.ErrUserNotFound) {
				response.NewError(w, r, log, err, http.StatusBadRequest, MsgUserNotFound)
				return
			}
			response.NewError(w, r, log, err, http.StatusInternalServerError, MsgInternalServerErr)
			return
		}

		type userResp struct {
			Id       string `json:"id"`
			Name     string `json:"name"`
			Email    string `json:"email"`
			Role     string `json:"role"`
			FamilyId string `json:"family_id"`
		}
		w.WriteHeader(http.StatusOK)

		var familyId string
		if user.FamilyId.Valid {
			familyId = user.FamilyId.String
		} else {
			familyId = "Not found"
		}
		render.JSON(
			w, r, &userResp{
				Id:       user.Id,
				Name:     user.Name,
				Email:    user.Email,
				Role:     user.Role,
				FamilyId: familyId,
			},
		)
	}
}
