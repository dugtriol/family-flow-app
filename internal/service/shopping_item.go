package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

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
	ID          string
	Title       string
	Description string
	Status      string
	Visibility  string
	IsArchived  bool
}

func (s *ShoppingService) Update(ctx context.Context, log *slog.Logger, input ShoppingUpdateInput) error {
	log.Info("Service - ShoppingService - Update")
	log.Info(fmt.Sprintf("ShoppingService - Update - IsArchived - %s", input.IsArchived))

	err := s.shoppingRepo.Update(
		ctx, log, entity.ShoppingItem{
			ID:          input.ID,
			Title:       input.Title,
			Description: input.Description,
			Status:      input.Status,
			Visibility:  input.Visibility,
			IsArchived:  input.IsArchived,
			UpdatedAt:   time.Now().Add(time.Hour * 3),
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

type ShoppingUpdateReservedByInput struct {
	Id         string
	ReservedBy string
}

func (s *ShoppingService) UpdateReservedBy(
	ctx context.Context, log *slog.Logger, input ShoppingUpdateReservedByInput,
) error {
	log.Info("Service - ShoppingService - UpdateReservedBy")

	err := s.shoppingRepo.UpdateReservedBy(ctx, log, input.Id, input.ReservedBy, time.Now().Add(time.Hour*3))
	if err != nil {
		log.Error("Service - ShoppingService - UpdateReservedBy: %v", err)
		return err
	}

	return nil
}

type ShoppingUpdateBuyerIdInput struct {
	Id      string
	BuyerId string
}

func (s *ShoppingService) UpdateBuyerId(
	ctx context.Context, log *slog.Logger, input ShoppingUpdateBuyerIdInput,
) error {
	log.Info("Service - ShoppingService - UpdateBuyerId")

	err := s.shoppingRepo.UpdateBuyerId(ctx, log, input.Id, input.BuyerId, time.Now().Add(time.Hour*3))
	if err != nil {
		log.Error("Service - ShoppingService - UpdateBuyerId: %v", err)
		return err
	}

	return nil
}

// get archived items by user id
func (s *ShoppingService) GetArchivedByUserID(
	ctx context.Context, log *slog.Logger, userID string,
) ([]entity.ShoppingItem, error) {
	log.Info("Service - ShoppingService - GetArchivedByUserID")

	items, err := s.shoppingRepo.GetArchivedByUserID(ctx, log, userID)
	if err != nil {
		log.Error("Service - ShoppingService - GetArchivedByUserID: %v", err)
		return nil, err
	}

	return items, nil
}

type ShoppingCancelUpdateReservedByInput struct {
	Id string
}

func (s *ShoppingService) CancelUpdateReservedBy(
	ctx context.Context, log *slog.Logger, input ShoppingCancelUpdateReservedByInput,
) error {
	log.Info("Service - ShoppingService - UpdateReservedBy")

	err := s.shoppingRepo.CancelUpdateReservedBy(ctx, log, input.Id, time.Now().Add(time.Hour*3))
	if err != nil {
		log.Error("Service - ShoppingService - UpdateReservedBy: %v", err)
		return err
	}

	return nil
}

// get by id
func (s *ShoppingService) GetByID(ctx context.Context, log *slog.Logger, id string) (entity.ShoppingItem, error) {
	log.Info("Service - ShoppingService - GetByID")

	items, err := s.shoppingRepo.GetByID(ctx, log, id)
	if err != nil {
		log.Error("Service - ShoppingService - GetByID: %v", err)
		return entity.ShoppingItem{}, err
	}

	return items, nil
}
