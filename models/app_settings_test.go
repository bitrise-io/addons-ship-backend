package models_test

import (
	"encoding/json"
	"testing"

	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/c2fo/testify/require"
)

func Test_IosSettings(t *testing.T) {
	t.Run("when ios settings is valid", func(t *testing.T) {
		validIosSettings := models.IosSettings{AppSKU: "2019061"}
		require.True(t, validIosSettings.Valid())
	})

	t.Run("when ios settings is invalid", func(t *testing.T) {
		validIosSettings := models.IosSettings{}
		require.False(t, validIosSettings.Valid())
	})
}

func Test_IosSettings_ValidateSelectedProvisioningProfileSlugs(t *testing.T) {
	t.Run("when selected slugs' list contains not existing", func(t *testing.T) {
		iosSetting := models.IosSettings{SelectedAppStoreProvisioningProfiles: []string{"prov-1-slug", "prov-2-slug", "prov-3-slug"}}
		iosSetting.ValidateSelectedProvisioningProfileSlugs([]string{"prov-1-slug", "prov-3-slug"})
		require.Equal(t, []string{"prov-1-slug", "prov-3-slug"}, iosSetting.SelectedAppStoreProvisioningProfiles)
	})

	t.Run("when selected slugs' list contains only existing ones", func(t *testing.T) {
		iosSetting := models.IosSettings{SelectedAppStoreProvisioningProfiles: []string{"prov-1-slug", "prov-2-slug"}}
		iosSetting.ValidateSelectedProvisioningProfileSlugs([]string{"prov-1-slug", "prov-2-slug", "prov-3-slug"})
		require.Equal(t, []string{"prov-1-slug", "prov-2-slug"}, iosSetting.SelectedAppStoreProvisioningProfiles)
	})
}

func Test_AppSettings_IosSettings(t *testing.T) {
	t.Run("when ios settings is valid", func(t *testing.T) {
		testAppSettings := models.AppSettings{IosSettingsData: json.RawMessage(`{"app_sku":"2019061"}`)}
		iosSettings, err := testAppSettings.IosSettings()
		require.NoError(t, err)
		require.Equal(t, models.IosSettings{AppSKU: "2019061"}, iosSettings)
	})

	t.Run("when ios settings is invalid", func(t *testing.T) {
		testAppSettings := models.AppSettings{IosSettingsData: json.RawMessage(`invalid json`)}
		iosSettings, err := testAppSettings.IosSettings()
		require.EqualError(t, err, "invalid character 'i' looking for beginning of value")
		require.Equal(t, models.IosSettings{}, iosSettings)
	})
}

func Test_AndroidSettings(t *testing.T) {
	t.Run("when ios settings is valid", func(t *testing.T) {
		validAndroidSettings := models.AndroidSettings{Track: "2019061"}
		require.True(t, validAndroidSettings.Valid())
	})

	t.Run("when ios settings is invalid", func(t *testing.T) {
		validAndroidSettings := models.AndroidSettings{}
		require.False(t, validAndroidSettings.Valid())
	})
}

func Test_AppSettings_AndroidSettings(t *testing.T) {
	t.Run("when ios settings is valid", func(t *testing.T) {
		testAppSettings := models.AppSettings{AndroidSettingsData: json.RawMessage(`{"track":"2019061"}`)}
		iosSettings, err := testAppSettings.AndroidSettings()
		require.NoError(t, err)
		require.Equal(t, models.AndroidSettings{Track: "2019061"}, iosSettings)
	})

	t.Run("when ios settings is invalid", func(t *testing.T) {
		testAppSettings := models.AppSettings{AndroidSettingsData: json.RawMessage(`invalid json`)}
		iosSettings, err := testAppSettings.AndroidSettings()
		require.EqualError(t, err, "invalid character 'i' looking for beginning of value")
		require.Equal(t, models.AndroidSettings{}, iosSettings)
	})
}
