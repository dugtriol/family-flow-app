package service

import (
	"context"
	"log/slog"

	"family-flow-app/internal/entity"
	"family-flow-app/internal/repo"
)

type WishlistService struct {
	wishlistRepo repo.WishlistItem
}

func NewWishlistService(wishlistRepo repo.WishlistItem) *WishlistService {
	return &WishlistService{wishlistRepo: wishlistRepo}
}

type WishlistCreateInput struct {
	Name        string
	Description string
	Link        string
	CreatedBy   string
}

func (w *WishlistService) Create(ctx context.Context, log *slog.Logger, input WishlistCreateInput) (string, error) {
	log.Info("Service - WishlistService - Create")

	item := entity.WishlistItem{
		Name:        input.Name,
		Description: input.Description,
		Link:        input.Link,
		IsReserved:  false,
		CreatedBy:   input.CreatedBy,
	}

	id, err := w.wishlistRepo.Create(ctx, log, item)
	if err != nil {
		log.Error("Service - WishlistService - Create: %v", err)
		return "", err
	}

	return id, nil
}

func (w *WishlistService) Delete(ctx context.Context, log *slog.Logger, id string) error {
	log.Info("Service - WishlistService - Delete")

	err := w.wishlistRepo.Delete(ctx, log, id)
	if err != nil {
		log.Error("Service - WishlistService - Delete: %v", err)
		return err
	}

	return nil
}

type WishlistUpdateInput struct {
	Name        string
	Description string
	Link        string
	Status      string
	IsReserved  bool
	CreatedBy   string
}

func (w *WishlistService) Update(ctx context.Context, log *slog.Logger, input WishlistUpdateInput) error {
	log.Info("Service - WishlistService - Update")

	err := w.wishlistRepo.Update(
		ctx, log, entity.WishlistItem{
			Name:        input.Name,
			Description: input.Description,
			Link:        input.Link,
			Status:      input.Status,
			IsReserved:  input.IsReserved,
			CreatedBy:   input.CreatedBy,
		},
	)
	if err != nil {
		log.Error("Service - WishlistService - Update: %v", err)
		return err
	}

	return nil
}

func (w *WishlistService) GetByID(ctx context.Context, log *slog.Logger, id string) ([]entity.WishlistItem, error) {
	log.Info("Service - WishlistService - GetByID")

	items, err := w.wishlistRepo.GetByUserID(ctx, log, id)
	if err != nil {
		log.Error("Service - WishlistService - GetByID: %v", err)
		return nil, err
	}

	return items, nil
}
