package services

import (
	"context"
	fat_secret_wrapper "procal/api/wrappers/FatSecretWrapper"
)

type nutritionService struct {
	fatSecretWrapper fat_secret_wrapper.FatSecretWrapper
}

type NutritionService interface {
	FindById(ctx context.Context) (fat_secret_wrapper.Food, error)
}

func NewNutritionService(fatSecretWrapper fat_secret_wrapper.FatSecretWrapper) NutritionService {
	return &nutritionService{
		fatSecretWrapper: fatSecretWrapper,
	}
}

func (service *nutritionService) FindById(ctx context.Context) (fat_secret_wrapper.Food, error) {
	food, _ := service.fatSecretWrapper.GetFoodFromId(1641)
	return food, nil
}
