package models

import (
	"fmt"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

// AppVersionEvent ...
type AppVersionEvent struct {
	Record
	Status         string `json:"status"`
	Text           string `json:"event_text" gorm:"column:event_text"`
	IsLogAvailable bool   `json:"is_log_available"`

	AppVersionID uuid.UUID  `db:"app_version_id" json:"-"`
	AppVersion   AppVersion `gorm:"foreignkey:AppVersionID" json:"-"`
}

// BeforeCreate ...
func (a *AppVersionEvent) BeforeCreate() error {
	if uuid.Equal(a.ID, uuid.UUID{}) {
		a.ID = uuid.NewV4()
	}
	return nil
}

// LogAWSPath ...
func (a *AppVersionEvent) LogAWSPath() (string, error) {
	if a.AppVersion.App.AppSlug == "" {
		return "", errors.New("App has empty App Slug, App has to be preloaded")
	}
	return fmt.Sprintf("logs/%s/%s/%s.log", a.AppVersion.App.AppSlug, a.AppVersionID, a.ID), nil
}
