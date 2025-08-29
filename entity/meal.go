package entity

import (
	"time"

	"gorm.io/gorm"
)

type Meal struct {
	gorm.Model
	UserID   string      `json:"UserId" gorm:"type:varchar(100);not_null;"`
	Date     time.Time   `json:"Date" gorm:"type:date;not_null;"`
	MealType MealType    `json:"MealType" gorm:"type:varchar(100);not_null;"`
	Foods    *[]MealFood `json:"Foods" gorm:"->"`
}

type MealType string

const (
	Breakfast MealType = "Breakfast"
	Lunch     MealType = "Lunch"
	Dinner    MealType = "Dinner"
	Snack     MealType = "Snack"
)

type MealFood struct {
	gorm.Model
	MealId     string  `json:"MealId" gorm:"type:string"`
	FoodId     string  `json:"FoodId" gorm:"type:string"`
	FoodName   string  `json:"FoodName" gorm:"type:string"`
	FoodAmount string  `json:"FoodAmount" gorm:"type:varchar(100);not_null;"`
	Protein    float32 `json:"Protein" gorm:"type:float;not_null;"`
	Calories   float32 `json:"Calories" gorm:"type:float;not_null;"`
	Fat        float32 `json:"Fat" gorm:"type:float;not_null;"`
}
