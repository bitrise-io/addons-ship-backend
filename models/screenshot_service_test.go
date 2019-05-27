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

func Test_ScreenshotService_BatchCreate(t *testing.T) {
	dbCloseCallbackMethod := prepareDB(t)
	defer dbCloseCallbackMethod()

	screenshotService := models.ScreenshotService{DB: dataservices.GetDB()}

	t.Run("ok", func(t *testing.T) {
		testScreenshots := []*models.Screenshot{
			&models.Screenshot{
				Filename: "screenshot.png",
				Filesize: 1234,
			},
		}
		createdScreeshots, verrs, err := screenshotService.BatchCreate(testScreenshots)
		require.Empty(t, verrs)
		require.NoError(t, err)
		require.False(t, createdScreeshots[0].ID.String() == "")
		require.False(t, createdScreeshots[0].CreatedAt.String() == "")
		require.False(t, createdScreeshots[0].UpdatedAt.String() == "")
	})

	t.Run("when filesize is too big", func(t *testing.T) {
		testScreenshot := []*models.Screenshot{
			&models.Screenshot{
				Filename: "screenshot.png",
				Filesize: models.MaxScreenshotFileByteSize + 1,
			},
		}
		createdScreeshot, verrs, err := screenshotService.BatchCreate(testScreenshot)
		require.Equal(t, 1, len(verrs))
		require.Equal(t, "filesize: Must be smaller than 10 megabytes", verrs[0].Error())
		require.NoError(t, err)
		require.Nil(t, createdScreeshot)
	})

	t.Run("when error happens at creation of any screenshot, transaction gets rolled back", func(t *testing.T) {
		testAppVersion := createTestAppVersion(t, &models.AppVersion{
			Platform: "iOS",
			Version:  "v1.0",
		})
		testScreenshots := []*models.Screenshot{
			&models.Screenshot{
				AppVersion: *testAppVersion,
				Filename:   "screenshot.png",
				Filesize:   1234,
			},
			&models.Screenshot{
				AppVersion: *testAppVersion,
				Filename:   "screenshot.png",
				Filesize:   models.MaxScreenshotFileByteSize + 1,
			},
		}
		createdScreeshots, verrs, err := screenshotService.BatchCreate(testScreenshots)
		require.Equal(t, 1, len(verrs))
		require.Equal(t, "filesize: Must be smaller than 10 megabytes", verrs[0].Error())
		require.NoError(t, err)
		require.Empty(t, createdScreeshots)

		foundScreenshots, err := screenshotService.FindAll(testAppVersion)
		require.NoError(t, err)
		require.Empty(t, foundScreenshots)
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

func Test_ScreenshotService_BatchUpdate(t *testing.T) {
	dbCloseCallbackMethod := prepareDB(t)
	defer dbCloseCallbackMethod()

	screenshotService := models.ScreenshotService{DB: dataservices.GetDB()}

	t.Run("ok", func(t *testing.T) {
		testAppVersions := []*models.AppVersion{
			createTestAppVersion(t, &models.AppVersion{Platform: "iOS", Version: "v1.0"}),
			createTestAppVersion(t, &models.AppVersion{Platform: "Android", Version: "v1.2"}),
		}
		testScreenshotsOfVersion1 := []*models.Screenshot{
			createTestScreenshot(t, &models.Screenshot{
				Filename:   "screenshot1.png",
				AppVersion: *testAppVersions[0],
			}),
			createTestScreenshot(t, &models.Screenshot{
				Filename:   "screenshot2.png",
				AppVersion: *testAppVersions[0],
			}),
		}
		testScreenshotsOfVersion2 := createTestScreenshot(t, &models.Screenshot{
			Filename:   "screenshot3.png",
			AppVersion: *testAppVersions[1],
		})
		testScreenshotsOfVersion1[0].Uploaded = true
		testScreenshotsOfVersion1[1].Uploaded = true
		updatedScreenshots, verrs, err := screenshotService.BatchUpdate(testScreenshotsOfVersion1, []string{"Uploaded"})
		require.Empty(t, verrs)
		require.NoError(t, err)
		require.True(t, updatedScreenshots[0].Uploaded)
		require.True(t, updatedScreenshots[1].Uploaded)
	})

	// t.Run("when filesize is too big", func(t *testing.T) {
	// 	testScreenshot := &models.Screenshot{
	// 		Filename: "screenshot.png",
	// 		Filesize: models.MaxScreenshotFileByteSize + 1,
	// 	}
	// 	createdScreeshot, verrs, err := screenshotService.Create(testScreenshot)
	// 	require.Equal(t, 1, len(verrs))
	// 	require.Equal(t, "filesize: Must be smaller than 10 megabytes", verrs[0].Error())
	// 	require.NoError(t, err)
	// 	require.Nil(t, createdScreeshot)
	// })

	// t.Run("when error happens at creation of any screenshot, transaction get rolled back", func(t *testing.T) {
	// 	testAppVersion := createTestAppVersion(t, &models.AppVersion{
	// 		Platform: "iOS",
	// 		Version:  "v1.0",
	// 	})
	// 	testScreenshots := []*models.Screenshot{
	// 		&models.Screenshot{
	// 			AppVersion: *testAppVersion,
	// 			Filename:   "screenshot.png",
	// 			Filesize:   1234,
	// 		},
	// 		&models.Screenshot{
	// 			AppVersion: *testAppVersion,
	// 			Filename:   "screenshot.png",
	// 			Filesize:   models.MaxScreenshotFileByteSize + 1,
	// 		},
	// 	}
	// 	createdScreeshots, verrs, err := screenshotService.BatchCreate(testScreenshots)
	// 	require.Equal(t, 1, len(verrs))
	// 	require.Equal(t, "filesize: Must be smaller than 10 megabytes", verrs[0].Error())
	// 	require.NoError(t, err)
	// 	require.Empty(t, createdScreeshots)

	// 	foundScreenshots, err := screenshotService.FindAll(testAppVersion)
	// 	require.NoError(t, err)
	// 	require.Empty(t, foundScreenshots)
	// })
}
