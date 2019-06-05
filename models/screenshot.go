package models

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bitrise-io/api-utils/constants"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

const (
	// MaxScreenshotFileByteSize ...
	MaxScreenshotFileByteSize = 10 * constants.MegaByte
)

// Screenshot ...
type Screenshot struct {
	Record
	UploadableObject
	DeviceType string `json:"device_type"`
	ScreenSize string `json:"screen_size"`

	AppVersionID uuid.UUID  `db:"app_version_id" json:"-"`
	AppVersion   AppVersion `gorm:"foreignkey:AppVersionID" json:"-"`
}

// BeforeCreate ...
func (s *Screenshot) BeforeCreate(scope *gorm.Scope) error {
	s.ID = uuid.NewV4()
	return nil
}

// BeforeSave ...
func (s *Screenshot) BeforeSave(scope *gorm.Scope) error {
	err := s.validate(scope)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (s *Screenshot) validate(scope *gorm.Scope) error {
	var err error
	if s.Filesize > MaxScreenshotFileByteSize {
		err = scope.DB().AddError(NewValidationError("filesize: Must be smaller than 10 megabytes"))
	}
	if err != nil {
		return errors.New("Validation failed")
	}
	return nil
}

// AWSPath ...
func (s *Screenshot) AWSPath() string {
	pathElements := []string{
		s.AppVersion.App.AppSlug,
		s.AppVersion.ID.String(),
		fmt.Sprintf("%s (%s)", s.DeviceType, s.ScreenSize),
		s.ID.String() + filepath.Ext(s.Filename),
	}
	return strings.Join(pathElements, "/")
}
