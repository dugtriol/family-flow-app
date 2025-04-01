package service

import (
	"context"
	"log/slog"

	"family-flow-app/internal/entity"
	"family-flow-app/internal/repo"
)

type ShoppingService struct {
	shoppingRepo repo.ShoppingItem
}

func NewShoppingService(shoppingRepo repo.ShoppingItem) *ShoppingService {
	return &ShoppingService{shoppingRepo: shoppingRepo}
}

type ShoppingCreateInput struct {
	FamilyID    string
	Title       string
	Description string
	Visibility  string
	CreatedBy   string
}

func (s *ShoppingService) Create(ctx context.Context, log *slog.Logger, input ShoppingCreateInput) (string, error) {
	log.Info("Service - ShoppingService - Create")

	item := entity.ShoppingItem{
		FamilyID:    input.FamilyID,
		Title:       input.Title,
		Description: input.Description,
		Visibility:  input.Visibility,
		CreatedBy:   input.CreatedBy,
	}

	id, err := s.shoppingRepo.Create(ctx, log, item)
	if err != nil {
		log.Error("Service - ShoppingService - Create: %v", err)
		return "", err
	}

	return id, nil
}

func (s *ShoppingService) Delete(ctx context.Context, log *slog.Logger, id string) error {
	log.Info("Service - ShoppingService - Delete")

	err := s.shoppingRepo.Delete(ctx, log, id)
	if err != nil {
		log.Error("Service - ShoppingService - Delete: %v", err)
		return err
	}

	return nil
}

type ShoppingUpdateInput struct {
	Title       string
	Description string
	Status      string
	Visibility  string
}

func (s *ShoppingService) Update(ctx context.Context, log *slog.Logger, input ShoppingUpdateInput) error {
	log.Info("Service - ShoppingService - Update")

	err := s.shoppingRepo.Update(
		ctx, log, entity.ShoppingItem{
			Title:       input.Title,
			Description: input.Description,
			Status:      input.Status,
			Visibility:  input.Visibility,
		},
	)
	if err != nil {
		log.Error("Service - ShoppingService - Update: %v", err)
		return err
	}

	return nil
}

func (s *ShoppingService) GetPublicByFamilyID(
	ctx context.Context, log *slog.Logger, familyID string,
) ([]entity.ShoppingItem, error) {
	log.Info("Service - ShoppingService - GetPublicByFamilyID")

	items, err := s.shoppingRepo.GetPublicByFamilyID(ctx, log, familyID)
	if err != nil {
		log.Error("Service - ShoppingService - GetPublicByFamilyID: %v", err)
		return nil, err
	}

	return items, nil
}

func (s *ShoppingService) GetPrivateByCreatedBy(
	ctx context.Context, log *slog.Logger, createdBy string,
) ([]entity.ShoppingItem, error) {
	log.Info("Service - ShoppingService - GetPrivateByCreatedBy")

	items, err := s.shoppingRepo.GetPrivateByCreatedBy(ctx, log, createdBy)
	if err != nil {
		log.Error("Service - ShoppingService - GetPrivateByCreatedBy: %v", err)
		return nil, err
	}

	return items, nil
}
