package v1

import (
	"context"
	"net/http"

	"family-flow-app/internal/entity"
	"family-flow-app/internal/service"
	"family-flow-app/pkg/response"

	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type RewardsRoutes struct {
	rewardsService service.Rewards
}

func NewRewardsRoutes(ctx context.Context, log *slog.Logger, route chi.Router, rewardsService service.Rewards) {
	routes := &RewardsRoutes{
		rewardsService: rewardsService,
	}

	route.Route("/rewards", func(r chi.Router) {
		r.Post("/", routes.createReward(ctx, log))                     // Создать вознаграждение
		r.Get("/", routes.getRewardsByFamilyID(ctx, log))              // Получить список вознаграждений семьи
		r.Get("/points", routes.getPoints(ctx, log))                   // Получить очки пользователя
		r.Post("/{rewardID}/redeem", routes.redeemReward(ctx, log))    // Обменять очки на вознаграждение
		r.Get("/redemptions", routes.getRedemptionsByUserID(ctx, log)) // Получить список обменов пользователя
	})
}

// createReward создает новое вознаграждение
func (r *RewardsRoutes) createReward(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		user, err := GetCurrentUserFromContext(req.Context())
		if err != nil {
			response.NewError(w, req, log, err, http.StatusUnauthorized, ErrNoUserInContextMsg)
			return
		}

		var input entity.Reward
		if err := render.DecodeJSON(req.Body, &input); err != nil {
			response.NewError(w, req, log, err, http.StatusBadRequest, "Invalid request payload")
			return
		}

		// Устанавливаем familyID из контекста пользователя
		input.FamilyID = user.FamilyId.String

		rewardID, err := r.rewardsService.Create(ctx, log, input)
		if err != nil {
			response.NewError(w, req, log, err, http.StatusInternalServerError, "Failed to create reward")
			return
		}

		render.JSON(w, req, map[string]string{"reward_id": rewardID})
	}
}

// getRewardsByFamilyID возвращает список вознаграждений для семьи
func (r *RewardsRoutes) getRewardsByFamilyID(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		user, err := GetCurrentUserFromContext(req.Context())
		if err != nil {
			response.NewError(w, req, log, err, http.StatusUnauthorized, ErrNoUserInContextMsg)
			return
		}

		rewards, err := r.rewardsService.GetRewardsByFamilyID(ctx, log, user.FamilyId.String)
		if err != nil {
			response.NewError(w, req, log, err, http.StatusInternalServerError, "Failed to get rewards")
			return
		}

		render.JSON(w, req, rewards)
	}
}

// getPoints возвращает количество очков пользователя
func (r *RewardsRoutes) getPoints(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		user, err := GetCurrentUserFromContext(req.Context())
		if err != nil {
			response.NewError(w, req, log, err, http.StatusUnauthorized, ErrNoUserInContextMsg)
			return
		}

		points, err := r.rewardsService.GetPoints(ctx, log, user.Id)
		if err != nil {
			response.NewError(w, req, log, err, http.StatusInternalServerError, "Failed to get points")
			return
		}

		render.JSON(w, req, map[string]int{"points": points})
	}
}

// redeemReward обменивает очки на вознаграждение
func (r *RewardsRoutes) redeemReward(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		user, err := GetCurrentUserFromContext(req.Context())
		if err != nil {
			response.NewError(w, req, log, err, http.StatusUnauthorized, ErrNoUserInContextMsg)
			return
		}

		rewardID := chi.URLParam(req, "rewardID")
		if rewardID == "" {
			response.NewError(w, req, log, nil, http.StatusBadRequest, "Reward ID is required")
			return
		}

		err = r.rewardsService.Redeem(ctx, log, user.Id, rewardID)
		if err != nil {
			response.NewError(w, req, log, err, http.StatusInternalServerError, "Failed to redeem reward")
			return
		}

		render.JSON(w, req, map[string]string{"message": "Reward redeemed successfully"})
	}
}

// getRedemptionsByUserID возвращает список вознаграждений, которые пользователь обменял
func (r *RewardsRoutes) getRedemptionsByUserID(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		user, err := GetCurrentUserFromContext(req.Context())
		if err != nil {
			response.NewError(w, req, log, err, http.StatusUnauthorized, ErrNoUserInContextMsg)
			return
		}

		redemptions, err := r.rewardsService.GetRedemptionsByUserID(ctx, log, user.Id)
		if err != nil {
			response.NewError(w, req, log, err, http.StatusInternalServerError, "Failed to get redemptions")
			return
		}

		render.JSON(w, req, redemptions)
	}
}
