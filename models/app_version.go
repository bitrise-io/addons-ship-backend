package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// AppVersion ...
type AppVersion struct {
	Record
	Version     string    `json:"version"`
	Platform    string    `json:"platform"`
	BuildNumber string    `json:"build_number"`
	BuildSlug   string    `json:"build_slug"`
	LastUpdate  time.Time `json:"last_update"`
	Description string    `json:"description"`

	AppID uuid.UUID `db:"app_id" json:"-"`
	App   App       `gorm:"foreignkey:AppID" json:"-"`
}

// BeforeCreate ...
func (a *AppVersion) BeforeCreate() error {
	a.ID = uuid.NewV4()
	return nil
}
