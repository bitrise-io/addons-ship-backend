package models_test

import (
	"encoding/json"
	"testing"

	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/c2fo/testify/require"
)

func Test_AppVersion_AppStoreInfo(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		testAppVersion := &models.AppVersion{AppStoreInfoData: json.RawMessage(`{"short_description":"Some shorter description"}`)}
		appStoreInfo, err := testAppVersion.AppStoreInfo()
		require.NoError(t, err)
		require.Equal(t, models.AppStoreInfo{ShortDescription: "Some shorter description"}, appStoreInfo)
	})

	t.Run("error unmarshaling store info", func(t *testing.T) {
		testAppVersion := &models.AppVersion{}
		appStoreInfo, err := testAppVersion.AppStoreInfo()
		require.EqualError(t, err, "unexpected end of JSON input")
		require.Equal(t, models.AppStoreInfo{}, appStoreInfo)
	})
}

func Test_AppVersion_AppInfo(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		testAppVersion := &models.AppVersion{AppInfoData: json.RawMessage(`{"minimum_os":"11.0"}`)}
		appInfo, err := testAppVersion.AppInfo()
		require.NoError(t, err)
		require.Equal(t, models.AppInfo{MinimumOS: "11.0"}, appInfo)
	})

	t.Run("error unmarshaling store info", func(t *testing.T) {
		testAppVersion := &models.AppVersion{}
		appInfo, err := testAppVersion.AppInfo()
		require.EqualError(t, err, "unexpected end of JSON input")
		require.Equal(t, models.AppInfo{}, appInfo)
	})
}

func Test_AppVersion_ProvisioningInfo(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		testAppVersion := &models.AppVersion{ProvisioningInfoData: json.RawMessage(`{"distribution_type":"development"}`)}
		provisioningInfo, err := testAppVersion.ProvisioningInfo()
		require.NoError(t, err)
		require.Equal(t, models.ProvisioningInfo{DistributionType: "development"}, provisioningInfo)
	})

	t.Run("error unmarshaling store info", func(t *testing.T) {
		testAppVersion := &models.AppVersion{}
		provisioningInfo, err := testAppVersion.ProvisioningInfo()
		require.EqualError(t, err, "unexpected end of JSON input")
		require.Equal(t, models.ProvisioningInfo{}, provisioningInfo)
	})
}
