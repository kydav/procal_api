package entity

import (
	"time"

	"gorm.io/gorm"
)

type JournalEntry struct {
	gorm.Model
	UserID     string    `json:"UserId" gorm:"type:varchar(100);not_null;"`
	FoodId     string    `json:"FoodId" gorm:"type:varchar(100);not_null;"`
	FoodAmount string    `json:"FoodAmount" gorm:"type:varchar(100);not_null;"`
	Date       time.Time `json:"Date" gorm:"type:date;not_null;"`
	Protein    float32   `json:"Protein" gorm:"type:float;not_null;"`
	Calories   float32   `json:"Calories" gorm:"type:float;not_null;"`
	Fat        float32   `json:"Fat" gorm:"type:float;not_null;"`
	Meal       string    `json:"Meal" gorm:"type:varchar(100);not_null;"`
}
