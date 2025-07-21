package services

import (
	"context"
	"errors"
	fat_secret_wrapper "procal/wrappers/FatSecretWrapper"
	"strconv"
)

type nutritionService struct {
	fatSecretWrapper fat_secret_wrapper.FatSecretWrapper
}

type NutritionService interface {
	FindById(ctx context.Context, id int) (fat_secret_wrapper.FatSecretFood, error)
	SearchByFoodName(ctx context.Context, search string, pageNumber string) (fat_secret_wrapper.FoodsSearch, error)
	FindByBarcode(ctx context.Context, barcode string) (fat_secret_wrapper.FatSecretFood, error)
}

func NewNutritionService(fatSecretWrapper fat_secret_wrapper.FatSecretWrapper) NutritionService {
	return &nutritionService{
		fatSecretWrapper: fatSecretWrapper,
	}
}

func (service *nutritionService) FindById(ctx context.Context, id int) (fat_secret_wrapper.FatSecretFood, error) {
	food, error := service.fatSecretWrapper.GetFoodFromId(ctx, id)
	if error != nil {
		return fat_secret_wrapper.FatSecretFood{}, error
	}
	return food, nil
}

func (service *nutritionService) FindByBarcode(ctx context.Context, barcode string) (fat_secret_wrapper.FatSecretFood, error) {
	id, error := service.fatSecretWrapper.GetFoodIdFromBarcode(ctx, barcode)
	if error != nil {
		return fat_secret_wrapper.FatSecretFood{}, error
	}

	if id.FoodId.Value == "" {
		return fat_secret_wrapper.FatSecretFood{}, errors.New("no foods found")
	}
	foodId, err := strconv.Atoi(id.FoodId.Value)
	if err != nil {
		return fat_secret_wrapper.FatSecretFood{}, err
	}
	food, err := service.fatSecretWrapper.GetFoodFromId(ctx, foodId)
	if err != nil {
		return fat_secret_wrapper.FatSecretFood{}, err
	}

	return food, nil
}

func (service *nutritionService) SearchByFoodName(ctx context.Context, search string, pageNumber string) (fat_secret_wrapper.FoodsSearch, error) {
	food, error := service.fatSecretWrapper.SearchFoodsByName(ctx, search, &pageNumber)
	if error != nil {
		return fat_secret_wrapper.FoodsSearch{}, error
	}
	return food.FoodsSearch, nil
}
