package pgdb

import (
	"context"
	"log/slog"
	"time"

	"family-flow-app/internal/entity"
	"family-flow-app/pkg/postgres"
)

const (
	diaryTable = "diary_items"
)

type DiaryRepo struct {
	*postgres.Database
}

func NewDiaryRepo(db *postgres.Database) *DiaryRepo {
	return &DiaryRepo{db}
}

// Создание записи в дневнике
func (r *DiaryRepo) Create(ctx context.Context, log *slog.Logger, item entity.DiaryItem) (string, error) {
	log.Info("DiaryRepo - Create")
	sql, args, _ := r.Builder.Insert(diaryTable).Columns(
		"title",
		"description",
		"emoji",
		"created_by",
	).Values(
		item.Title,
		item.Description,
		item.Emoji,
		item.CreatedBy,
	).Suffix("RETURNING id").ToSql()

	var id string
	err := r.Cluster.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

// Получение списка записей по идентификатору пользователя
func (r *DiaryRepo) GetByUserID(ctx context.Context, log *slog.Logger, userID string) ([]entity.DiaryItem, error) {
	log.Info("DiaryRepo - GetByUserID")
	sql, args, _ := r.Builder.Select("*").From(diaryTable).Where("created_by = ?", userID).ToSql()

	rows, err := r.Cluster.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []entity.DiaryItem
	for rows.Next() {
		var item entity.DiaryItem
		err = rows.Scan(
			&item.ID,
			&item.Title,
			&item.Description,
			&item.Emoji,
			&item.CreatedBy,
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

// Обновление записи в дневнике
func (r *DiaryRepo) Update(ctx context.Context, log *slog.Logger, item entity.DiaryItem) error {
	log.Info("DiaryRepo - Update")
	sql, args, _ := r.Builder.Update(diaryTable).Set(
		"title", item.Title,
	).Set(
		"description", item.Description,
	).Set(
		"emoji", item.Emoji,
	).Set(
		"updated_at", time.Now(),
	).Where("id = ?", item.ID).ToSql()

	_, err := r.Cluster.Exec(ctx, sql, args...)
	return err
}

// Удаление записи из дневника
func (r *DiaryRepo) Delete(ctx context.Context, log *slog.Logger, id string) error {
	log.Info("DiaryRepo - Delete")
	sql, args, _ := r.Builder.Delete(diaryTable).Where("id = ?", id).ToSql()
	_, err := r.Cluster.Exec(ctx, sql, args...)
	return err
}
