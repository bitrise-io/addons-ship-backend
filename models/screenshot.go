package models

import uuid "github.com/satori/go.uuid"

// Screenshot ...
type Screenshot struct {
	Record
	FileName string `json:"filename"`
	FileSize string `json:"filesize"`
	Uploaded bool   `json:"uploaded"`

	AppVersionID uuid.UUID  `db:"app_version_id" json:"-"`
	AppVersion   AppVersion `gorm:"foreignkey:AppVersionID" json:"-"`
}

// BeforeCreate ...
func (s *Screenshot) BeforeCreate() error {
	s.ID = uuid.NewV4()
	return nil
}
