package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// AppVersion ...
type AppVersion struct {
	Record
	Version         string    `json:"version"`
	Platform        string    `json:"platform"`
	BuildNumber     string    `json:"build_number"`
	BuildSlug       string    `json:"build_slug"`
	LastUpdate      time.Time `json:"last_update"`
	Description     string    `json:"description"`
	WhatsNew        string    `json:"whats_new"`
	PromotionalText string    `json:"promotional_text"`
	Keywords        string    `json:"keywords"`
	ReviewNotes     string    `json:"review_notes"`
	SupportURL      string    `json:"support_url"`
	MarketingURL    string    `json:"marketing_url"`
	Scheme          string    `json:"scheme"`
	Configuration   string    `json:"configuration"`

	AppID uuid.UUID `db:"app_id" json:"-"`
	App   App       `gorm:"foreignkey:AppID" json:"-"`
}

// BeforeCreate ...
func (a *AppVersion) BeforeCreate() error {
	a.ID = uuid.NewV4()
	return nil
}
