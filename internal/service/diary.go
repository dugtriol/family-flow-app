package service

import (
	"context"
	"log/slog"
	"time"

	"family-flow-app/internal/entity"
	"family-flow-app/internal/repo"
)

type DiaryService struct {
	diaryRepo repo.Diary
}

func NewDiaryService(diaryRepo repo.Diary) *DiaryService {
	return &DiaryService{diaryRepo: diaryRepo}
}

type DiaryCreateInput struct {
	Title       string
	Description string
	Emoji       string
	CreatedBy   string
}

func (s *DiaryService) Create(ctx context.Context, log *slog.Logger, input DiaryCreateInput) (string, error) {
	log.Info("Service - DiaryService - Create")

	item := entity.DiaryItem{
		Title:       input.Title,
		Description: input.Description,
		Emoji:       input.Emoji,
		CreatedBy:   input.CreatedBy,
	}

	id, err := s.diaryRepo.Create(ctx, log, item)
	if err != nil {
		log.Error("Service - DiaryService - Create: %v", err)
		return "", err
	}

	return id, nil
}

func (s *DiaryService) GetByUserID(ctx context.Context, log *slog.Logger, userID string) ([]entity.DiaryItem, error) {
	log.Info("Service - DiaryService - GetByUserID")

	items, err := s.diaryRepo.GetByUserID(ctx, log, userID)
	if err != nil {
		log.Error("Service - DiaryService - GetByUserID: %v", err)
		return nil, err
	}

	return items, nil
}

type DiaryUpdateInput struct {
	ID          string
	Title       string
	Description string
	Emoji       string
}

func (s *DiaryService) Update(ctx context.Context, log *slog.Logger, input DiaryUpdateInput) error {
	log.Info("Service - DiaryService - Update")

	item := entity.DiaryItem{
		ID:          input.ID,
		Title:       input.Title,
		Description: input.Description,
		Emoji:       input.Emoji,
		UpdatedAt:   time.Now(),
	}

	err := s.diaryRepo.Update(ctx, log, item)
	if err != nil {
		log.Error("Service - DiaryService - Update: %v", err)
		return err
	}

	return nil
}

func (s *DiaryService) Delete(ctx context.Context, log *slog.Logger, id string) error {
	log.Info("Service - DiaryService - Delete")

	err := s.diaryRepo.Delete(ctx, log, id)
	if err != nil {
		log.Error("Service - DiaryService - Delete: %v", err)
		return err
	}

	return nil
}
