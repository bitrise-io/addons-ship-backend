// +build mage

package main

import (
	"encoding/json"
	"fmt"

	"github.com/bitrise-io/addons-ship-backend/dataservices"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/api-utils/logging"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type PreviousIosSettings struct {
	AppSKU                              string `json:"app_sku"`
	AppleDeveloperAccountEmail          string `json:"apple_developer_account_email"`
	ApplSpecificPassword                string `json:"app_specific_password"`
	SelectedAppStoreProvisioningProfile string `json:"selected_app_store_provisioning_profile"`
	SelectedCodeSigningIdentity         string `json:"selected_code_signing_identity"`
	IncludeBitCode                      bool   `json:"include_bit_code"`
}

func MigrateSelectedProvisioningProfileSlugToArray() error {
	logger := logging.WithContext(nil)
	defer func() {
		err := logger.Sync()
		if err != nil {
			fmt.Println(err)
		}
	}()
	err := dataservices.InitializeConnection(dataservices.ConnectionParams{}, true)
	if err != nil {
		logger.Error("Failed to initialize database connection", zap.Error(errors.WithStack(err)))
		return errors.WithStack(err)
	}

	appSettingsToMigrate := []models.AppSettings{}
	dataservices.GetDB().Where("ios_settings->>'selected_app_store_provisioning_profile' IS NOT NULL").Find(&appSettingssToMigrate)
	appSettingsService := models.AppSettingsService{DB: dataservices.GetDB()}
	for _, appSettings := range appSettingsToMigrate {
		iosSettings := PreviousIosSettings{}
		err := json.Unmarshal(appSettings.IosSettingsData, &iosSettings)
		if err != nil {
			logger.Error("Failed to unmarshal IosSettings struct", zap.Error(errors.WithStack(err)))
			return errors.WithStack(err)
		}
		newIosSettings := models.IosSettings{}
		newIosSettings.AppSKU = iosSettings.AppSKU
		newIosSettings.AppleDeveloperAccountEmail = iosSettings.AppleDeveloperAccountEmail
		newIosSettings.ApplSpecificPassword = iosSettings.ApplSpecificPassword
		newIosSettings.SelectedAppStoreProvisioningProfiles = []string{iosSettings.SelectedAppStoreProvisioningProfile}
		newIosSettings.SelectedCodeSigningIdentity = iosSettings.SelectedCodeSigningIdentity
		newIosSettings.IncludeBitCode = iosSettings.IncludeBitCode

		iosSettingsUpdateData, err := json.Marshal(newIosSettings)
		if err != nil {
			logger.Error("Failed to marshal new IosSettings struct", zap.Error(errors.WithStack(err)))
			return errors.WithStack(err)
		}
		appSettings.IosSettingsData = iosSettingsUpdateData
		appSettingsService.Update(&appSettings, []string{"IosSettingsData"})
	}

	return nil
}
