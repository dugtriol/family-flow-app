package v1

import (
	"context"
	"fmt"
	"net/http"

	"family-flow-app/internal/entity"
	"family-flow-app/internal/service"
	"family-flow-app/pkg/response"

	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type RewardsRoutes struct {
	rewardsService      service.Rewards
	notificationService service.Notification
	familyService       service.Family
}

func NewRewardsRoutes(ctx context.Context, log *slog.Logger, route chi.Router, rewardsService service.Rewards, notificationService service.Notification, familyService service.Family) {
	routes := &RewardsRoutes{
		rewardsService:      rewardsService,
		notificationService: notificationService,
		familyService:       familyService,
	}

	route.Route(
		"/rewards", func(r chi.Router) {
			r.Post("/", routes.createReward(ctx, log)) // Создать вознаграждение
			r.Get(
				"/",
				routes.getRewardsByFamilyID(ctx, log),
			) // Получить список вознаграждений семьи
			r.Get("/points", routes.getPoints(ctx, log)) // Получить очки пользователя
			r.Post(
				"/{rewardID}/redeem",
				routes.redeemReward(ctx, log),
			) // Обменять очки на вознаграждение
			r.Get(
				"/redemptions",
				routes.getRedemptionsByUserID(ctx, log),
			) // Получить список обменов пользователя
			r.Put("/{rewardID}", routes.updateReward(ctx, log)) // Обновить награду
			r.Get(
				"/redemptions/{userID}",
				routes.getRedemptionsByUserIDParam(ctx, log),
			) // Получить список обменов для указанного пользователя
		},
	)
}

type RewardCreateInput struct {
	FamilyID    string `json:"family_id" validate:"required"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" `
	Cost        int    `json:"cost" validate:"required"`
}

// createReward создает новое вознаграждение
func (r *RewardsRoutes) createReward(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		user, err := GetCurrentUserFromContext(req.Context())
		if err != nil {
			response.NewError(w, req, log, err, http.StatusUnauthorized, ErrNoUserInContextMsg)
			return
		}

		var input RewardCreateInput
		if err := render.DecodeJSON(req.Body, &input); err != nil {
			response.NewError(w, req, log, err, http.StatusBadRequest, "Invalid request payload")
			return
		}

		// Устанавливаем familyID из контекста пользователя
		input.FamilyID = user.FamilyId.String

		rewardID, err := r.rewardsService.Create(
			ctx, log, entity.Reward{
				FamilyID:    input.FamilyID,
				Title:       input.Title,
				Description: input.Description,
				Cost:        input.Cost,
			},
		)
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

		// Получаем информацию о награде
		reward, err := r.rewardsService.GetByID(ctx, log, rewardID)
		if err != nil {
			response.NewError(w, req, log, err, http.StatusInternalServerError, "Failed to get reward details")
			return
		}

		// Получаем информацию о семье
		family, err := r.familyService.GetByFamilyID(ctx, log, user.FamilyId.String)
		if err != nil {
			response.NewError(w, req, log, err, http.StatusInternalServerError, "Failed to get family details")
			return
		}
		for _, member := range family {
			if member.Role == "Parent" {
				// Отправляем уведомление всем членам семьи
				err = r.notificationService.SendNotification(
					ctx, log, service.NotificationCreateInput{
						UserID: member.Id,
						Title:  "Награда получена",
						Body:   fmt.Sprintf("Пользователь '%s' обменял очки на награду '%s'", user.Name, reward.Title),
					},
				)
				if err != nil {
					log.Error("Failed to send notification: %v", err)
				}
			}
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

type RewardUpdateInput struct {
	FamilyID    string `json:"family_id" validate:"required"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" `
	Cost        int    `json:"cost" validate:"required"`
}

// updateReward обновляет существующую награду
func (r *RewardsRoutes) updateReward(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		_, err := GetCurrentUserFromContext(req.Context())
		if err != nil {
			response.NewError(w, req, log, err, http.StatusUnauthorized, ErrNoUserInContextMsg)
			return
		}

		rewardID := chi.URLParam(req, "rewardID")
		if rewardID == "" {
			response.NewError(w, req, log, nil, http.StatusBadRequest, "Reward ID is required")
			return
		}

		var input RewardUpdateInput
		if err := render.DecodeJSON(req.Body, &input); err != nil {
			response.NewError(w, req, log, err, http.StatusBadRequest, "Invalid request payload")
			return
		}
		log.Info("RewardUpdateInput - input: %v", input)

		// // Проверяем, принадлежит ли награда семье пользователя
		// reward, err := r.rewardsService.GetByID(ctx, log, rewardID)
		// if err != nil {
		//     response.NewError(w, req, log, err, http.StatusNotFound, "Reward not found")
		//     return
		// }
		// if reward.FamilyID != user.FamilyId.String {
		//     response.NewError(w, req, log, nil, http.StatusForbidden, "You do not have permission to update this reward")
		//     return
		// }

		// Обновляем награду
		err = r.rewardsService.Update(
			ctx, log, entity.Reward{
				ID:          rewardID,
				FamilyID:    input.FamilyID,
				Title:       input.Title,
				Description: input.Description,
				Cost:        input.Cost,
			},
		)
		if err != nil {
			response.NewError(w, req, log, err, http.StatusInternalServerError, "Failed to update reward")
			return
		}

		render.JSON(w, req, map[string]string{"message": "Reward updated successfully"})
	}
}

// @Summary Get redemptions by user ID
// @Description Get redemptions by user ID
// @Tags rewards
// @Accept json
// @Produce json
// @Param userID path string true "User ID"
// @Success 200 {array} entity.RewardRedemption
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /rewards/redemptions/{userID} [get]
// getRedemptionsByUserIDParam возвращает список вознаграждений, которые пользователь обменял по userID

func (r *RewardsRoutes) getRedemptionsByUserIDParam(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// Получаем userID из параметров запроса
		userID := chi.URLParam(req, "userID")
		if userID == "" {
			response.NewError(w, req, log, nil, http.StatusBadRequest, "User ID is required")
			return
		}

		// Получаем список обменов для указанного пользователя
		redemptions, err := r.rewardsService.GetRedemptionsByUserID(ctx, log, userID)
		if err != nil {
			response.NewError(w, req, log, err, http.StatusInternalServerError, "Failed to get redemptions")
			return
		}

		render.JSON(w, req, redemptions)
	}
}
