package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/bitrise-io/addons-ship-backend/bitrise"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

const (
	maxCharNumberForAndroidShortDescription = 80
	maxCharNumberForAndroidFullDescription  = 80
	maxCharNumberForIOSFullDescription      = 255
)

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
	Version              string          `json:"version"`
	Platform             string          `json:"platform"`
	BuildNumber          string          `json:"build_number"`
	BuildSlug            string          `json:"build_slug"`
	LastUpdate           time.Time       `json:"last_update"`
	Scheme               string          `json:"scheme"`
	Configuration        string          `json:"configuration"`
	MinimumOS            string          `json:"minimum_os"`
	MinimumSDK           string          `json:"minimum_sdk"`
	CertificateExpiresAt time.Time       `json:"certificate_expires_at"`
	DistributionType     string          `json:"distribution_type"`
	SupportedDeviceTypes pq.StringArray  `json:"supported_device_types" gorm:"type:varchar(255)[]"`
	AppInfoData          json.RawMessage `json:"-" db:"app_info" gorm:"column:app_info;type:json"`
	ProvisioningInfoData json.RawMessage `json:"-" db:"provisioning_info" gorm:"column:provisioning_info;type:json"`
	AppStoreInfoData     json.RawMessage `json:"-" db:"app_store_info" gorm:"column:app_store_info;type:json"`

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

// AppInfo ...
func (a *AppVersion) AppInfo() (bitrise.AppInfo, error) {
	var appInfo bitrise.AppInfo
	err := json.Unmarshal(a.AppInfoData, &appInfo)
	if err != nil {
		return bitrise.AppInfo{}, err
	}
	return appInfo, nil
}

// ProvisioningInfo ...
func (a *AppVersion) ProvisioningInfo() (bitrise.ProvisioningInfo, error) {
	var provisioningInfo bitrise.ProvisioningInfo
	err := json.Unmarshal(a.ProvisioningInfoData, &provisioningInfo)
	if err != nil {
		return bitrise.ProvisioningInfo{}, err
	}
	return provisioningInfo, nil
}
