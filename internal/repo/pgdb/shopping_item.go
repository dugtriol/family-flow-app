package pgdb

import (
	"context"
	"log/slog"

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
	sql, args, _ := r.Builder.Update(shoppingTable).Set(
		"title", item.Title,
	).Set(
		"description", item.Description,
	).Set(
		"status", item.Status,
	).Set(
		"visibility", item.Visibility,
	).Where("id = ?", item.ID).ToSql()

	_, err := r.Cluster.Exec(ctx, sql, args...)
	return err
}

// получить списки visibility - public family_id
func (r *ShoppingRepo) GetPublicByFamilyID(
	ctx context.Context, log *slog.Logger, familyID string,
) ([]entity.ShoppingItem, error) {
	log.Info("ShoppingRepo - GetPublicByFamilyID")
	sql, args, _ := r.Builder.Select(
		"id",
		"family_id",
		"title",
		"description",
		"status",
		"visibility",
		"created_by",
		"created_at",
	).From(shoppingTable).Where("family_id = ? AND visibility = 'Public'", familyID).ToSql()

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
			&item.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

type InputGetPrivateByCreatedBy struct {
}

// получить списки visibility - private created_by
func (r *ShoppingRepo) GetPrivateByCreatedBy(
	ctx context.Context, log *slog.Logger, CreatedBy string,
) ([]entity.ShoppingItem, error) {
	log.Info("ShoppingRepo - GetPrivateByCreatedBy")
	sql, args, _ := r.Builder.Select("*").From(shoppingTable).Where(
		"created_by = ? AND visibility = 'Private'",
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
			&item.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}
