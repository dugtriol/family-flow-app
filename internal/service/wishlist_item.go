package service

import (
	"context"
	"log/slog"
	`time`

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
	ID          string
	Name        string
	Description string
	Link        string
	Status      string
	IsArchived  bool
	CreatedBy   string
	UpdatedAt   time.Time
}

func (w *WishlistService) Update(ctx context.Context, log *slog.Logger, input WishlistUpdateInput) error {
	log.Info("Service - WishlistService - Update")

	err := w.wishlistRepo.Update(
		ctx, log, entity.WishlistItem{
			ID:          input.ID,
			Name:        input.Name,
			Description: input.Description,
			Link:        input.Link,
			Status:      input.Status,
			IsArchived:  input.IsArchived,
			CreatedBy:   input.CreatedBy,
			UpdatedAt:   time.Now().Add(time.Hour * 3),
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

type WishlistUpdateReservedByInput struct {
	ID         string
	ReservedBy string
}

func (w *WishlistService) UpdateReservedBy(
	ctx context.Context, log *slog.Logger, input WishlistUpdateReservedByInput,
) error {
	log.Info("Service - WishlistService - UpdateReservedBy")
	err := w.wishlistRepo.UpdateReservedBy(ctx, log, input.ID, input.ReservedBy)
	if err != nil {
		log.Error("Service - WishlistService - UpdateReservedBy: %v", err)
		return err
	}

	return nil
}

type WishlistCancelUpdateReservedByInput struct {
	ID string
}

func (w *WishlistService) CancelUpdateReservedBy(
	ctx context.Context, log *slog.Logger, input WishlistCancelUpdateReservedByInput,
) error {
	log.Info("Service - WishlistService - UpdateReservedBy")
	err := w.wishlistRepo.CancelUpdateReservedBy(ctx, log, input.ID)
	if err != nil {
		log.Error("Service - WishlistService - UpdateReservedBy: %v", err)
		return err
	}

	return nil
}

// GetArchivedByUserID
func (w *WishlistService) GetArchivedByUserID(
	ctx context.Context, log *slog.Logger, userID string,
) ([]entity.WishlistItem, error) {
	log.Info("Service - WishlistService - GetArchivedByUserID")

	items, err := w.wishlistRepo.GetArchivedByUserID(ctx, log, userID)
	if err != nil {
		log.Error("Service - WishlistService - GetArchivedByUserID: %v", err)
		return nil, err
	}

	return items, nil
}
