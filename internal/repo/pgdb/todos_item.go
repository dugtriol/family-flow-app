package pgdb

import (
	"context"
	"log/slog"

	"family-flow-app/internal/entity"
	"family-flow-app/pkg/postgres"
)

const (
	todoTable = "todo_items"
)

type TodoRepo struct {
	*postgres.Database
}

func NewTodoRepo(db *postgres.Database) *TodoRepo {
	return &TodoRepo{db}
}

func (r *TodoRepo) Create(ctx context.Context, log *slog.Logger, item entity.TodoItem) (string, error) {
	log.Info("TodoRepo - Create")
	sql, args, _ := r.Builder.Insert(todoTable).Columns(
		"family_id",
		"title",
		"description",
		"deadline",
		"assigned_to",
		"created_by",
		"point",
	).Values(
		item.FamilyID,
		item.Title,
		item.Description,
		item.Deadline,
		item.AssignedTo,
		item.CreatedBy,
		item.Point,
	).Suffix("RETURNING id").ToSql()

	var id string
	err := r.Cluster.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *TodoRepo) Delete(ctx context.Context, log *slog.Logger, id string) error {
	log.Info("TodoRepo - Delete")
	sql, args, _ := r.Builder.Delete(todoTable).Where("id = ?", id).ToSql()
	_, err := r.Cluster.Exec(ctx, sql, args...)
	return err
}

func (r *TodoRepo) Update(ctx context.Context, log *slog.Logger, item entity.TodoItem) error {
	log.Info("TodoRepo - Update")
	sql, args, _ := r.Builder.Update(todoTable).Set(
		"title", item.Title,
	).Set(
		"description", item.Description,
	).Set(
		"status", item.Status,
	).Set(
		"deadline", item.Deadline,
	).Set(
		"assigned_to", item.AssignedTo,
	).Set("point", item.Point).
		Where("id = ?", item.ID).ToSql()
	_, err := r.Cluster.Exec(ctx, sql, args...)
	return err
}

func (r *TodoRepo) getByField(ctx context.Context, log *slog.Logger, field, value string) ([]entity.TodoItem, error) {
	log.Info("TodoRepo - getByField")
	sql, args, _ := r.Builder.Select("*").From(todoTable).Where(field+" = ?", value).ToSql()

	rows, err := r.Cluster.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []entity.TodoItem
	for rows.Next() {
		var item entity.TodoItem
		if err := rows.Scan(
			&item.ID,
			&item.FamilyID,
			&item.Title,
			&item.Description,
			&item.Status,
			&item.Deadline,
			&item.AssignedTo,
			&item.CreatedBy,
			&item.IsArchived,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.Point,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

// получить списки по assigned to
func (r *TodoRepo) GetByAssignedTo(ctx context.Context, log *slog.Logger, assignedTo string) (
	[]entity.TodoItem, error,
) {
	// проверка того, что assigned_to != creator_by
	return r.getByField(ctx, log, "assigned_to", assignedTo)
}

// получиь списки по created by
func (r *TodoRepo) GetByCreatedBy(ctx context.Context, log *slog.Logger, createdBy string) ([]entity.TodoItem, error) {
	return r.getByField(ctx, log, "created_by", createdBy)
}

// get by id
func (r *TodoRepo) GetByID(ctx context.Context, log *slog.Logger, id string) (entity.TodoItem, error) {
	log.Info("TodoRepo - GetByID")
	sql, args, _ := r.Builder.Select("*").From(todoTable).Where("id = ?", id).ToSql()

	row := r.Cluster.QueryRow(ctx, sql, args...)
	var item entity.TodoItem
	if err := row.Scan(
		&item.ID,
		&item.FamilyID,
		&item.Title,
		&item.Description,
		&item.Status,
		&item.Deadline,
		&item.AssignedTo,
		&item.CreatedBy,
		&item.IsArchived,
		&item.CreatedAt,
		&item.UpdatedAt,
		&item.Point,
	); err != nil {
		return entity.TodoItem{}, err
	}

	return item, nil
}
