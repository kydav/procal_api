package entity

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID                    string    `json:"ID" gorm:"primary_key;type:uuid;default:uuid_generate_v4();not_null;unique"`
	FirebaseUid           string    `json:"FirebaseUid" gorm:"type:varchar(255);not_null;unique"`
	Email                 string    `json:"Email" gorm:"type:varchar(70);index:user_email_ln_uniq,unique,where:deleted_at is null"`
	FirstName             string    `json:"FirstName" gorm:"type:varchar(70)"`
	LastName              string    `json:"LastName" gorm:"type:varchar(70)"`
	BirthDate             time.Time `json:"BirthDate" gorm:"type:date"`
	CurrentWeight         float32   `json:"CurrentWeight" gorm:"type:float;not_null;"`
	Height                float32   `json:"Height" gorm:"type:float;not_null;"`
	Age                   int       `json:"Age" gorm:"type:int;not_null;"`
	MeasurementPreference string    `json:"MeasurementPreference" gorm:"type:varchar(10);not_null;"`
	IsActive              bool      `json:"IsActive"`
	Gender                string    `json:"Gender" gorm:"type:varchar(10);not_null;"`
}
