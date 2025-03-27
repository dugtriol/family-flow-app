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
	userTable = "users"
)

type UserRepo struct {
	*postgres.Database
}

func NewUserRepo(db *postgres.Database) *UserRepo {
	return &UserRepo{db}
}

func (r *UserRepo) Create(ctx context.Context, user entity.User) (string, error) {
	sql, args, _ := r.Builder.Insert(userTable).Columns("name", "email", "password", "role").Values(
		user.Name,
		user.Email,
		user.Password,
		user.Role,
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
		return "", fmt.Errorf("UserRepo - Create - r.Cluster.QueryRow: %v", err)
	}
	return id, nil
}

func (u *UserRepo) GetByID(ctx context.Context, id string) (entity.User, error) {
	return u.getByField(ctx, "id", id)
}

func (u *UserRepo) GetByEmail(ctx context.Context, email string) (entity.User, error) {
	return u.getByField(ctx, "email", email)
}
func (u *UserRepo) getByField(ctx context.Context, field, value string) (entity.User, error) {
	var err error
	sql, args, _ := u.Builder.
		Select("*").
		From(userTable).
		Where(fmt.Sprintf("%v = ?", field), value).
		ToSql()
	log.Printf("UserRepo - GetByField - sql %s args %s \n", sql, args)

	var output entity.User
	err = u.Cluster.QueryRow(ctx, sql, args...).Scan(
		&output.Id,
		&output.Name,
		&output.Email,
		&output.Password,
		&output.Role,
		&output.FamilyId,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, repoerrs.ErrNotFound
		}
		return entity.User{}, fmt.Errorf("UserRepo - GetByField %s - r.Cluster.QueryRow: %v", field, err)
	}
	return output, nil
}

func (u *UserRepo) UpdateFamilyID(ctx context.Context, userID, familyID string) error {
	sql, args, _ := u.Builder.Update(userTable).Set("family_id", familyID).Where("id = ?", userID).ToSql()

	_, err := u.Cluster.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("UserRepo - UpdateFamilyID - r.Cluster.Exec: %v", err)
	}
	return nil
}

func (u *UserRepo) GetByFamilyID(ctx context.Context, familyID string) ([]entity.User, error) {
	sql, args, _ := u.Builder.
		Select("*").
		From(userTable).
		Where("family_id = ?", familyID).
		ToSql()

	rows, err := u.Cluster.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("UserRepo - GetByFamilyID - r.Cluster.Query: %v", err)
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var user entity.User
		err = rows.Scan(
			&user.Id,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.Role,
			&user.FamilyId,
		)
		if err != nil {
			return nil, fmt.Errorf("UserRepo - GetByFamilyID - rows.Scan: %v", err)
		}
		users = append(users, user)
	}
	return users, nil
}

func (u *UserRepo) Update(ctx context.Context, user entity.User) error {
	sql, args, _ := u.Builder.Update(userTable).
		Set("name", user.Name).
		Set("email", user.Email).
		Set("password", user.Password).
		Set("role", user.Role).
		Where("id = ?", user.Id).
		ToSql()

	_, err := u.Cluster.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("UserRepo - Update - r.Cluster.Exec: %v", err)
	}
	return nil
}

func (u *UserRepo) Delete(ctx context.Context, id string) error {
	sql, args, _ := u.Builder.Delete(userTable).
		Where("id = ?", id).
		ToSql()

	_, err := u.Cluster.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("UserRepo - Delete - r.Cluster.Exec: %v", err)
	}
	return nil
}
