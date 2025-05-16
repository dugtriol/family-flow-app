package v1

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"time"

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
	fileService service.File
}

func NewUserRoutes(
	ctx context.Context, log *slog.Logger, route chi.Router, userService service.User, fileService service.File,
) {
	u := UserRoutes{userService: userService, fileService: fileService}
	route.Route(
		userString, func(r chi.Router) {
			r.Get("/", u.get(ctx, log))
			r.Put("/", u.update(ctx, log))
			r.Put("/family_id", u.resetFamilyId(ctx, log))
			r.Put("/location", u.updateLocation(ctx, log))
		},
	)
}

// type userResponse struct {
// 	Id        string    `json:"id"`
// 	Name      string    `json:"name"`
// 	Email     string    `json:"email"`
// 	Role      string    `json:"role"`
// 	FamilyId  string    `json:"family_id"`
// 	Latitude  string    `json:"latitude" swaggerignore:"true"`
// 	Longitude string    `json:"longitude" swaggerignore:"true"`
// 	Gender    string    `json:"gender"`
// 	Point     int       `json:"point"`
// 	BirthDate time.Time `json:"birth_date"`
// 	Avatar    string    `json:"avatar"`
// }

// @Summary Get user info
// @Description Get user info
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} entity.User
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

		render.JSON(
			w, r, user,
		)
	}
}

type UpdateUserInput struct {
	Name      string                `json:"name"`
	Email     string                `json:"email" validate:"email"`
	Role      string                `json:"role" validate:"oneof=Parent Child Unknown"`
	Gender    string                `json:"gender" validate:"oneof=Male Female Unknown"`
	BirthDate string                `json:"birth_date"`
	Avatar    *multipart.FileHeader `json:"avatar"` // Поле для загружаемого файла
	AvatarURL string                `json:"avatar_url"`
}

// @Summary Update user info
// @Description Update user info
// @Tags user
// @Accept json
// @Produce json
// @Param input body UpdateUserInput true "Update user info"
// @Success 200 {object} entity.User
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /user [put]
func (u *UserRoutes) update(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получаем текущего пользователя из контекста
		user, err := GetCurrentUserFromContext(r.Context())
		if err != nil {
			response.NewError(w, r, log, err, http.StatusUnauthorized, "Failed to get current user")
			return
		}

		// Парсим multipart/form-data
		if err := r.ParseMultipartForm(10 << 20); err != nil { // Ограничение на размер файла: 10 MB
			response.NewError(w, r, log, err, http.StatusBadRequest, "Failed to parse form data")
			return
		}

		// Получаем данные из формы
		input := UpdateUserInput{
			Name:      r.FormValue("name"),
			Email:     r.FormValue("email"),
			Role:      r.FormValue("role"),
			Gender:    r.FormValue("gender"),
			BirthDate: r.FormValue("birth_date"),
		}

		// Преобразуем строку BirthDate в sql.NullTime
		log.Info("Parsing birth date", slog.String("birth_date", input.BirthDate))
		var birthDate sql.NullTime
		if input.BirthDate != "" {
			parsedDate, err := time.Parse("2006-01-02", input.BirthDate)
			if err != nil {
				response.NewError(w, r, log, err, http.StatusBadRequest, "Invalid date format. Use YYYY-MM-DD.")
				return
			}
			birthDate = sql.NullTime{Time: parsedDate, Valid: true}
		}

		var avatar sql.NullString

		avatar = sql.NullString{String: r.FormValue("avatar_url"), Valid: true}

		if avatar.String == "empty" {
			// Получаем файл аватара
			file, fileHeader, err := r.FormFile("avatar")
			if err != nil && err != http.ErrMissingFile {
				response.NewError(w, r, log, err, http.StatusBadRequest, "Failed to get avatar file")
				return
			}
			defer func() {
				if file != nil {
					file.Close()
				}
			}()
			input.Avatar = fileHeader

			// Загружаем аватар в облако, если он передан
			if input.Avatar != nil {
				fileBytes, err := io.ReadAll(file)
				if err != nil {
					response.NewError(w, r, log, err, http.StatusInternalServerError, "Failed to read avatar file")
					return
				}

				fileInput := service.FileUploadInput{
					FileName: input.Avatar.Filename,
					FileBody: fileBytes,
				}

				avatarPath, err := u.fileService.Upload(ctx, log, fileInput)
				if err != nil {
					response.NewError(w, r, log, err, http.StatusInternalServerError, "Failed to upload avatar")
					return
				}

				avatarURL := u.fileService.BuildImageURL(avatarPath)
				avatar = sql.NullString{String: avatarURL, Valid: true}
			}
		}

		err = u.userService.Update(
			ctx, log, service.UpdateUserInput{
				ID:        user.Id,
				Name:      input.Name,
				Email:     input.Email,
				Role:      input.Role,
				Gender:    input.Gender,
				BirthDate: birthDate,
				Avatar:    avatar,
			},
		)
		if err != nil {
			response.NewError(w, r, log, err, http.StatusInternalServerError, "Failed to update user")
			return
		}

		// Формируем успешный ответ
		w.WriteHeader(http.StatusOK)
		render.JSON(
			w,
			r,
			user,
		)
	}
}

type ResetFamilyIdInput struct {
	ID       string `json:"id"`
	FamilyId string `json:"family_id"`
}

// @Summary Reset user family id
// @Description Reset user family id
// @Tags user
// @Accept json
// @Produce json
// @Param input body ResetFamilyIdInput true "Reset user family id"
// @Success 200 {object} string
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /user/family_id [put]
func (u *UserRoutes) resetFamilyId(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := GetCurrentUserFromContext(r.Context())
		if err != nil {
			response.NewError(w, r, log, err, http.StatusUnauthorized, "Failed to get current user")
			return
		}

		var input ResetFamilyIdInput
		if err = render.DecodeJSON(r.Body, &input); err != nil {
			response.NewError(w, r, log, err, http.StatusBadRequest, MsgFailedParsing)
			return
		}
		if err = validator.New().Struct(input); err != nil {
			response.NewValidateError(w, r, log, http.StatusBadRequest, MsgInvalidReq, err)
			return
		}

		err = u.userService.ResetFamilyID(ctx, log, input.ID)
		if err != nil {
			if errors.Is(err, service.ErrInsufficientPermissions) {
				response.NewError(
					w,
					r,
					log,
					err,
					http.StatusForbidden,
					"Insufficient permissions to reset family ID",
				)
				return
			}
			response.NewError(
				w,
				r,
				log,
				err,
				http.StatusInternalServerError,
				"Failed to reset user family id",
			)
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, "Family ID reset successfully")
	}
}

type UpdateLocationInput struct {
	Latitude  float64 `json:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`
}

// @Summary Update user location
// @Description Update user location
// @Tags user
// @Accept json
// @Produce json
// @Param input body UpdateLocationInput true "Update user location"
// @Success 200 {string} string "Location updated successfully"
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /user/location [put]
func (u *UserRoutes) updateLocation(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := GetCurrentUserFromContext(r.Context())
		if err != nil {
			response.NewError(w, r, log, err, http.StatusUnauthorized, "Failed to get current user")
			return
		}

		var input UpdateLocationInput
		if err = render.DecodeJSON(r.Body, &input); err != nil {
			response.NewError(w, r, log, err, http.StatusBadRequest, "Failed to parse request")
			return
		}
		if err = validator.New().Struct(input); err != nil {
			response.NewValidateError(w, r, log, http.StatusBadRequest, "Invalid request", err)
			return
		}

		err = u.userService.UpdateLocation(
			ctx, log, service.UpdateLocationInput{
				UserID:    user.Id,
				Latitude:  input.Latitude,
				Longitude: input.Longitude,
			},
		)
		if err != nil {
			response.NewError(w, r, log, err, http.StatusInternalServerError, "Failed to update location")
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, "Location updated successfully")
	}
}
