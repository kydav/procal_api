package services

import (
	"context"
	"procal/entity"
	"procal/repository"
)

type MealFoodService interface {
	Create(ctx context.Context, mealFood *entity.MealFood) error
	CreateMealFoods(ctx context.Context, mealFoods []*entity.MealFood) error
	GetFoodsByMealID(ctx context.Context, mealID string) ([]*entity.MealFood, error)
	Update(ctx context.Context, mealFood *entity.MealFood) error
	Delete(ctx context.Context, id string) error
}

type mealFoodService struct {
	mealFoodRepo repository.MealFoodRepository
}

func NewMealFoodService(repo repository.MealFoodRepository) *mealFoodService {
	return &mealFoodService{mealFoodRepo: repo}
}

func (s *mealFoodService) Create(ctx context.Context, mealFood *entity.MealFood) error {
	return s.mealFoodRepo.CreateMealFood(ctx, mealFood)
}

func (s *mealFoodService) CreateMealFoods(ctx context.Context, mealFoods []*entity.MealFood) error {
	return s.mealFoodRepo.CreateMealFoods(ctx, mealFoods)
}

func (s *mealFoodService) GetFoodsByMealID(ctx context.Context, id string) ([]*entity.MealFood, error) {
	return s.mealFoodRepo.GetByMealID(ctx, id)
}

func (s *mealFoodService) Update(ctx context.Context, mealFood *entity.MealFood) error {
	return s.mealFoodRepo.Update(ctx, mealFood)
}

func (s *mealFoodService) Delete(ctx context.Context, id string) error {
	return s.mealFoodRepo.Delete(ctx, id)
}
