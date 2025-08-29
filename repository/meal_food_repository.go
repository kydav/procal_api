package repository

import (
	"context"

	"procal/entity"

	"gorm.io/gorm"
)

type MealFoodRepository interface {
	CreateMealFood(ctx context.Context, entry *entity.MealFood) error
	CreateMealFoods(ctx context.Context, foods []*entity.MealFood) error
	GetByMealID(ctx context.Context, mealID string) ([]*entity.MealFood, error)
	Update(ctx context.Context, entry *entity.MealFood) error
	Delete(ctx context.Context, id string) error
}

type mealFoodRepository struct {
	db *gorm.DB
}

func NewMealFoodRepository(db *gorm.DB) MealFoodRepository {
	err := db.AutoMigrate(&entity.MealFood{})
	if err != nil {
		panic("failed to auto migrate MealEntry: " + err.Error())
	}
	return &mealFoodRepository{
		db: db,
	}
}

func (r *mealFoodRepository) CreateMealFood(ctx context.Context, entry *entity.MealFood) error {
	if err := r.db.WithContext(ctx).Create(entry).Error; err != nil {
		return err
	}
	return nil
}

func (r *mealFoodRepository) CreateMealFoods(ctx context.Context, foods []*entity.MealFood) error {
	if err := r.db.WithContext(ctx).Create(foods).Error; err != nil {
		return err
	}
	return nil
}

func (r *mealFoodRepository) GetByMealID(ctx context.Context, mealID string) ([]*entity.MealFood, error) {
	var foods []*entity.MealFood
	err := r.db.WithContext(ctx).Where("meal_id = ?", mealID).Find(&foods).Error
	if err != nil {
		return nil, err
	}
	return foods, nil
}

func (r *mealFoodRepository) Update(ctx context.Context, entry *entity.MealFood) error {
	return r.db.WithContext(ctx).Save(entry).Error
}

func (r *mealFoodRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&entity.MealFood{}, id).Error
}
