package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// Record ...
type Record struct {
	ID        uuid.UUID  `json:"id" yml:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"created_at" gorm:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" sql:"index"`
}
