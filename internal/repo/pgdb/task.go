package pgdb

import (
	"context"
	"errors"
	"fmt"
	"log"

	"family-flow-app/internal/entity"
	"family-flow-app/internal/repo/repoerrs"
	"family-flow-app/pkg/postgres"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	taskTable = "tasks"
)

type TaskRepo struct {
	*postgres.Database
}

func NewTaskRepo(db *postgres.Database) *TaskRepo {
	return &TaskRepo{db}
}

func (r *TaskRepo) Create(ctx context.Context, task entity.Task) (string, error) {
	sql, args, _ := r.Builder.Insert(taskTable).Columns(
		"title", "description", "status", "deadline", "assigned_to", "created_by", "reward",
	).Values(
		task.Title,
		task.Description,
		task.Status,
		task.Deadline,
		task.AssignedTo,
		task.CreatedBy,
		task.Reward,
	).Suffix("RETURNING id").ToSql()

	var id string
	err := r.Cluster.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return "", repoerrs.ErrAlreadyExists
			}
		}
		return "", fmt.Errorf("TaskRepo - Create - r.Cluster.QueryRow: %v", err)
	}
	return id, nil
}

func (t *TaskRepo) GetByID(ctx context.Context, id string) (entity.Task, error) {
	return t.getByField(ctx, "id", id)
}

func (t *TaskRepo) GetByAssignedTo(ctx context.Context, assignedTo string) ([]entity.Task, error) {
	return t.getByFieldList(ctx, "assigned_to", assignedTo)
}

func (t *TaskRepo) GetByCreatedBy(ctx context.Context, createdBy string) ([]entity.Task, error) {
	return t.getByFieldList(ctx, "created_by", createdBy)
}

func (t *TaskRepo) getByField(ctx context.Context, field, value string) (entity.Task, error) {
	sql, args, _ := t.Builder.
		Select("*").
		From(taskTable).
		Where(fmt.Sprintf("%v = ?", field), value).
		ToSql()
	log.Printf("TaskRepo - GetByField - sql %s args %s \n", sql, args)

	var output entity.Task
	err := t.Cluster.QueryRow(ctx, sql, args...).Scan(
		&output.Id,
		&output.Title,
		&output.Description,
		&output.Status,
		&output.Deadline,
		&output.AssignedTo,
		&output.CreatedBy,
		&output.Reward,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Task{}, repoerrs.ErrNotFound
		}
		return entity.Task{}, fmt.Errorf("TaskRepo - GetByField %s - r.Cluster.QueryRow: %v", field, err)
	}
	return output, nil
}

func (t *TaskRepo) getByFieldList(ctx context.Context, field, value string) ([]entity.Task, error) {
	sql, args, _ := t.Builder.
		Select("*").
		From(taskTable).
		Where(fmt.Sprintf("%v = ?", field), value).
		ToSql()
	log.Printf("TaskRepo - GetByFieldList - sql %s args %s \n", sql, args)

	rows, err := t.Cluster.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("TaskRepo - GetByFieldList %s - r.Cluster.Query: %v", field, err)
	}
	defer rows.Close()

	var tasks []entity.Task
	for rows.Next() {
		var task entity.Task
		err := rows.Scan(
			&task.Id,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Deadline,
			&task.AssignedTo,
			&task.CreatedBy,
			&task.Reward,
		)
		if err != nil {
			return nil, fmt.Errorf("TaskRepo - GetByFieldList - rows.Scan: %v", err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("TaskRepo - GetByFieldList - rows.Err: %v", err)
	}

	return tasks, nil
}

func (t *TaskRepo) Update(ctx context.Context, task entity.Task) error {
	sql, args, _ := t.Builder.Update(taskTable).
		Set("title", task.Title).
		Set("description", task.Description).
		Set("status", task.Status).
		Set("deadline", task.Deadline).
		Set("assigned_to", task.AssignedTo).
		Set("reward", task.Reward).
		Where("id = ?", task.Id).
		ToSql()

	_, err := t.Cluster.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("TaskRepo - Update - r.Cluster.Exec: %v", err)
	}
	return nil
}

func (t *TaskRepo) Delete(ctx context.Context, id string) error {
	sql, args, _ := t.Builder.Delete(taskTable).
		Where("id = ?", id).
		ToSql()

	_, err := t.Cluster.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("TaskRepo - Delete - r.Cluster.Exec: %v", err)
	}
	return nil
}

func (t *TaskRepo) GetByStatus(ctx context.Context, status string) ([]entity.Task, error) {
	sql, args, _ := t.Builder.
		Select("*").
		From(taskTable).
		Where("status = ?", status).
		ToSql()
	log.Printf("TaskRepo - GetByStatus - sql %s args %s \n", sql, args)

	rows, err := t.Cluster.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("TaskRepo - GetByStatus - r.Cluster.Query: %v", err)
	}
	defer rows.Close()

	var tasks []entity.Task
	for rows.Next() {
		var task entity.Task
		err := rows.Scan(
			&task.Id,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Deadline,
			&task.AssignedTo,
			&task.CreatedBy,
			&task.Reward,
		)
		if err != nil {
			return nil, fmt.Errorf("TaskRepo - GetByStatus - rows.Scan: %v", err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("TaskRepo - GetByStatus - rows.Err: %v", err)
	}

	return tasks, nil
}

func (t *TaskRepo) GetOverdueTasks(ctx context.Context) ([]entity.Task, error) {
	sql, args, _ := t.Builder.
		Select("*").
		From(taskTable).
		Where("status = ? AND deadline < NOW()", "active").
		ToSql()
	log.Printf("TaskRepo - GetOverdueTasks - sql %s args %s \n", sql, args)

	rows, err := t.Cluster.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("TaskRepo - GetOverdueTasks - r.Cluster.Query: %v", err)
	}
	defer rows.Close()

	var tasks []entity.Task
	for rows.Next() {
		var task entity.Task
		err := rows.Scan(
			&task.Id,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Deadline,
			&task.AssignedTo,
			&task.CreatedBy,
			&task.Reward,
		)
		if err != nil {
			return nil, fmt.Errorf("TaskRepo - GetOverdueTasks - rows.Scan: %v", err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("TaskRepo - GetOverdueTasks - rows.Err: %v", err)
	}

	return tasks, nil
}

func (t *TaskRepo) Complete(ctx context.Context, id string, userId string) error {
	// Проверяем, что задача назначена на текущего пользователя
	sqlCheck, argsCheck, _ := t.Builder.
		Select("assigned_to").
		From(taskTable).
		Where("id = ?", id).
		ToSql()

	var assignedTo string
	err := t.Cluster.QueryRow(ctx, sqlCheck, argsCheck...).Scan(&assignedTo)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repoerrs.ErrNotFound
		}
		return fmt.Errorf("TaskRepo - Complete - r.Cluster.QueryRow: %v", err)
	}

	if assignedTo != userId {
		return repoerrs.ErrForbidden // Задача не назначена на текущего пользователя
	}

	// Обновляем статус задачи
	sqlUpdate, argsUpdate, _ := t.Builder.
		Update(taskTable).
		Set("status", "completed").
		Where("id = ?", id).
		ToSql()

	_, err = t.Cluster.Exec(ctx, sqlUpdate, argsUpdate...)
	if err != nil {
		return fmt.Errorf("TaskRepo - Complete - r.Cluster.Exec: %v", err)
	}

	return nil
}
