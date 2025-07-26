package repository

import (
	"context"
	"procal/entity"

	"gorm.io/gorm"
)

type GoalRepository interface {
	Create(ctx context.Context, goal *entity.Goal) error
	GetByID(ctx context.Context, id string) (*entity.Goal, error)
	Update(ctx context.Context, goal *entity.Goal) error
	Delete(ctx context.Context, id string) error
}

type goalRepository struct {
	db *gorm.DB
}

func NewGoalRepository(db *gorm.DB) GoalRepository {
	err := db.AutoMigrate(&entity.Goal{})
	if err != nil {
		panic("failed to auto migrate Goal: " + err.Error())
	}
	return &goalRepository{
		db: db,
	}
}

func (r *goalRepository) Create(ctx context.Context, goal *entity.Goal) error {
	if result := r.db.WithContext(ctx).Create(goal); result.Error != nil {
		return result.Error
	}
	return nil
}
func (r *goalRepository) GetByID(ctx context.Context, id string) (*entity.Goal, error) {
	var goal entity.Goal
	if result := r.db.WithContext(ctx).Where("user_id = ?", id).First(&goal); result.Error != nil {
		return nil, result.Error
	}
	return &goal, nil
}
func (r *goalRepository) Update(ctx context.Context, goal *entity.Goal) error {
	if result := r.db.WithContext(ctx).Save(goal); result.Error != nil {
		return result.Error
	}
	return nil
}
func (r *goalRepository) Delete(ctx context.Context, id string) error {
	if result := r.db.WithContext(ctx).Where("user_id = ?", id).Delete(&entity.Goal{}); result.Error != nil {
		return result.Error
	}
	return nil
}
