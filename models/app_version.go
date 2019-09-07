package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

const (
	maxCharNumberForAndroidShortDescription = 80
	maxCharNumberForAndroidFullDescription  = 80
	maxCharNumberForIOSFullDescription      = 255
)

// ArtifactInfo ...
type ArtifactInfo struct {
	Version              string    `json:"version"`
	VersionCode          string    `json:"version_code"`
	MinimumOS            string    `json:"minimum_os"`
	MinimumSDK           string    `json:"minimum_sdk"`
	Size                 int64     `json:"size"`
	BundleID             string    `json:"bundle_id"`
	SupportedDeviceTypes []string  `json:"supported_device_types"`
	PackageName          string    `json:"package_name"`
	ExpireDate           time.Time `json:"expire_date"`
	IPAExportMethod      string    `json:"ipa_export_method"`
}

// AppStoreInfo ...
type AppStoreInfo struct {
	ShortDescription string `json:"short_description"`
	FullDescription  string `json:"full_description"`
	WhatsNew         string `json:"whats_new"`
	PromotionalText  string `json:"promotional_text"`
	Keywords         string `json:"keywords"`
	ReviewNotes      string `json:"review_notes"`
	SupportURL       string `json:"support_url"`
	MarketingURL     string `json:"marketing_url"`
}

// AppVersion ...
type AppVersion struct {
	Record
	Platform         string          `json:"platform"`
	BuildNumber      string          `json:"build_number"`
	BuildSlug        string          `json:"build_slug"`
	LastUpdate       time.Time       `json:"last_update"`
	Scheme           string          `json:"scheme"`
	Configuration    string          `json:"configuration"`
	CommitMessage    string          `json:"commit_message"`
	ArtifactInfoData json.RawMessage `json:"-" db:"artifact_info" gorm:"column:artifact_info;type:json"`
	AppStoreInfoData json.RawMessage `json:"-" db:"app_store_info" gorm:"column:app_store_info;type:json"`

	AppID uuid.UUID `db:"app_id" json:"-"`
	App   App       `gorm:"foreignkey:AppID" json:"-"`
}

// BeforeCreate ...
func (a *AppVersion) BeforeCreate(scope *gorm.Scope) error {
	if uuid.Equal(a.ID, uuid.UUID{}) {
		a.ID = uuid.NewV4()
	}
	if a.AppStoreInfoData == nil {
		a.AppStoreInfoData = json.RawMessage(`{}`)
	}
	if a.ArtifactInfoData == nil {
		a.ArtifactInfoData = json.RawMessage(`{}`)
	}
	err := a.validate(scope)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// BeforeUpdate ...
func (a *AppVersion) BeforeUpdate(scope *gorm.Scope) error {
	err := a.validate(scope)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (a *AppVersion) validate(scope *gorm.Scope) error {
	appStoreInfo, err := a.AppStoreInfo()
	if err != nil {
		return err
	}
	if a.Platform == "android" {
		if len(appStoreInfo.ShortDescription) > maxCharNumberForAndroidShortDescription {
			err = scope.DB().AddError(NewValidationError(fmt.Sprintf("short_description: Mustn't be longer than %d characters", maxCharNumberForAndroidShortDescription)))
		}
		if len(appStoreInfo.FullDescription) > maxCharNumberForAndroidFullDescription {
			err = scope.DB().AddError(NewValidationError(fmt.Sprintf("full_description: Mustn't be longer than %d characters", maxCharNumberForAndroidFullDescription)))
		}
	}
	if a.Platform == "ios" {
		if len(appStoreInfo.FullDescription) > maxCharNumberForIOSFullDescription {
			err = scope.DB().AddError(NewValidationError(fmt.Sprintf("full_description: Mustn't be longer than %d characters", maxCharNumberForIOSFullDescription)))
		}
	}
	artifactInfo, err := a.ArtifactInfo()
	if err != nil {
		return errors.WithStack(err)
	}
	err = artifactInfo.validate(scope)
	if err != nil {
		return errors.New("Validation failed")
	}
	return nil
}

// AppStoreInfo ...
func (a *AppVersion) AppStoreInfo() (AppStoreInfo, error) {
	var appStoreInfo AppStoreInfo
	err := json.Unmarshal(a.AppStoreInfoData, &appStoreInfo)
	if err != nil {
		return AppStoreInfo{}, err
	}
	return appStoreInfo, nil
}

// ArtifactInfo ...
func (a *AppVersion) ArtifactInfo() (ArtifactInfo, error) {
	var artifactInfo ArtifactInfo
	err := json.Unmarshal(a.ArtifactInfoData, &artifactInfo)
	if err != nil {
		return ArtifactInfo{}, err
	}
	return artifactInfo, nil
}

func (a *ArtifactInfo) validate(scope *gorm.Scope) error {
	if a.Version == "" {
		return scope.DB().AddError(NewValidationError("version: Cannot be empty"))
	}
	return nil
}
