package models

import (
	"fmt"

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

// BeforeCreate ...
func (a *AppEvent) BeforeCreate() error {
	if uuid.Equal(a.ID, uuid.UUID{}) {
		a.ID = uuid.NewV4()
	}
	return nil
}

// LogAWSPath ...
func (a *AppEvent) LogAWSPath() string {
	return fmt.Sprintf("logs/%s/%s.log", a.App.AppSlug, a.ID)
}
