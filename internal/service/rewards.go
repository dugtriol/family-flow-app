package service

import (
	"context"
	"fmt"
	"log/slog"

	"family-flow-app/internal/entity"
	"family-flow-app/internal/repo"
)

type RewardsService struct {
	rewardsRepo repo.Rewards
	userRepo    repo.User
}

func NewRewardsService(rewardsRepo repo.Rewards, userRepo repo.User) *RewardsService {
	return &RewardsService{
		rewardsRepo: rewardsRepo,
		userRepo:    userRepo,
	}
}

// CreateReward создает новое вознаграждение
func (s *RewardsService) Create(ctx context.Context, log *slog.Logger, input entity.Reward) (string, error) {
	log.Info("Service - RewardsService - CreateReward", "familyID", input.FamilyID, "title", input.Title)

	id, err := s.rewardsRepo.Create(ctx, input)
	if err != nil {
		log.Error("Service - RewardsService - CreateReward - Failed to create reward", "error", err)
		return "", fmt.Errorf("failed to create reward: %w", err)
	}

	log.Info("Service - RewardsService - CreateReward - Reward created successfully", "rewardID", id)
	return id, nil
}

// GetRewardsByFamilyID возвращает список вознаграждений для семьи
func (s *RewardsService) GetRewardsByFamilyID(ctx context.Context, log *slog.Logger, familyID string) (
	[]entity.Reward, error,
) {
	log.Info("Service - RewardsService - GetRewardsByFamilyID", "familyID", familyID)

	rewards, err := s.rewardsRepo.GetByFamilyID(ctx, familyID)
	if err != nil {
		log.Error("Service - RewardsService - GetRewardsByFamilyID - Failed to get rewards", "error", err)
		return nil, fmt.Errorf("failed to get rewards: %w", err)
	}

	log.Info("Service - RewardsService - GetRewardsByFamilyID - Rewards retrieved successfully", "count", len(rewards))
	return rewards, nil
}

// AddPoints добавляет очки пользователю
func (s *RewardsService) AddPoints(ctx context.Context, log *slog.Logger, userID string, points int) error {
	log.Info("Service - RewardsService - AddPoints", "userID", userID, "points", points)

	err := s.rewardsRepo.AddPoints(ctx, userID, points)
	if err != nil {
		log.Error("Service - RewardsService - AddPoints - Failed to add points", "error", err)
		return fmt.Errorf("failed to add points: %w", err)
	}

	log.Info("Service - RewardsService - AddPoints - Points added successfully")
	return nil
}

// SubtractPoints списывает очки у пользователя
func (s *RewardsService) SubtractPoints(ctx context.Context, log *slog.Logger, userID string, points int) error {
	log.Info("Service - RewardsService - SubtractPoints", "userID", userID, "points", points)

	err := s.rewardsRepo.SubtractPoints(ctx, userID, points)
	if err != nil {
		log.Error("Service - RewardsService - SubtractPoints - Failed to subtract points", "error", err)
		return fmt.Errorf("failed to subtract points: %w", err)
	}

	log.Info("Service - RewardsService - SubtractPoints - Points subtracted successfully")
	return nil
}

// GetPoints возвращает количество очков пользователя
func (s *RewardsService) GetPoints(ctx context.Context, log *slog.Logger, userID string) (int, error) {
	log.Info("Service - RewardsService - GetPoints", "userID", userID)

	points, err := s.rewardsRepo.GetPoints(ctx, userID)
	if err != nil {
		log.Error("Service - RewardsService - GetPoints - Failed to get points", "error", err)
		return 0, fmt.Errorf("failed to get points: %w", err)
	}

	log.Info("Service - RewardsService - GetPoints - Points retrieved successfully", "points", points)
	return points, nil
}

// RedeemReward обменивает очки на вознаграждение
func (s *RewardsService) Redeem(ctx context.Context, log *slog.Logger, userID, rewardID string) error {
	log.Info("Service - RewardsService - RedeemReward", "userID", userID, "rewardID", rewardID)

	// Получаем информацию о вознаграждении
	reward, err := s.rewardsRepo.GetByID(ctx, rewardID)
	if err != nil {
		log.Error("Service - RewardsService - RedeemReward - Failed to get reward", "error", err)
		return fmt.Errorf("failed to get reward: %w", err)
	}

	// Проверяем, достаточно ли очков у пользователя
	points, err := s.rewardsRepo.GetPoints(ctx, userID)
	if err != nil {
		log.Error("Service - RewardsService - RedeemReward - Failed to get points", "error", err)
		return fmt.Errorf("failed to get points: %w", err)
	}

	if points < reward.Cost {
		log.Error(
			"Service - RewardsService - RedeemReward - Not enough points",
			"userID",
			userID,
			"points",
			points,
			"cost",
			reward.Cost,
		)
		return fmt.Errorf("not enough points to redeem reward")
	}

	// Списываем очки
	err = s.rewardsRepo.SubtractPoints(ctx, userID, reward.Cost)
	if err != nil {
		log.Error("Service - RewardsService - RedeemReward - Failed to subtract points", "error", err)
		return fmt.Errorf("failed to subtract points: %w", err)
	}

	// Добавляем запись об обмене
	err = s.rewardsRepo.Redeem(ctx, userID, rewardID)
	if err != nil {
		log.Error("Service - RewardsService - RedeemReward - Failed to redeem reward", "error", err)
		return fmt.Errorf("failed to redeem reward: %w", err)
	}

	log.Info("Service - RewardsService - RedeemReward - Reward redeemed successfully")
	return nil
}

// GetRedemptionsByUserID возвращает список вознаграждений, которые пользователь обменял
func (s *RewardsService) GetRedemptionsByUserID(
	ctx context.Context, log *slog.Logger, userID string,
) ([]entity.RewardRedemption, error) {
	log.Info("Service - RewardsService - GetRedemptionsByUserID", "userID", userID)

	redemptions, err := s.rewardsRepo.GetRedemptionsByUserID(ctx, userID)
	if err != nil {
		log.Error("Service - RewardsService - GetRedemptionsByUserID - Failed to get redemptions", "error", err)
		return nil, fmt.Errorf("failed to get redemptions: %w", err)
	}

	log.Info(
		"Service - RewardsService - GetRedemptionsByUserID - Redemptions retrieved successfully",
		"count",
		len(redemptions),
	)
	return redemptions, nil
}

// GetByID 
// возвращает информацию о вознаграждении по его ID
func (s *RewardsService) GetByID(ctx context.Context, log *slog.Logger, id string) (entity.Reward, error) {
	log.Info("Service - RewardsService - GetByID", "id", id)

	reward, err := s.rewardsRepo.GetByID(ctx, id)
	if err != nil {
		log.Error("Service - RewardsService - GetByID - Failed to get reward", "error", err)
		return entity.Reward{}, fmt.Errorf("failed to get reward: %w", err)
	}
	log.Info("Service - RewardsService - GetByID - Reward retrieved successfully", "reward", reward)
	return reward, nil
}
