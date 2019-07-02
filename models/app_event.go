package models

import (
	"fmt"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

// AppEvent ...
type AppEvent struct {
	Record
	Status string `json:"status"`
	Text   string `json:"text"`

	AppID uuid.UUID `db:"app_id" json:"-"`
	App   App       `gorm:"foreignkey:AppID" json:"-"`
}

// LogAWSPath ...
func (a *AppEvent) LogAWSPath() (string, error) {
	if a.App.AppSlug == "" {
		return "", errors.New("App has empty App Slug, App has to be preloaded")
	}
	return fmt.Sprintf("/logs/%s/%s", a.App.AppSlug, a.ID), nil
}
