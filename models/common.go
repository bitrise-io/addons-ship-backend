package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// Record ...
type Record struct {
	ID        uuid.UUID `json:"id" yml:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at" gorm:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"updated_at"`
}

func removeFromArray(arr []string, elements ...string) []string {
	diff := []string{}

	for _, item := range arr {
		if !contains(elements, item) {
			diff = append(diff, item)
		}
	}

	return diff
}

func contains(arr []string, element string) bool {
	for _, el := range arr {
		if el == element {
			return true
		}
	}
	return false
}
