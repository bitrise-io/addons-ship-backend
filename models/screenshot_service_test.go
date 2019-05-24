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

func Test_ScreenshotService_Create(t *testing.T) {
	dbCloseCallbackMethod := prepareDB(t)
	defer dbCloseCallbackMethod()

	screenshotService := models.ScreenshotService{DB: dataservices.GetDB()}

	t.Run("ok", func(t *testing.T) {
		testScreenshot := &models.Screenshot{
			Filename: "screenshot.png",
			Filesize: 1234,
		}
		createdScreeshot, verrs, err := screenshotService.Create(testScreenshot)
		require.Empty(t, verrs)
		require.NoError(t, err)
		require.False(t, createdScreeshot.ID.String() == "")
		require.False(t, createdScreeshot.CreatedAt.String() == "")
		require.False(t, createdScreeshot.UpdatedAt.String() == "")
	})

	t.Run("when filesize is too big", func(t *testing.T) {
		testScreenshot := &models.Screenshot{
			Filename: "screenshot.png",
			Filesize: services.MaxScreenshotFileByteSize + 1,
		}
		createdScreeshot, verrs, err := screenshotService.Create(testScreenshot)
		require.Equal(t, 1, len(verrs))
		require.Equal(t, "filesize: Must be smaller than 10 megabytes", verrs[0].Error())
		require.NoError(t, err)
		require.Nil(t, createdScreeshot)
	})
}

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
