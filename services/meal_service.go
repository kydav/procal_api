package services

import (
	"context"
	"time"

	"procal/entity"
	"procal/repository"
)

type MealService interface {
	CreateEntry(ctx context.Context, entry *entity.Meal) (entity.Meal, error)
	GetEntryByUserAndDate(ctx context.Context, userID string, date time.Time) ([]*entity.Meal, error)
	UpdateEntry(ctx context.Context, entry *entity.Meal) error
	DeleteEntry(ctx context.Context, id string) error
}

type mealService struct {
	mealRepo repository.MealRepository
}

func NewMealService(repo repository.MealRepository) MealService {
	return &mealService{
		mealRepo: repo,
	}
}

func (s *mealService) CreateEntry(ctx context.Context, entry *entity.Meal) (entity.Meal, error) {
	err := s.mealRepo.CreateMeal(ctx, entry)
	return *entry, err
}

func (s *mealService) GetEntryByUserAndDate(ctx context.Context, userID string, date time.Time) ([]*entity.Meal, error) {
	return s.mealRepo.GetByUserAndDate(ctx, userID, date)
}

func (s *mealService) UpdateEntry(ctx context.Context, entry *entity.Meal) error {
	return s.mealRepo.Update(ctx, entry)
}

func (s *mealService) DeleteEntry(ctx context.Context, id string) error {
	return s.mealRepo.Delete(ctx, id)
}
