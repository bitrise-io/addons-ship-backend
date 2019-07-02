package models

import uuid "github.com/satori/go.uuid"

// PublishTask ...
type PublishTask struct {
	Record
	TaskID string `json:"task_id"`

	AppVersionID uuid.UUID  `db:"app_version_id" json:"-"`
	AppVersion   AppVersion `gorm:"foreignkey:AppVersionID" json:"-"`
}

// BeforeCreate ...
func (t *PublishTask) BeforeCreate() error {
	if uuid.Equal(t.ID, uuid.UUID{}) {
		t.ID = uuid.NewV4()
	}
	return nil
}
