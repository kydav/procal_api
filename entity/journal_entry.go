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
	Protein    int       `json:"Protein" gorm:"type:int;not_null;"`
	Calories   int       `json:"Calories" gorm:"type:int;not_null;"`
}
