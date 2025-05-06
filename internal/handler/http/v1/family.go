package v1

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"family-flow-app/internal/entity"
	"family-flow-app/pkg/response"
	"github.com/go-chi/chi/v5"

	"family-flow-app/internal/service"
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
			r.Post("/members", u.getMembers(ctx, log))
			r.Get("/{familyId}", u.getByFamilyId(ctx, log))
			r.Post("/invite", u.inviteMember(ctx, log))
		},
	)
}

// createFamilyInput структура для создания семьи
type createFamilyInput struct {
	Name string `json:"name" validate:"required"`
}

// @Summary Create family
// @Description Create family
// @Tags family
// @Accept json
// @Produce json
// @Param name body string true "Name"
// @Success 201 {string} string "Family created"
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /family [post]
func (u *FamilyRoutes) create(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input createFamilyInput
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

		_, err = u.familyService.Create(
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
	Role      string `json:"role" validate:"required"`
}

// @Summary Add member to family
// @Description Add member to family
// @Tags family
// @Accept json
// @Produce json
// @Param email_user body string true "Email user"
// @Param family_id body string true "Family ID"
// @Success 200 {string} string "Member added to family"
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /family/add [post]
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

		err = u.familyService.AddMember(
			ctx, log, service.AddMemberToFamilyInput{
				FamilyId:  input.FamilyId,
				UserEmail: input.EmailUser,
				Role:      input.Role,
			},
		)

		if errors.Is(err, service.ErrUserNotFound) {
			var family entity.Family

			if family, err = u.familyService.GetFamilyByUserID(ctx, log, input.FamilyId); err != nil {
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
			return
		} else if err != nil {
			response.NewError(w, r, log, err, http.StatusInternalServerError, "Failed to add member to family")
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, "Member added to family")
	}
}

type inputGetMembers struct {
	FamilyId string `json:"family_id"`
}

// @Summary Get members
// @Description Get members
// @Tags family
// @Accept json
// @Produce json
// @Param familyId body string true "Family ID"
// @Success 200 {object} []entity.User
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /family/members [get]
func (u *FamilyRoutes) getMembers(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input inputGetMembers
		var err error

		_, err = GetCurrentUserFromContext(r.Context())
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

		var users []entity.User
		if users, err = u.familyService.GetByFamilyID(ctx, log, input.FamilyId); err != nil {
			response.NewError(w, r, log, err, http.StatusInternalServerError, "Failed to get family members")
			return
		}

		output := formatOutput{Users: users}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, output)
	}

	// присоединится к семье
	// просто по идентификатору семьи
}

type formatOutput struct {
	Users []entity.User `json:"list"`
}

func formatUsers(users []entity.User) []map[string]interface{} {
	formattedUsers := make([]map[string]interface{}, len(users))
	for i, user := range users {

		formattedUsers[i] = map[string]interface{}{
			"id":        user.Id,
			"name":      user.Name,
			"email":     user.Email,
			"role":      user.Role,
			"family_id": user.FamilyId.String,
			"latitude":  user.Latitude,
			"longitude": user.Longitude,
		}
	}
	return formattedUsers
}

// @Summary Get family by ID
// @Description Get family by ID
// @Tags family
// @Accept json
// @Produce json
// @Param familyId body string true "Family ID"
// @Success 200 {object} entity.Family
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /family/{familyId} [get]
func (u *FamilyRoutes) getByFamilyId(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		_, err = GetCurrentUserFromContext(r.Context())
		if err != nil {
			response.NewError(w, r, log, err, http.StatusUnauthorized, "Failed to get current user")
			return
		}

		familyId := chi.URLParam(r, "familyId")
		if familyId == "" {
			response.NewError(w, r, log, errors.New("missing familyId"), http.StatusBadRequest, "Missing familyId")
			return
		}
		if err = validator.New().Var(familyId, "required,uuid"); err != nil {
			response.NewValidateError(w, r, log, http.StatusBadRequest, MsgInvalidReq, err)
			return
		}

		var family entity.Family
		if family, err = u.familyService.GetByID(ctx, log, familyId); err != nil {
			response.NewError(w, r, log, err, http.StatusInternalServerError, "Failed to get family members")
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, family)
	}
}

// @Summary Invite member to family
// @Description Invite member to family
// @Tags family
// @Accept json
// @Produce json
// @Param email_user body string true "Email user"
// @Param family_id body string true "Family ID"
// @Success 200 {string} string "Member added to family"
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /family/invite [post]
func (u *FamilyRoutes) inviteMember(ctx context.Context, log *slog.Logger) http.HandlerFunc {
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

		var family entity.Family

		if family, err = u.familyService.GetFamilyByUserID(ctx, log, input.FamilyId); err != nil {
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
		return

	}
}
