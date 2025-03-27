package v1

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"family-flow-app/internal/entity"
	"family-flow-app/internal/service"
	"family-flow-app/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

const (
	familyString = "/family"
)

type FamilyRoutes struct {
	familyService service.Family
	emailService  service.Email
}

func NewFamilyRoutes(
	ctx context.Context, log *slog.Logger, route chi.Router, emailService service.Email, familyService service.Family,
) {
	u := FamilyRoutes{familyService: familyService, emailService: emailService}
	route.Route(
		familyString, func(r chi.Router) {
			r.Post("/add", u.addMember(ctx, log))
			r.Post("/", u.create(ctx, log))
		},
	)
}

// createFamilyInput структура для создания семьи
type createFamilyInput struct {
	Name string `json:"name" validate:"required"`
}

// CreateFamily создать семью
func (u *FamilyRoutes) create(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input createFamilyInput
		var err error

		user, err := GetCurrentUserFromContext(r.Context())
		if err != nil {
			response.NewError(w, r, log, err, http.StatusUnauthorized, "Failed to get current user")
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

		_, err = u.familyService.CreateFamily(
			ctx, log, service.FamilyCreateInput{
				Name:          input.Name,
				CreatorUserId: user.Id,
			},
		)

		if err != nil {
			response.NewError(w, r, log, err, http.StatusInternalServerError, "Failed to create family")
			return
		}

		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, "Family created")
	}
}

type inputAddMemberToFamily struct {
	EmailUser string `json:"email_user" validate:"required,email"`
	FamilyId  string `json:"family_id" validate:"required,uuid"`
}

// AddMemberToFamily добавить пользователя в семью
func (u *FamilyRoutes) addMember(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input inputAddMemberToFamily
		var err error

		user, err := GetCurrentUserFromContext(r.Context())
		if err != nil {
			response.NewError(w, r, log, err, http.StatusUnauthorized, "Failed to get current user")
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

		err = u.familyService.AddMemberToFamily(
			ctx, log, service.AddMemberToFamilyInput{
				FamilyId: input.FamilyId,
				UserId:   input.EmailUser,
			},
		)

		if errors.Is(err, service.ErrUserNotFound) {
			var family entity.Family

			if family, err = u.familyService.GetFamilyByID(ctx, log, input.FamilyId); err != nil {
				response.NewError(w, r, log, err, http.StatusNotFound, "Family not found")
				return
			}
			if err = u.emailService.SendInvite(
				ctx, service.InputSendInvite{
					To:         []string{input.EmailUser},
					From:       user.Email,
					FromName:   user.Name,
					FamilyName: family.Name,
				},
			); err != nil {
				response.NewError(w, r, log, err, http.StatusInternalServerError, "Failed to send invite")
				return
			}
			w.WriteHeader(http.StatusOK)
			render.JSON(w, r, "Invite sent")
		} else if err != nil {
			response.NewError(w, r, log, err, http.StatusInternalServerError, "Failed to add member to family")
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, "Member added to family")
	}
}
