package service

import (
	"context"
	"fmt"
	"log/slog"

	"family-flow-app/internal/entity"
	"family-flow-app/internal/repo"
)

type FamilyService struct {
	familyRepo repo.Family
	userRepo   repo.User
}

func NewFamilyService(familyRepo repo.Family, userRepo repo.User) *FamilyService {
	return &FamilyService{familyRepo: familyRepo, userRepo: userRepo}
}

func (f *FamilyService) CreateFamily(ctx context.Context, log *slog.Logger, input FamilyCreateInput) (string, error) {
	log.Info("Service - FamilyService - CreateFamily")

	family := entity.Family{
		Name: input.Name,
	}

	id, err := f.familyRepo.Create(ctx, family)
	if err != nil {
		log.Error(fmt.Sprintf("Service - FamilyService - CreateFamily: %v", err))
		return "", ErrCannotCreateFamily
	}

	if err = f.userRepo.UpdateFamilyID(ctx, input.CreatorUserId, id); err != nil {
		log.Error(fmt.Sprintf("Service - FamilyService - UpdateFamilyID: %v", err))
		return "", ErrCannotCreateFamily
	}

	log.Info(fmt.Sprintf("Service - FamilyService - familyRepo.Create - id: %s", id))
	return id, nil
}

func (f *FamilyService) GetFamilyByID(ctx context.Context, log *slog.Logger, id string) (entity.Family, error) {
	log.Info("Service - FamilyService - GetFamilyByID")

	family, err := f.familyRepo.GetByID(ctx, id)
	if err != nil {
		log.Error(fmt.Sprintf("Service - FamilyService - GetFamilyByID: %v", err))
		return entity.Family{}, ErrFamilyNotFound
	}

	return family, nil
}

func (f *FamilyService) AddMemberToFamily(ctx context.Context, log *slog.Logger, input AddMemberToFamilyInput) error {
	log.Info("Service - FamilyService - AddMemberToFamily")

	if _, err := f.isExistUser(ctx, log, input.UserId); err != nil {
		return ErrUserNotFound
	}

	if err := f.userRepo.UpdateFamilyID(ctx, input.UserId, input.FamilyId); err != nil {
		log.Error(fmt.Sprintf("Service - FamilyService - AddMemberToFamily: %v", err))
		return ErrCannotAddMemberToFamily
	}

	return nil
}

func (f *FamilyService) isExistUser(ctx context.Context, log *slog.Logger, userId string) (entity.User, error) {
	log.Info(fmt.Sprintf("Service - FamilyService - isExistUser"))

	user, err := f.userRepo.GetByID(ctx, userId)
	if err != nil {
		log.Error(fmt.Sprintf("Service - FamilyService - isExistUser - GetByID: %v", err))
		return entity.User{}, ErrUserNotFound
	}

	return user, err
}
