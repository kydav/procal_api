package repository

import (
	"context"
	"time"

	"procal/entity"

	"gorm.io/gorm"
)

type MealRepository interface {
	CreateMeal(ctx context.Context, entry entity.Meal) (entity.Meal, error)
	GetByUserAndDate(ctx context.Context, userID string, date time.Time) ([]entity.Meal, error)
	Update(ctx context.Context, entry entity.Meal) error
	Delete(ctx context.Context, id string) error
}

type mealRepository struct {
	db *gorm.DB
}

func NewMealRepository(db *gorm.DB) MealRepository {
	err := db.AutoMigrate(&entity.Meal{})
	if err != nil {
		panic("failed to auto migrate MealEntry: " + err.Error())
	}
	return &mealRepository{
		db: db,
	}
}

func (r *mealRepository) CreateMeal(ctx context.Context, entry entity.Meal) (entity.Meal, error) {
	if err := r.db.WithContext(ctx).Create(&entry).Error; err != nil {
		return entity.Meal{}, err
	}

	return entry, nil
}

func (r *mealRepository) GetByUserAndDate(ctx context.Context, userID string, date time.Time) ([]entity.Meal, error) {
	var entries []entity.Meal
	err := r.db.WithContext(ctx).Where("user_id = ? AND date = ?", userID, date).Find(&entries).Error
	if err != nil {
		return nil, err
	}
	return entries, nil
}

func (r *mealRepository) Update(ctx context.Context, entry entity.Meal) error {
	return r.db.WithContext(ctx).Save(entry).Error
}

func (r *mealRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&entity.Meal{}, id).Error
}
