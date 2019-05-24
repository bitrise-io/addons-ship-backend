package models

import (
	"fmt"
	"strings"

	uuid "github.com/satori/go.uuid"
)

// Screenshot ...
type Screenshot struct {
	Record
	Filename   string `json:"filename"`
	Filesize   int64  `json:"filesize"`
	Uploaded   bool   `json:"uploaded"`
	DeviceType string `json:"device_type"`
	ScreenSize string `json:"screen_size"`

	AppVersionID uuid.UUID  `db:"app_version_id" json:"-"`
	AppVersion   AppVersion `gorm:"foreignkey:AppVersionID" json:"-"`
}

// BeforeCreate ...
func (s *Screenshot) BeforeCreate() error {
	s.ID = uuid.NewV4()
	return nil
}

// AWSPath ...
func (s *Screenshot) AWSPath() string {
	pathElements := []string{
		s.AppVersion.App.AppSlug,
		s.AppVersion.ID.String(),
		fmt.Sprintf("%s (%s)", s.DeviceType, s.ScreenSize),
		s.Filename,
	}
	return strings.Join(pathElements, "/")
}
