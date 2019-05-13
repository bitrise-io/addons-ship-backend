package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// Model ...
type Model struct {
	ID        uuid.UUID  `gorm:"primary_key"`
	CreatedAt time.Time  `gorm:"created_at" db:"created_at"`
	UpdatedAt time.Time  `gorm:"updated_at"`
	DeletedAt *time.Time `sql:"index"`
}
