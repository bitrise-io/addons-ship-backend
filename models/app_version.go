package models

import (
	"encoding/json"
	"time"

	uuid "github.com/satori/go.uuid"
)

// AppStoreInfo ...
type AppStoreInfo struct {
	ShortDescription string `json:"short_description"`
	FullDescription  string `json:"full_description"`
	WhatsNew         string `json:"whats_new"`
	PromotionalText  string `json:"promotional_text"`
	Keywords         string `json:"keywords"`
	ReviewNotes      string `json:"review_notes"`
	SupportURL       string `json:"support_url"`
	MarketingURL     string `json:"marketing_url"`
}

// AppVersion ...
type AppVersion struct {
	Record
	Version          string          `json:"version"`
	Platform         string          `json:"platform"`
	BuildNumber      string          `json:"build_number"`
	BuildSlug        string          `json:"build_slug"`
	LastUpdate       time.Time       `json:"last_update"`
	Scheme           string          `json:"scheme"`
	Configuration    string          `json:"configuration"`
	AppStoreInfoData json.RawMessage `json:"-" gorm:"column:app_store_info;type:json"`

	AppID uuid.UUID `db:"app_id" json:"-"`
	App   App       `gorm:"foreignkey:AppID" json:"-"`
}

// BeforeCreate ...
func (a *AppVersion) BeforeCreate() error {
	a.ID = uuid.NewV4()
	if a.AppStoreInfoData == nil {
		a.AppStoreInfoData = json.RawMessage(`{}`)
	}
	return nil
}

// AppStoreInfo ...
func (a *AppVersion) AppStoreInfo() (AppStoreInfo, error) {
	var appStoreInfo AppStoreInfo
	err := json.Unmarshal(a.AppStoreInfoData, &appStoreInfo)
	if err != nil {
		return AppStoreInfo{}, err
	}
	return appStoreInfo, nil
}
