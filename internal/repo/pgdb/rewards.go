package pgdb

import (
	"context"
	"fmt"

	"family-flow-app/internal/entity"
	"family-flow-app/pkg/postgres"

	"github.com/Masterminds/squirrel"
)

const (
	rewardsTable = "rewards"
	//userRewardsTable       = "user_rewards"
	rewardRedemptionsTable = "reward_redemptions"
)

type RewardsRepo struct {
	*postgres.Database
}

func NewRewardsRepo(db *postgres.Database) *RewardsRepo {
	return &RewardsRepo{db}
}

// CreateReward создает новое вознаграждение
func (r *RewardsRepo) Create(ctx context.Context, reward entity.Reward) (string, error) {
	sql, args, _ := r.Builder.Insert(rewardsTable).Columns(
		"family_id",
		"title",
		"description",
		"cost",
	).Values(
		reward.FamilyID,
		reward.Title,
		reward.Description,
		reward.Cost,
	).Suffix("RETURNING id").ToSql()

	var id string
	err := r.Cluster.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("failed to create reward: %w", err)
	}
	return id, nil
}

// GetRewardsByFamilyID возвращает список вознаграждений для семьи
func (r *RewardsRepo) GetByFamilyID(ctx context.Context, familyID string) ([]entity.Reward, error) {
	sql, args, _ := r.Builder.Select(
		"id",
		"family_id",
		"title",
		"description",
		"cost",
		"created_at",
		"updated_at",
	).From(rewardsTable).Where(
		squirrel.Eq{"family_id": familyID},
	).ToSql()

	rows, err := r.Cluster.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get rewards: %w", err)
	}
	defer rows.Close()

	var rewards []entity.Reward
	for rows.Next() {
		var reward entity.Reward
		err := rows.Scan(
			&reward.ID,
			&reward.FamilyID,
			&reward.Title,
			&reward.Description,
			&reward.Cost,
			&reward.CreatedAt,
			&reward.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan reward: %w", err)
		}
		rewards = append(rewards, reward)
	}
	return rewards, nil
}

// AddPoints добавляет очки пользователю
func (r *RewardsRepo) AddPoints(ctx context.Context, userID string, points int) error {
	sql, args, _ := r.Builder.Update(userTable).
		Set("point", squirrel.Expr("point + ?", points)).
		Where(squirrel.Eq{"id": userID}).
		ToSql()

	_, err := r.Cluster.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to add points: %w", err)
	}
	return nil
}

// SubtractPoints списывает очки у пользователя
func (r *RewardsRepo) SubtractPoints(ctx context.Context, userID string, points int) error {
	sql, args, _ := r.Builder.Update(userTable).
		Set("point", squirrel.Expr("point - ?", points)).
		Where(squirrel.Eq{"id": userID}).
		ToSql()

	_, err := r.Cluster.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to subtract points: %w", err)
	}
	return nil
}

// GetPoints возвращает количество очков пользователя
func (r *RewardsRepo) GetPoints(ctx context.Context, userID string) (int, error) {
	sql, args, _ := r.Builder.Select("point").
		From(userTable).
		Where(squirrel.Eq{"id": userID}).
		ToSql()

	var points int
	err := r.Cluster.QueryRow(ctx, sql, args...).Scan(&points)
	if err != nil {
		return 0, fmt.Errorf("failed to get points: %w", err)
	}
	return points, nil
}

// RedeemReward обменивает очки на вознаграждение
func (r *RewardsRepo) Redeem(ctx context.Context, userID, rewardID string) error {
	sql, args, _ := r.Builder.Insert(rewardRedemptionsTable).Columns(
		"user_id",
		"reward_id",
	).Values(
		userID,
		rewardID,
	).ToSql()

	_, err := r.Cluster.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to redeem reward: %w", err)
	}
	return nil
}

// GetRedemptionsByUserID возвращает список вознаграждений, которые пользователь обменял
func (r *RewardsRepo) GetRedemptionsByUserID(ctx context.Context, userID string) ([]entity.RewardRedemption, error) {
	sql, args, _ := r.Builder.Select(
		"rr.id",
		"rr.user_id",
		"rr.reward_id",
		"rr.redeemed_at",
		"r.title",
		"r.description",
		"r.cost",
	).From(rewardRedemptionsTable + " rr").
		Join(rewardsTable + " r ON rr.reward_id = r.id").
		Where(squirrel.Eq{"rr.user_id": userID}).
		ToSql()

	rows, err := r.Cluster.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get redemptions: %w", err)
	}
	defer rows.Close()

	var redemptions []entity.RewardRedemption
	for rows.Next() {
		var redemption entity.RewardRedemption
		err := rows.Scan(
			&redemption.ID,
			&redemption.UserID,
			&redemption.RewardID,
			&redemption.RedeemedAt,
			&redemption.Reward.Title,
			&redemption.Reward.Description,
			&redemption.Reward.Cost,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan redemption: %w", err)
		}
		redemptions = append(redemptions, redemption)
	}
	return redemptions, nil
}

// get by id
func (r *RewardsRepo) GetByID(ctx context.Context, id string) (entity.Reward, error) {
	sql, args, _ := r.Builder.Select(
		"id",
		"family_id",
		"title",
		"description",
		"cost",
	).From(rewardsTable).Where(
		squirrel.Eq{"id": id},
	).ToSql()

	var reward entity.Reward
	err := r.Cluster.QueryRow(ctx, sql, args...).Scan(
		&reward.ID,
		&reward.FamilyID,
		&reward.Title,
		&reward.Description,
		&reward.Cost,
	)
	if err != nil {
		return entity.Reward{}, fmt.Errorf("failed to get reward: %w", err)
	}
	return reward, nil
}

func (r *RewardsRepo) Update(ctx context.Context, reward entity.Reward) error {
	sql, args, _ := r.Builder.Update(rewardsTable).
		Set("title", reward.Title).
		Set("description", reward.Description).
		Set("cost", reward.Cost).
		Set("updated_at", "NOW()").
		Where(squirrel.Eq{"id": reward.ID}).
		ToSql()

	_, err := r.Cluster.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to update reward: %w", err)
	}
	return nil
}

// Delete удаляет вознаграждение
func (r *RewardsRepo) Delete(ctx context.Context, id string) error {
	sql, args, _ := r.Builder.Delete(rewardsTable).Where(
		squirrel.Eq{"id": id},
	).ToSql()

	_, err := r.Cluster.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to delete reward: %w", err)
	}
	return nil
}
