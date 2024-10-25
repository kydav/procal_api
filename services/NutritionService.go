package services

import (
	"context"
	fat_secret_wrapper "procal/api/wrappers/FatSecretWrapper"
)

type nutritionService struct {
	fatSecretWrapper fat_secret_wrapper.FatSecretWrapper
}

type NutritionService interface {
	FindById(ctx context.Context) (fat_secret_wrapper.FatSecretFood, error)
}

func NewNutritionService(fatSecretWrapper fat_secret_wrapper.FatSecretWrapper) NutritionService {
	return &nutritionService{
		fatSecretWrapper: fatSecretWrapper,
	}
}

func (service *nutritionService) FindById(ctx context.Context) (fat_secret_wrapper.FatSecretFood, error) {
	food, error := service.fatSecretWrapper.GetFoodFromId(33691)
	if error != nil {
		return fat_secret_wrapper.FatSecretFood{}, error
	}
	return food, nil
}
