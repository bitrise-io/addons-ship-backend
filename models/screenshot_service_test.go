// +build database

package models_test

import (
	"reflect"
	"testing"

	"github.com/bitrise-io/addons-ship-backend/dataservices"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/c2fo/testify/require"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func Test_ScreenshotService_BatchCreate(t *testing.T) {
	dbCloseCallbackMethod := prepareDB(t)
	defer dbCloseCallbackMethod()

	screenshotService := models.ScreenshotService{DB: dataservices.GetDB()}
	testApp := createTestApp(t, &models.App{AppSlug: "test-app-slug"})
	testAppVersion := createTestAppVersion(t, &models.AppVersion{AppID: testApp.ID, Platform: "ios"})

	t.Run("ok", func(t *testing.T) {
		testScreenshots := []*models.Screenshot{
			&models.Screenshot{
				UploadableObject: models.UploadableObject{
					Filename: "screenshot.png",
					Filesize: 1234,
				},
				AppVersionID: testAppVersion.ID,
			},
		}
		createdScreeshots, verrs, err := screenshotService.BatchCreate(testScreenshots)
		require.Empty(t, verrs)
		require.NoError(t, err)
		require.False(t, createdScreeshots[0].ID.String() == "")
		require.False(t, createdScreeshots[0].CreatedAt.String() == "")
		require.False(t, createdScreeshots[0].UpdatedAt.String() == "")
		require.Equal(t, "ios", createdScreeshots[0].AppVersion.Platform)
		require.Equal(t, "test-app-slug", createdScreeshots[0].AppVersion.App.AppSlug)
	})

	t.Run("when filesize is too big", func(t *testing.T) {
		testScreenshot := []*models.Screenshot{
			&models.Screenshot{
				UploadableObject: models.UploadableObject{
					Filename: "screenshot.png",
					Filesize: models.MaxScreenshotFileByteSize + 1,
				},
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
				UploadableObject: models.UploadableObject{
					Filename: "screenshot.png",
					Filesize: 1234,
				},
			},
			&models.Screenshot{
				AppVersion: *testAppVersion,
				UploadableObject: models.UploadableObject{
					Filename: "screenshot.png",
					Filesize: models.MaxScreenshotFileByteSize + 1,
				},
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

func Test_ScreenshotService_Find(t *testing.T) {
	dbCloseCallbackMethod := prepareDB(t)
	defer dbCloseCallbackMethod()

	screenshotService := models.ScreenshotService{DB: dataservices.GetDB()}
	testAppVersion := createTestAppVersion(t, &models.AppVersion{
		AppID:    uuid.NewV4(),
		Platform: "iOS",
	})

	testScreenshot := createTestScreenshot(t, &models.Screenshot{
		AppVersion: *testAppVersion,
		DeviceType: "iPhone XS Max",
		ScreenSize: "6.5 inch",
	})

	t.Run("when querying a screenshot that belongs to an app version", func(t *testing.T) {
		foundScreenshot, err := screenshotService.Find(&models.Screenshot{Record: models.Record{ID: testScreenshot.ID}, AppVersionID: testAppVersion.ID})
		require.NoError(t, err)
		reflect.DeepEqual(testScreenshot, foundScreenshot)
	})

	t.Run("error - when screenshot is not found", func(t *testing.T) {
		otherTestAppVersion := createTestAppVersion(t, &models.AppVersion{
			AppID:    uuid.NewV4(),
			Platform: "iOS",
		})

		foundScreenshot, err := screenshotService.Find(&models.Screenshot{Record: models.Record{ID: testScreenshot.ID}, AppVersionID: otherTestAppVersion.ID})
		require.Equal(t, errors.Cause(err), gorm.ErrRecordNotFound)
		require.Nil(t, foundScreenshot)
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
		testScreenshotsOfVersion1 := []models.Screenshot{
			*createTestScreenshot(t, &models.Screenshot{
				UploadableObject: models.UploadableObject{Filename: "screenshot1.png"},
				AppVersion:       *testAppVersions[0],
			}),
			*createTestScreenshot(t, &models.Screenshot{
				UploadableObject: models.UploadableObject{Filename: "screenshot2.png"},
				AppVersion:       *testAppVersions[0],
			}),
		}
		createTestScreenshot(t, &models.Screenshot{
			UploadableObject: models.UploadableObject{Filename: "screenshot3.png"},
			AppVersion:       *testAppVersions[1],
		})

		testScreenshotsOfVersion1[0].Uploaded = true
		testScreenshotsOfVersion1[1].Uploaded = true
		verrs, err := screenshotService.BatchUpdate(testScreenshotsOfVersion1, []string{"Uploaded"})
		require.Empty(t, verrs)
		require.NoError(t, err)

		t.Log("check if screenshots got updated")
		foundScreenshots, err := screenshotService.FindAll(testAppVersions[0])
		require.NoError(t, err)
		require.True(t, foundScreenshots[0].Uploaded)
		require.True(t, foundScreenshots[1].Uploaded)

		t.Log("check if no other screenshots were updated")
		foundScreenshots, err = screenshotService.FindAll(testAppVersions[1])
		require.NoError(t, err)
		require.False(t, foundScreenshots[0].Uploaded)
	})

	t.Run("when filesize is too big", func(t *testing.T) {
		testScreenshots := []models.Screenshot{
			*createTestScreenshot(t, &models.Screenshot{
				UploadableObject: models.UploadableObject{
					Filename: "screenshot1.png",
					Filesize: 1234,
				},
			}),
		}
		testScreenshots[0].Filesize = models.MaxScreenshotFileByteSize + 1
		verrs, err := screenshotService.BatchUpdate(testScreenshots, []string{"Filesize"})
		require.Equal(t, 1, len(verrs))
		require.Equal(t, "filesize: Must be smaller than 10 megabytes", verrs[0].Error())
		require.NoError(t, err)
	})

	t.Run("when trying to update non-existing field", func(t *testing.T) {
		testScreenshots := []models.Screenshot{
			*createTestScreenshot(t, &models.Screenshot{
				UploadableObject: models.UploadableObject{
					Filename: "screenshot1.png",
					Filesize: 1234,
				},
			}),
		}
		verrs, err := screenshotService.BatchUpdate(testScreenshots, []string{"NonExistingField"})
		require.EqualError(t, err, "Attribute name doesn't exist in the model")
		require.Equal(t, 0, len(verrs))
	})
}

func Test_ScreenshotService_Delete(t *testing.T) {
	dbCloseCallbackMethod := prepareDB(t)
	defer dbCloseCallbackMethod()

	screenshotService := models.ScreenshotService{DB: dataservices.GetDB()}

	testScreenshot := createTestScreenshot(t, &models.Screenshot{
		DeviceType: "iPhone XS Max",
		ScreenSize: "6.5 inch",
	})

	t.Run("when deleting a screenshot", func(t *testing.T) {
		err := screenshotService.Delete(&models.Screenshot{Record: models.Record{ID: testScreenshot.ID}})
		require.NoError(t, err)
	})

	t.Run("error - when screenshot is not found", func(t *testing.T) {
		err := screenshotService.Delete(&models.Screenshot{Record: models.Record{ID: uuid.NewV4()}})

		require.Equal(t, err, gorm.ErrRecordNotFound)
	})
}
