package services

import (
	"context"
	"procal/entity"
	"procal/repository"
)

type GoalService interface {
	CreateGoal(ctx context.Context, goal *entity.Goal) error
	GetGoalByID(ctx context.Context, id string) (*entity.Goal, error)
	UpdateGoal(ctx context.Context, goal *entity.Goal) error
	DeleteGoal(ctx context.Context, id string) error
}

type goalService struct {
	repo repository.GoalRepository
}

func NewGoalService(repo repository.GoalRepository) GoalService {
	return &goalService{
		repo: repo,
	}
}
func (s *goalService) CreateGoal(ctx context.Context, goal *entity.Goal) error {
	return s.repo.Create(ctx, goal)
}
func (s *goalService) GetGoalByID(ctx context.Context, id string) (*entity.Goal, error) {
	return s.repo.GetByID(ctx, id)
}
func (s *goalService) UpdateGoal(ctx context.Context, goal *entity.Goal) error {
	return s.repo.Update(ctx, goal)
}
func (s *goalService) DeleteGoal(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
