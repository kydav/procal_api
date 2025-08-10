package entity

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID                    string    `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4();not_null;unique"`
	FirebaseUid           string    `json:"firebase_uid" gorm:"type:varchar(255);not_null;unique"`
	Email                 string    `json:"email" gorm:"type:varchar(70);index:user_email_ln_uniq,unique,where:deleted_at is null"`
	FirstName             string    `json:"first_name" gorm:"type:varchar(70)"`
	LastName              string    `json:"last_name" gorm:"type:varchar(70)"`
	BirthDate             time.Time `json:"birth_date" gorm:"type:date"`
	CurrentWeight         int       `json:"current_weight" gorm:"type:int;not_null;"`
	Height                int       `json:"height" gorm:"type:int;not_null;"`
	Age                   int       `json:"age" gorm:"type:int;not_null;"`
	MeasurementPreference string    `json:"measurement_preference" gorm:"type:varchar(10);not_null;"`
	IsActive              bool      `json:"is_active"`
}
