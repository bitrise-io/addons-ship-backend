package models

import (
	"encoding/json"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

const (
	maxCharNumberForEmail = 254
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
	NotificationPreferencesData json.RawMessage `gorm:"column:notification_preferences;type:json" json:"notification_preferences"`
	ConfirmedAt                 time.Time       `json:"confirmed_at"`
	ConfirmationToken           *string         `db:"confirmation_token" json:"-"`

	AppID uuid.UUID `db:"app_id" json:"-"`
	App   *App      `gorm:"foreignkey:AppID" json:"-"`
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

// BeforeSave ...
func (a *AppContact) BeforeSave(scope *gorm.Scope) error {
	err := a.validate(scope)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (a *AppContact) validate(scope *gorm.Scope) error {
	var err error
	if len(a.Email) > maxCharNumberForEmail {
		err = scope.DB().AddError(NewValidationError("email: Too long"))
	}
	ev := EmailVerifier{Email: a.Email}
	if !ev.Verify() {
		err = scope.DB().AddError(NewValidationError("email: Wrong format"))
	}
	if err != nil {
		return errors.New("Validation failed")
	}
	return nil
}

// NotificationPreferences ...
func (a *AppContact) NotificationPreferences() (NotificationPreferences, error) {
	var notificationPreferences NotificationPreferences
	err := json.Unmarshal(a.NotificationPreferencesData, &notificationPreferences)
	if err != nil {
		return NotificationPreferences{}, err
	}
	return notificationPreferences, nil
}
