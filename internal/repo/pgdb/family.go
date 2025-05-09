package pgdb

import (
	"context"
	"fmt"

	"family-flow-app/internal/entity"
	"family-flow-app/pkg/postgres"
)

const (
	familyTable = "families"
)

type FamilyRepo struct {
	*postgres.Database
}

func NewFamilyRepo(db *postgres.Database) *FamilyRepo {
	return &FamilyRepo{db}
}

func (r *FamilyRepo) Create(ctx context.Context, family entity.Family) (string, error) {
	sql, args, _ := r.Builder.Insert(familyTable).Columns(
		"name",
	).Values(
		family.Name,
	).Suffix("RETURNING id").ToSql()

	var id string
	err := r.Cluster.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

// GetByID returns a family by its id
func (r *FamilyRepo) GetByID(ctx context.Context, id string) (entity.Family, error) {
	return r.getByField(ctx, "id", id)
}

func (r *FamilyRepo) getByField(ctx context.Context, field, value string) (entity.Family, error) {
	sql, args, _ := r.Builder.Select("*").From(familyTable).Where(
		field+" = ?", value,
	).ToSql()

	var family entity.Family
	err := r.Cluster.QueryRow(ctx, sql, args...).Scan(&family.Id, &family.Name, &family.CreatedAt, &family.Photo)
	if err != nil {
		return entity.Family{}, err
	}
	return family, nil
}

func (r *FamilyRepo) UpdatePhoto(ctx context.Context, familyId, photoURL string) error {
	sql, args, _ := r.Builder.Update(familyTable).
		Set("photo", photoURL).
		Where("id = ?", familyId).
		ToSql()

	_, err := r.Cluster.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to update family photo: %w", err)
	}

	return nil
}
