package v1

import (
	"context"
	"log/slog"
	"net/http"

	"family-flow-app/internal/service"
	"family-flow-app/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
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
			r.Get("/", u.get(ctx, log))
		},
	)
}

type userResponse struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	FamilyId string `json:"family_id"`
}

// @Summary Get user info
// @Description Get user info
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} userResponse
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /user [get]
func (u *UserRoutes) get(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := GetCurrentUserFromContext(r.Context())
		if err != nil {
			response.NewError(w, r, log, err, http.StatusUnauthorized, "Failed to get current user")
			return
		}

		w.WriteHeader(http.StatusOK)

		var familyId string
		if user.FamilyId.Valid {
			familyId = user.FamilyId.String
		} else {
			familyId = ""
		}
		render.JSON(
			w, r, &userResponse{
				Id:       user.Id,
				Name:     user.Name,
				Email:    user.Email,
				Role:     user.Role,
				FamilyId: familyId,
			},
		)
	}
}
