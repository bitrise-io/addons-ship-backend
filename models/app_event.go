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

// LogAWSPath ...
func (a *AppEvent) LogAWSPath() string {
	return fmt.Sprintf("/logs/%s/%s", a.App.ID, a.ID)
}
