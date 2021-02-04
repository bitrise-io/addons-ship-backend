package models

import (
	"encoding/json"
	"reflect"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/thoas/go-funk"
)

// IosSettings ...
type IosSettings struct {
	AppSKU                               string   `json:"app_sku"`
	AppleDeveloperAccountEmail           string   `json:"apple_developer_account_email"`
	ApplSpecificPassword                 string   `json:"app_specific_password"`
	SelectedAppStoreProvisioningProfiles []string `json:"selected_app_store_provisioning_profiles"`
	SelectedCodeSigningIdentity          string   `json:"selected_code_signing_identity"`
	IncludeBitCode                       bool     `json:"include_bit_code"`
}

// Valid ...
func (s IosSettings) Valid() bool {
	return !reflect.DeepEqual(s, IosSettings{})
}

// ValidateProvisioningProfileSlugs ...
func (s *IosSettings) ValidateSelectedProvisioningProfileSlugs(provProfiles []string) {
	if len(provProfiles) == 0 && len(s.SelectedAppStoreProvisioningProfiles) > 0 {
		s.SelectedAppStoreProvisioningProfiles = []string{}
		return
	}
	for _, slug := range s.SelectedAppStoreProvisioningProfiles {
		valid := false
		for _, provProfileSlug := range provProfiles {
			if provProfileSlug == slug {
				valid = true
				break
			}
		}
		if !valid {
			s.SelectedAppStoreProvisioningProfiles = funk.SubtractString(s.SelectedAppStoreProvisioningProfiles, []string{slug})
		}
	}
}

// AndroidSettings ...
type AndroidSettings struct {
	Track                  string `json:"track"`
	SelectedKeystoreFile   string `json:"selected_keystore_file"`
	SelectedServiceAccount string `json:"selected_service_account"`
	Module                 string `json:"module"`
}

// Valid ...
func (s AndroidSettings) Valid() bool {
	return s != (AndroidSettings{})
}

// AppSettings ...
type AppSettings struct {
	Record
	IosWorkflow         string          `json:"ios_workflow"`
	AndroidWorkflow     string          `json:"android_workflow"`
	IosSettingsData     json.RawMessage `json:"-" db:"ios_settings" gorm:"column:ios_settings;type:json"`
	AndroidSettingsData json.RawMessage `json:"-" db:"android_settings" gorm:"column:android_settings;type:json"`

	AppID uuid.UUID `db:"app_id" json:"-"`
	App   *App      `gorm:"foreignkey:AppID" json:"-"`
}

// BeforeCreate ...
func (a *AppSettings) BeforeCreate(scope *gorm.Scope) error {
	if uuid.Equal(a.ID, uuid.UUID{}) {
		a.ID = uuid.NewV4()
	}
	if a.IosSettingsData == nil {
		a.IosSettingsData = json.RawMessage(`{}`)
	}
	if a.AndroidSettingsData == nil {
		a.AndroidSettingsData = json.RawMessage(`{}`)
	}
	return nil
}

// IosSettings ...
func (a *AppSettings) IosSettings() (IosSettings, error) {
	var iosSettings IosSettings
	err := json.Unmarshal(a.IosSettingsData, &iosSettings)
	if err != nil {
		return IosSettings{}, err
	}
	return iosSettings, nil
}

// AndroidSettings ...
func (a *AppSettings) AndroidSettings() (AndroidSettings, error) {
	var androidSettings AndroidSettings
	err := json.Unmarshal(a.AndroidSettingsData, &androidSettings)
	if err != nil {
		return AndroidSettings{}, err
	}
	return androidSettings, nil
}
