package pgdb

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"family-flow-app/internal/entity"
	"family-flow-app/pkg/postgres"
)

const (
	shoppingTable = "shopping_items"
)

type ShoppingRepo struct {
	*postgres.Database
}

func NewShoppingRepo(db *postgres.Database) *ShoppingRepo {
	return &ShoppingRepo{db}
}

func (r *ShoppingRepo) Create(ctx context.Context, log *slog.Logger, item entity.ShoppingItem) (string, error) {
	log.Info("ShoppingRepo - Create")
	sql, args, _ := r.Builder.Insert(shoppingTable).Columns(
		"family_id",
		"title",
		"description",
		"visibility",
		"created_by",
	).Values(
		item.FamilyID,
		item.Title,
		item.Description,
		item.Visibility,
		item.CreatedBy,
	).Suffix("RETURNING id").ToSql()

	var id string
	err := r.Cluster.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *ShoppingRepo) Delete(ctx context.Context, log *slog.Logger, id string) error {
	log.Info("ShoppingRepo - Delete")
	sql, args, _ := r.Builder.Delete(shoppingTable).Where("id = ?", id).ToSql()
	_, err := r.Cluster.Exec(ctx, sql, args...)
	return err
}

func (r *ShoppingRepo) Update(ctx context.Context, log *slog.Logger, item entity.ShoppingItem) error {
	log.Info("ShoppingRepo - Update")

	log.Info(fmt.Sprintf("ShoppingRepo - Update - IsArchived - %s", item.IsArchived))
	sql, args, _ := r.Builder.Update(shoppingTable).Set(
		"title", item.Title,
	).Set(
		"description", item.Description,
	).Set(
		"status", item.Status,
	).Set(
		"visibility", item.Visibility,
	).Set("is_archived", item.IsArchived).
		Set("updated_at", item.UpdatedAt).
		Where("id = ?", item.ID).ToSql()

	_, err := r.Cluster.Exec(ctx, sql, args...)
	return err
}

// получить списки visibility - public family_id
func (r *ShoppingRepo) GetPublicByFamilyID(
	ctx context.Context, log *slog.Logger, familyID string,
) ([]entity.ShoppingItem, error) {
	log.Info("ShoppingRepo - GetPublicByFamilyID")
	sql, args, _ := r.Builder.Select(
		"*",
	).From(shoppingTable).Where("family_id = ? AND visibility = 'Public' AND is_archived = false", familyID).ToSql()

	rows, err := r.Cluster.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]entity.ShoppingItem, 0)
	for rows.Next() {
		var item entity.ShoppingItem
		err = rows.Scan(
			&item.ID,
			&item.FamilyID,
			&item.Title,
			&item.Description,
			&item.Status,
			&item.Visibility,
			&item.CreatedBy,
			&item.ReservedBy,
			&item.BuyerId,
			&item.IsArchived,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

// получить списки visibility - private created_by
func (r *ShoppingRepo) GetPrivateByCreatedBy(
	ctx context.Context, log *slog.Logger, CreatedBy string,
) ([]entity.ShoppingItem, error) {
	log.Info("ShoppingRepo - GetPrivateByCreatedBy")
	sql, args, _ := r.Builder.Select("*").From(shoppingTable).Where(
		"created_by = ? AND visibility = 'Private' AND is_archived = false",
		CreatedBy,
	).ToSql()

	rows, err := r.Cluster.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []entity.ShoppingItem
	for rows.Next() {
		var item entity.ShoppingItem
		err = rows.Scan(
			&item.ID,
			&item.FamilyID,
			&item.Title,
			&item.Description,
			&item.Status,
			&item.Visibility,
			&item.CreatedBy,
			&item.ReservedBy,
			&item.BuyerId,
			&item.IsArchived,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

// update reserved_by
func (r *ShoppingRepo) UpdateReservedBy(
	ctx context.Context, log *slog.Logger, id string, reservedBy string, updatedAt time.Time,
) error {
	log.Info("ShoppingRepo - UpdateReservedBy")
	sql, args, _ := r.Builder.Update(shoppingTable).Set(
		"reserved_by", reservedBy,
	).Set("status", "Reserved").Set("updated_at", updatedAt).Where("id = ?", id).ToSql()

	_, err := r.Cluster.Exec(ctx, sql, args...)
	return err
}

// update reserved_by
func (r *ShoppingRepo) CancelUpdateReservedBy(
	ctx context.Context, log *slog.Logger, id string, updatedAt time.Time,
) error {
	log.Info("ShoppingRepo - UpdateReservedBy")
	sql, args, _ := r.Builder.Update(shoppingTable).Set(
		"reserved_by", nil,
	).Set("status", "Active").Set("updated_at", updatedAt).Where("id = ?", id).ToSql()

	_, err := r.Cluster.Exec(ctx, sql, args...)
	return err
}

// update buyer_id
func (r *ShoppingRepo) UpdateBuyerId(
	ctx context.Context, log *slog.Logger, id string, buyerId string, updatedAt time.Time,
) error {
	log.Info("ShoppingRepo - UpdateBuyerId")
	sql, args, _ := r.Builder.Update(shoppingTable).Set(
		"buyer_id", buyerId,
	).Set("status", "Completed").
		Set("updated_at", updatedAt).
		Where("id = ?", id).ToSql()

	_, err := r.Cluster.Exec(ctx, sql, args...)
	return err
}

// get archived items by user id
func (r *ShoppingRepo) GetArchivedByUserID(
	ctx context.Context, log *slog.Logger, userID string,
) ([]entity.ShoppingItem, error) {
	log.Info("ShoppingRepo - GetArchivedByUserID")
	sql, args, _ := r.Builder.Select("*").From(shoppingTable).Where(
		"created_by = ? AND is_archived = true",
		userID,
	).ToSql()

	rows, err := r.Cluster.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []entity.ShoppingItem
	for rows.Next() {
		var item entity.ShoppingItem
		err = rows.Scan(
			&item.ID,
			&item.FamilyID,
			&item.Title,
			&item.Description,
			&item.Status,
			&item.Visibility,
			&item.CreatedBy,
			&item.ReservedBy,
			&item.BuyerId,
			&item.IsArchived,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

// get by id
func (r *ShoppingRepo) GetByID(
	ctx context.Context, log *slog.Logger, id string,
) (entity.ShoppingItem, error) {
	log.Info("ShoppingRepo - GetByID")
	sql, args, _ := r.Builder.Select("*").From(shoppingTable).Where(
		"id = ?",
		id,
	).ToSql()

	row := r.Cluster.QueryRow(ctx, sql, args...)

	var item entity.ShoppingItem
	err := row.Scan(
		&item.ID,
		&item.FamilyID,
		&item.Title,
		&item.Description,
		&item.Status,
		&item.Visibility,
		&item.CreatedBy,
		&item.ReservedBy,
		&item.BuyerId,
		&item.IsArchived,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		return entity.ShoppingItem{}, err
	}
	return item, nil
}
