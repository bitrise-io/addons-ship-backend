package models

import (
	"encoding/json"
	"time"

	uuid "github.com/satori/go.uuid"
)

// NotificationPreferences ...
type NotificationPreferences struct {
	NewVersion        bool `json:"new_version"`
	SuccessfulPublish bool `json:"successful_publish"`
	FailedPublish     bool `json:"failed_publish"`
}

// AppContact ...
type AppContact struct {
	Record
	Email                       string          `json:"email"`
	NotificationPreferencesData json.RawMessage `json:"notification_preferences"`
	ConfirmedAt                 *time.Time      `json:"confirmed_at"`
	ConfirmationToken           string          `json:"-"`

	AppID uuid.UUID `db:"app_id" json:"-"`
	App   App       `gorm:"foreignkey:AppID" json:"-"`
}

// BeforeCreate ...
func (a *AppContact) BeforeCreate() error {
	if uuid.Equal(a.ID, uuid.UUID{}) {
		a.ID = uuid.NewV4()
	}

	if a.NotificationPreferencesData == nil {
		a.NotificationPreferencesData = json.RawMessage(`{}`)
	}

	return nil
}
