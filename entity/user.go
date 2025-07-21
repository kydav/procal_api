package entity

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID        string `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4();not_null;unique"`
	Email     string `json:"email" gorm:"type:varchar(70);index:user_email_ln_uniq,unique,where:deleted_at is null"`
	FirstName string `json:"first_name" gorm:"type:varchar(70)"`
	LastName  string `json:"last_name" gorm:"type:varchar(70)"`
	IsActive  bool   `json:"is_active"`
}
