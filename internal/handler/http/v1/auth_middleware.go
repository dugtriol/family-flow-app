package v1

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strings"

	"family-flow-app/internal/entity"
	"family-flow-app/pkg/response"

	"family-flow-app/internal/service"
	"family-flow-app/pkg/token"
	"github.com/golang-jwt/jwt/v5"
)

const CurrentUserKey = "currentUser"

func AuthMiddleware(
	ctx context.Context, log *slog.Logger, service service.User,
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			var err error
			// вытаскиваем jwt, проверяем его, возвращаем id пользователя, вставляем его в context value
			header := r.Header.Get("Authorization")
			parseToken, err, done := authValidateToken(w, r, log, header)
			if done {
				return
			}
			var userID string
			if userID, err = parseToken.Claims.GetSubject(); err != nil {
				response.NewError(
					w,
					r,
					log,
					ErrInvalidToken,
					http.StatusUnauthorized,
					"AuthMiddleware: Failed to get user id",
				)
				return
			}
			var output entity.User
			if output, err = service.GetById(ctx, log, userID); err != nil {
				log.Info("AuthMiddleware - service.GetByID", "error", err.Error())
				response.NewError(
					w,
					r,
					log,
					ErrUserGet,
					http.StatusInternalServerError,
					"AuthMiddleware: Failed to get user by id",
				)
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), CurrentUserKey, output))
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func authValidateToken(w http.ResponseWriter, r *http.Request, log *slog.Logger, header string) (
	*jwt.Token, error, bool,
) {
	var err error
	arr := strings.Split(header, " ")

	if len(arr) != 2 {
		response.NewError(
			w,
			r,
			log,
			ErrInvalidToken,
			http.StatusUnauthorized,
			"AuthMiddleware: Invalid format of token",
		)
		return nil, nil, true
	}

	tokenString := arr[1]
	var parseToken *jwt.Token
	if parseToken, err = token.Check(tokenString); err != nil {
		response.NewError(
			w,
			r,
			log,
			ErrInvalidToken,
			http.StatusUnauthorized,
			"AuthMiddleware: Bad token",
		)
		return nil, nil, true
	}
	return parseToken, err, false
}

func GetCurrentUserFromContext(ctx context.Context) (*entity.User, error) {
	user, ok := ctx.Value(CurrentUserKey).(entity.User)
	log.Printf(fmt.Sprintf("GetCurrentUserFromContext - ctx.Value(CurrentUserKey).(*entity.User): %v", user))
	if !ok || user.Id == "" {
		return nil, ErrNoUserInContext
	}

	return &user, nil
}
