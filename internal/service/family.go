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

func (f *FamilyService) Create(ctx context.Context, log *slog.Logger, input FamilyCreateInput) (string, error) {
	log.Info("Service - FamilyService - Create")

	family := entity.Family{
		Name: input.Name,
	}

	id, err := f.familyRepo.Create(ctx, family)
	if err != nil {
		log.Error(fmt.Sprintf("Service - FamilyService - Create: %v", err))
		return "", ErrCannotCreateFamily
	}

	if err = f.userRepo.UpdateFamilyID(ctx, input.CreatorUserId, id); err != nil {
		log.Error(fmt.Sprintf("Service - FamilyService - UpdateFamilyID: %v", err))
		return "", ErrCannotCreateFamily
	}

	log.Info(fmt.Sprintf("Service - FamilyService - familyRepo.Create - id: %s", id))
	return id, nil
}

func (f *FamilyService) GetFamilyByUserID(ctx context.Context, log *slog.Logger, id string) (entity.Family, error) {
	log.Info("Service - FamilyService - GetFamilyByUserID")

	family, err := f.familyRepo.GetByID(ctx, id)
	if err != nil {
		log.Error(fmt.Sprintf("Service - FamilyService - GetFamilyByUserID: %v", err))
		return entity.Family{}, ErrFamilyNotFound
	}

	return family, nil
}

func (f *FamilyService) AddMember(ctx context.Context, log *slog.Logger, input AddMemberToFamilyInput) error {
	log.Info("Service - FamilyService - AddMember")
	var user entity.User
	var err error

	if user, err = f.isExistUser(ctx, log, input.UserEmail); err != nil {
		return ErrUserNotFound
	}

	err = f.userRepo.UpdateRole(ctx, input.UserEmail, input.Role)
	if err != nil {
		log.Error(fmt.Sprintf("Service - UserService - AddMember: %v", err))
		return ErrCannotUpdateUser
	}

	if err = f.userRepo.UpdateFamilyID(ctx, user.Id, input.FamilyId); err != nil {
		log.Error(fmt.Sprintf("Service - FamilyService - AddMember: %v", err))
		return ErrCannotAddMemberToFamily
	}

	return nil
}

func (f *FamilyService) isExistUser(ctx context.Context, log *slog.Logger, email string) (entity.User, error) {
	log.Info(fmt.Sprintf("Service - FamilyService - isExistUser"))

	user, err := f.userRepo.GetByEmail(ctx, email)
	if err != nil {
		log.Error(fmt.Sprintf("Service - FamilyService - isExistUser - GetByID: %v", err))
		return entity.User{}, ErrUserNotFound
	}

	return user, err
}

func (f *FamilyService) GetByFamilyID(ctx context.Context, log *slog.Logger, familyId string) ([]entity.User, error) {
	log.Info(fmt.Sprintf("Service - FamilyService - GetByFamilyID"))

	users, err := f.userRepo.GetByFamilyID(ctx, familyId)
	if err != nil {
		log.Error(fmt.Sprintf("Service - FamilyService - GetByFamilyID: %v", err))
		return nil, ErrCannotGetFamilyMembers
	}

	return users, nil
}

func (f *FamilyService) GetByID(ctx context.Context, log *slog.Logger, id string) (entity.Family, error) {
	log.Info(fmt.Sprintf("Service - FamilyService - GetByID"))

	family, err := f.familyRepo.GetByID(ctx, id)
	if err != nil {
		log.Error(fmt.Sprintf("Service - FamilyService - GetByID: %v", err))
		return entity.Family{}, ErrFamilyNotFound
	}

	return family, nil
}

func (s *FamilyService) UpdatePhoto(ctx context.Context, log *slog.Logger, familyId, photoURL string) error {
	log.Info("Updating family photo", "familyId", familyId, "photoURL", photoURL)

	err := s.familyRepo.UpdatePhoto(ctx, familyId, photoURL)
	if err != nil {
		log.Error("Failed to update family photo", "error", err)
		return err
	}

	log.Info("Family photo updated successfully", "familyId", familyId)
	return nil
}
