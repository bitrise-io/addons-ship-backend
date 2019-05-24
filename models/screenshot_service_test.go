// +build database

package models_test

import (
	"reflect"
	"testing"

	"github.com/bitrise-io/addons-ship-backend/dataservices"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/c2fo/testify/require"
	uuid "github.com/satori/go.uuid"
)

func Test_ScreenshotService_FindAll(t *testing.T) {
	dbCloseCallbackMethod := prepareDB(t)
	defer dbCloseCallbackMethod()

	screenshotService := models.ScreenshotService{DB: dataservices.GetDB()}
	testAppVersionIOS := createTestAppVersion(t, &models.AppVersion{
		AppID:    uuid.NewV4(),
		Platform: "iOS",
	})
	testAppVersionAndroid := createTestAppVersion(t, &models.AppVersion{
		AppID:    uuid.NewV4(),
		Platform: "android",
	})
	testScreenshot1 := createTestScreenshot(t, &models.Screenshot{
		AppVersion: *testAppVersionIOS,
		DeviceType: "iPhone XS Max",
		ScreenSize: "6.5 inch",
	})
	testScreenshot2 := createTestScreenshot(t, &models.Screenshot{
		AppVersion: *testAppVersionIOS,
		DeviceType: "iPad Pro",
		ScreenSize: "12.9 inch",
	})
	createTestScreenshot(t, &models.Screenshot{
		AppVersion: *testAppVersionAndroid,
		DeviceType: "Google Pixel 3",
		ScreenSize: "5.5 inch",
	})

	t.Run("when query all screenshots of test iOS app version", func(t *testing.T) {
		foundScreenshots, err := screenshotService.FindAll(testAppVersionIOS)
		require.NoError(t, err)
		reflect.DeepEqual([]models.Screenshot{*testScreenshot2, *testScreenshot1}, foundScreenshots)
	})
}
