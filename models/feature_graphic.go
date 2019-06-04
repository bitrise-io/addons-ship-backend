package models

import (
	"github.com/bitrise-io/api-utils/constants"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

const (
	// MaxFeatureGraphicFileByteSize ...
	MaxFeatureGraphicFileByteSize = 10 * constants.MegaByte
)

// FeatureGraphic ...
type FeatureGraphic struct {
	Record
	UploadableObject

	AppVersionID uuid.UUID  `db:"app_version_id" json:"-"`
	AppVersion   AppVersion `gorm:"foreignkey:AppVersionID" json:"-"`
}

// BeforeCreate ...
func (f *FeatureGraphic) BeforeCreate(scope *gorm.Scope) error {
	f.ID = uuid.NewV4()
	err := f.validate(scope, "create")
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// BeforeUpdate ...
func (f *FeatureGraphic) BeforeUpdate(scope *gorm.Scope) error {
	err := f.validate(scope, "update")
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (f *FeatureGraphic) validate(scope *gorm.Scope, action string) error {
	var err error
	if f.Filesize > MaxFeatureGraphicFileByteSize {
		err = scope.DB().AddError(NewValidationError("filesize: Must be smaller than 10 megabytes"))
	}
	if action == "create" {
		var featureGraphicCnt int64
		err = scope.DB().Where(&FeatureGraphic{AppVersionID: f.AppVersionID}).Count(&featureGraphicCnt).Error
		if featureGraphicCnt > 0 {
			err = scope.DB().AddError(NewValidationError("feature_graphics: Maximum count of feature graphics is 1"))
		}
	}
	if err != nil {
		return errors.New("Validation failed")
	}
	return nil
}
