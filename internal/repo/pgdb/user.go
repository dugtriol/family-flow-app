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
