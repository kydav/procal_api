package repository

import (
	"context"
	"time"

	"procal/entity"

	"gorm.io/gorm"
)

type JournalRepository interface {
	Create(ctx context.Context, entry *entity.JournalEntry) error
	GetByUserAndDate(ctx context.Context, userID string, date time.Time) ([]*entity.JournalEntry, error)
	Update(ctx context.Context, entry *entity.JournalEntry) error
	Delete(ctx context.Context, id string) error
}

type journalRepository struct {
	db *gorm.DB
}

func NewJournalRepository(db *gorm.DB) JournalRepository {
	err := db.AutoMigrate(&entity.JournalEntry{})
	if err != nil {
		panic("failed to auto migrate JournalEntry: " + err.Error())
	}
	return &journalRepository{
		db: db,
	}
}

func (r *journalRepository) Create(ctx context.Context, entry *entity.JournalEntry) error {
	return r.db.WithContext(ctx).Create(entry).Error
}

func (r *journalRepository) GetByUserAndDate(ctx context.Context, userID string, date time.Time) ([]*entity.JournalEntry, error) {
	var entries []*entity.JournalEntry
	err := r.db.WithContext(ctx).Where("user_id = ? AND date = ?", userID, date).Find(&entries).Error
	if err != nil {
		return nil, err
	}
	return entries, nil
}

func (r *journalRepository) Update(ctx context.Context, entry *entity.JournalEntry) error {
	return r.db.WithContext(ctx).Save(entry).Error
}

func (r *journalRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&entity.JournalEntry{}, id).Error
}
