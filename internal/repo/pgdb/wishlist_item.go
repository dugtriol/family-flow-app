package pgdb

import (
	"context"
	"log/slog"
	"time"

	"family-flow-app/internal/entity"
	"family-flow-app/pkg/postgres"
)

const (
	wishlistTable = "wishlist_items"
)

type WishlistRepo struct {
	*postgres.Database
}

func NewWishlistRepo(db *postgres.Database) *WishlistRepo {
	return &WishlistRepo{db}
}

func (r *WishlistRepo) Create(ctx context.Context, log *slog.Logger, item entity.WishlistItem) (string, error) {
	log.Info("WishlistRepo - Create")
	sql, args, _ := r.Builder.Insert(wishlistTable).Columns(
		"name",
		"description",
		"link",
		"created_by",
	).Values(
		item.Name,
		item.Description,
		item.Link,
		item.CreatedBy,
	).Suffix("RETURNING id").ToSql()

	var id string
	err := r.Cluster.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *WishlistRepo) Delete(ctx context.Context, log *slog.Logger, id string) error {
	log.Info("WishlistRepo - Delete")
	sql, args, _ := r.Builder.Delete(wishlistTable).Where("id = ?", id).ToSql()
	_, err := r.Cluster.Exec(ctx, sql, args...)
	return err
}

func (r *WishlistRepo) Update(ctx context.Context, log *slog.Logger, item entity.WishlistItem) error {
	log.Info("WishlistRepo - Update")
	sql, args, _ := r.Builder.Update(wishlistTable).Set(
		"name", item.Name,
	).Set(
		"description", item.Description,
	).Set(
		"link", item.Link,
	).Set(
		"status", item.Status,
	).Set("is_archived", item.IsArchived).
		Set("updated_at", item.UpdatedAt).
		Where("id = ?", item.ID).ToSql()
	_, err := r.Cluster.Exec(ctx, sql, args...)
	return err
}

// получить список юзера
func (r *WishlistRepo) GetByUserID(ctx context.Context, log *slog.Logger, userID string) (
	[]entity.WishlistItem, error,
) {
	log.Info("WishlistRepo - GetByUserID")
	sql, args, _ := r.Builder.Select("*").From(wishlistTable).Where(
		"created_by = ? AND is_archived = false",
		userID,
	).ToSql()

	rows, err := r.Cluster.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []entity.WishlistItem
	for rows.Next() {
		var item entity.WishlistItem
		err = rows.Scan(
			&item.ID,
			&item.Name,
			&item.Description,
			&item.Link,
			&item.Status,
			&item.CreatedBy,
			&item.ReservedBy,
			&item.IsArchived,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.Photo,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

// update reserved by
func (r *WishlistRepo) UpdateReservedBy(ctx context.Context, log *slog.Logger, id, reservedBy string) error {
	log.Info("WishlistRepo - UpdateReservedBy")
	sql, args, _ := r.Builder.Update(wishlistTable).Set(
		"reserved_by", reservedBy,
	).Set("Status", "Reserved").
		Set("updated_at", time.Now().Add(time.Hour*3)).
		Where("id = ?", id).ToSql()
	_, err := r.Cluster.Exec(ctx, sql, args...)
	return err
}

// update reserved by
func (r *WishlistRepo) CancelUpdateReservedBy(ctx context.Context, log *slog.Logger, id string) error {
	log.Info("WishlistRepo - UpdateReservedBy")
	sql, args, _ := r.Builder.Update(wishlistTable).Set(
		"reserved_by", nil,
	).Set("Status", "Active").
		Set("updated_at", time.Now().Add(time.Hour*3)).
		Where("id = ?", id).ToSql()
	_, err := r.Cluster.Exec(ctx, sql, args...)
	return err
}

// get is_archived by user id
func (r *WishlistRepo) GetArchivedByUserID(ctx context.Context, log *slog.Logger, userID string) (
	[]entity.WishlistItem, error,
) {
	log.Info("WishlistRepo - GetArchivedByUserID")
	sql, args, _ := r.Builder.Select("*").From(wishlistTable).Where(
		"created_by = ? AND is_archived = true",
		userID,
	).ToSql()

	rows, err := r.Cluster.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []entity.WishlistItem
	for rows.Next() {
		var item entity.WishlistItem
		err = rows.Scan(
			&item.ID,
			&item.Name,
			&item.Description,
			&item.Link,
			&item.Status,
			&item.CreatedBy,
			&item.ReservedBy,
			&item.IsArchived,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.Photo,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

// get by wishlist id
func (r *WishlistRepo) GetByID(ctx context.Context, log *slog.Logger, id string) (entity.WishlistItem, error) {
	log.Info("WishlistRepo - GetByID")
	sql, args, _ := r.Builder.Select("*").From(wishlistTable).Where(
		"id = ?",
		id,
	).ToSql()

	rows, err := r.Cluster.Query(ctx, sql, args...)
	if err != nil {
		return entity.WishlistItem{}, err
	}
	defer rows.Close()

	var item entity.WishlistItem
	if rows.Next() {
		err = rows.Scan(
			&item.ID,
			&item.Name,
			&item.Description,
			&item.Link,
			&item.Status,
			&item.CreatedBy,
			&item.ReservedBy,
			&item.IsArchived,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.Photo,
		)
		if err != nil {
			return entity.WishlistItem{}, err
		}
	}
	return item, nil
}
