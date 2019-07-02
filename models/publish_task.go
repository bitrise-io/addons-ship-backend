package models

import uuid "github.com/satori/go.uuid"

// PublishTask ...
type PublishTask struct {
	Record
	TaskID    string `json:"task_id"`
	Completed bool   `json:"completed"`

	AppVersionID uuid.UUID  `db:"app_version_id" json:"-"`
	AppVersion   AppVersion `gorm:"foreignkey:AppVersionID" json:"-"`
}
