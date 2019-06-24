package models

import uuid "github.com/satori/go.uuid"

// App ...
type App struct {
	Record
	AppSlug         string `json:"app_slug"`
	Plan            string `json:"plan"`
	BitriseAPIToken string `json:"-"`
	APIToken        string `json:"-"`

	AppVersions []AppVersion `gorm:"foreignkey:AppID" json:"app_versions"`
	AppSettings AppSettings  `gorm:"foreignkey:AppsID" json:"app_settings"`
}

// BeforeCreate ...
func (a *App) BeforeCreate() error {
	if uuid.Equal(a.ID, uuid.UUID{}) {
		a.ID = uuid.NewV4()
	}
	return nil
}
