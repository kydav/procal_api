package entity

import "gorm.io/gorm"

type Goal struct {
	gorm.Model
	UserID      string `json:"UserId" gorm:"type:uuid;not_null;unique"`
	ProteinGoal int    `json:"ProteinGoal" gorm:"type:int;not_null;"`
	CalorieGoal int    `json:"CalorieGoal" gorm:"type:int;not_null;"`
	WeightGoal  int    `json:"WeightGoal" gorm:"type:int;not_null;"`
	Objective   string `json:"Objective" gorm:"type:varchar(50);not_null;"`
}
