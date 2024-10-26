package services

import (
	"context"
	"errors"
	fat_secret_wrapper "procal/api/wrappers/FatSecretWrapper"
	"strconv"
)

type nutritionService struct {
	fatSecretWrapper fat_secret_wrapper.FatSecretWrapper
}

type NutritionService interface {
	FindById(ctx context.Context, id int) (fat_secret_wrapper.FatSecretFood, error)
	SearchByFoodName(ctx context.Context, search string, pageNumber int) ([]fat_secret_wrapper.Food, error)
	FindByBarcode(ctx context.Context, barcode string) (fat_secret_wrapper.FatSecretFood, error)
}

func NewNutritionService(fatSecretWrapper fat_secret_wrapper.FatSecretWrapper) NutritionService {
	return &nutritionService{
		fatSecretWrapper: fatSecretWrapper,
	}
}

func (service *nutritionService) FindById(ctx context.Context, id int) (fat_secret_wrapper.FatSecretFood, error) {
	food, error := service.fatSecretWrapper.GetFoodFromId(id)
	if error != nil {
		return fat_secret_wrapper.FatSecretFood{}, error
	}
	return food, nil
}

func (service *nutritionService) FindByBarcode(ctx context.Context, barcode string) (fat_secret_wrapper.FatSecretFood, error) {
	id, error := service.fatSecretWrapper.GetFoodIdFromBarcode(barcode)
	if error != nil {
		return fat_secret_wrapper.FatSecretFood{}, error
	}

	if id.FoodId.Value == "" {
		return fat_secret_wrapper.FatSecretFood{}, errors.New("No foods found")
	}
	foodId, err := strconv.Atoi(id.FoodId.Value)
	if err != nil {
		return fat_secret_wrapper.FatSecretFood{}, err
	}
	food, error := service.fatSecretWrapper.GetFoodFromId(foodId)

	return food, nil
}

func (service *nutritionService) SearchByFoodName(ctx context.Context, search string, pageNumber int) ([]fat_secret_wrapper.Food, error) {
	food, error := service.fatSecretWrapper.SearchFoodsByName(search, &pageNumber)
	if error != nil {
		return []fat_secret_wrapper.Food{}, error
	}
	return food.FoodsSearch.Results.Food, nil
}
