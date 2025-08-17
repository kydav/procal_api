package services

import (
	"context"
	"time"

	"procal/entity"
	"procal/repository"
)

type JournalService interface {
	CreateEntry(ctx context.Context, entry *entity.JournalEntry) (entity.JournalEntry, error)
	GetEntryByUserAndDate(ctx context.Context, userID string, date time.Time) ([]*entity.JournalEntry, error)
	UpdateEntry(ctx context.Context, entry *entity.JournalEntry) error
	DeleteEntry(ctx context.Context, id string) error
}

type journalService struct {
	repo repository.JournalRepository
}

func NewJournalService(repo repository.JournalRepository) JournalService {
	return &journalService{
		repo: repo,
	}
}

func (s *journalService) CreateEntry(ctx context.Context, entry *entity.JournalEntry) (entity.JournalEntry, error) {
	err := s.repo.Create(ctx, entry)
	return *entry, err
}

func (s *journalService) GetEntryByUserAndDate(ctx context.Context, userID string, date time.Time) ([]*entity.JournalEntry, error) {
	return s.repo.GetByUserAndDate(ctx, userID, date)
}

func (s *journalService) UpdateEntry(ctx context.Context, entry *entity.JournalEntry) error {
	return s.repo.Update(ctx, entry)
}

func (s *journalService) DeleteEntry(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
