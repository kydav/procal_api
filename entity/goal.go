package entity

import "gorm.io/gorm"

type Goal struct {
	gorm.Model
	UserID      string `json:"userId" gorm:"type:uuid;not_null;unique"`
	ProteinGoal int    `json:"proteinGoal" gorm:"type:int;not_null;"`
	CalorieGoal int    `json:"calorieGoal" gorm:"type:int;not_null;"`
	WeightGoal  int    `json:"weightGoal" gorm:"type:int;not_null;"`
	Objective   string `json:"objective" gorm:"type:varchar(50);not_null;"`
}
