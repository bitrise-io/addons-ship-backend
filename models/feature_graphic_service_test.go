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

func Test_FeatureGraphicService_Create(t *testing.T) {
	dbCloseCallbackMethod := prepareDB(t)
	defer dbCloseCallbackMethod()

	featureGraphicService := models.FeatureGraphicService{DB: dataservices.GetDB()}

	t.Run("ok", func(t *testing.T) {
		testAppVersion := createTestAppVersion(t, &models.AppVersion{
			AppID:    uuid.NewV4(),
			Platform: "iOS",
		})
		testFeatureGraphic := &models.FeatureGraphic{
			AppVersion: *testAppVersion,
			UploadableObject: models.UploadableObject{
				Filename: "feature_graphic.png",
				Filesize: 1234,
			},
		}

		createdFeatureGraphic, verrs, err := featureGraphicService.Create(testFeatureGraphic)
		require.Empty(t, verrs)
		require.NoError(t, err)
		require.False(t, createdFeatureGraphic.ID.String() == "")
		require.False(t, createdFeatureGraphic.CreatedAt.String() == "")
		require.False(t, createdFeatureGraphic.UpdatedAt.String() == "")
	})

	t.Run("when filesize is too big", func(t *testing.T) {
		testAppVersion := createTestAppVersion(t, &models.AppVersion{
			AppID:    uuid.NewV4(),
			Platform: "iOS",
		})
		testFeatureGraphic := &models.FeatureGraphic{
			AppVersion: *testAppVersion,
			UploadableObject: models.UploadableObject{
				Filename: "feature_graphic.png",
				Filesize: models.MaxFeatureGraphicFileByteSize + 1,
			},
		}
		createdFeatureGraphic, verrs, err := featureGraphicService.Create(testFeatureGraphic)
		require.Equal(t, 1, len(verrs))
		require.Equal(t, "filesize: Must be smaller than 10 megabytes", verrs[0].Error())
		require.NoError(t, err)
		require.Nil(t, createdFeatureGraphic)
	})

	t.Run("when app version already has feature graphic", func(t *testing.T) {
		testAppVersion := createTestAppVersion(t, &models.AppVersion{
			AppID:    uuid.NewV4(),
			Platform: "iOS",
		})
		createTestFeatureGraphic(t, &models.FeatureGraphic{
			AppVersion: *testAppVersion,
			UploadableObject: models.UploadableObject{
				Filename: "feature_graphic.png",
				Filesize: 1234,
			},
		})
		testFeatureGraphic := &models.FeatureGraphic{
			AppVersionID: testAppVersion.ID,
			UploadableObject: models.UploadableObject{
				Filename: "feature_graphic.png",
				Filesize: 1234,
			},
		}

		createdFeatureGraphic, verrs, err := featureGraphicService.Create(testFeatureGraphic)
		require.Equal(t, []error{errors.New("feature_graphics: Maximum count of feature graphics is 1")}, verrs)
		require.NoError(t, err)
		require.Nil(t, createdFeatureGraphic)
	})
}

func Test_FeatureGraphicService_Find(t *testing.T) {
	dbCloseCallbackMethod := prepareDB(t)
	defer dbCloseCallbackMethod()

	featureGraphicService := models.FeatureGraphicService{DB: dataservices.GetDB()}
	testAppVersion := createTestAppVersion(t, &models.AppVersion{
		AppID:    uuid.NewV4(),
		Platform: "iOS",
	})

	testFeatureGraphic := createTestFeatureGraphic(t, &models.FeatureGraphic{
		AppVersion: *testAppVersion,
		UploadableObject: models.UploadableObject{
			Filename: "feature_graphic.png",
			Filesize: 1234,
		},
	})

	t.Run("when querying a feature graphic that belongs to an app version", func(t *testing.T) {
		foundFeatureGraphic, err := featureGraphicService.Find(&models.FeatureGraphic{Record: models.Record{ID: testFeatureGraphic.ID}, AppVersionID: testAppVersion.ID})
		require.NoError(t, err)
		reflect.DeepEqual(testFeatureGraphic, foundFeatureGraphic)
	})

	t.Run("error - when feature graphic is not found", func(t *testing.T) {
		otherTestAppVersion := createTestAppVersion(t, &models.AppVersion{
			AppID:    uuid.NewV4(),
			Platform: "iOS",
		})

		foundFeatureGraphic, err := featureGraphicService.Find(&models.FeatureGraphic{Record: models.Record{ID: testFeatureGraphic.ID}, AppVersionID: otherTestAppVersion.ID})
		require.Equal(t, errors.Cause(err), gorm.ErrRecordNotFound)
		require.Nil(t, foundFeatureGraphic)
	})
}

func Test_FeatureGraphicService_Update(t *testing.T) {
	dbCloseCallbackMethod := prepareDB(t)
	defer dbCloseCallbackMethod()

	featureGraphicService := models.FeatureGraphicService{DB: dataservices.GetDB()}

	t.Run("ok", func(t *testing.T) {
		testAppVersions := []*models.AppVersion{
			createTestAppVersion(t, &models.AppVersion{Platform: "iOS", Version: "v1.0"}),
			createTestAppVersion(t, &models.AppVersion{Platform: "Android", Version: "v1.2"}),
		}
		testFeatureGraphicToUpdate := *createTestFeatureGraphic(t, &models.FeatureGraphic{
			UploadableObject: models.UploadableObject{Filename: "screenshot1.png"},
			AppVersion:       *testAppVersions[0],
		})
		testFeatureGraphicNotToUpdate := createTestFeatureGraphic(t, &models.FeatureGraphic{
			UploadableObject: models.UploadableObject{Filename: "screenshot3.png"},
			AppVersion:       *testAppVersions[1],
		})

		testFeatureGraphicToUpdate.Uploaded = true
		verrs, err := featureGraphicService.Update(testFeatureGraphicToUpdate, []string{"Uploaded"})
		require.Empty(t, verrs)
		require.NoError(t, err)

		t.Log("check if feature graphic got updated")
		foundFeatureGraphic, err := featureGraphicService.Find(&models.FeatureGraphic{Record: models.Record{ID: testFeatureGraphicToUpdate.ID}})
		require.NoError(t, err)
		require.True(t, foundFeatureGraphic.Uploaded)

		t.Log("check if no other feature graphics were updated")
		foundFeatureGraphic, err = featureGraphicService.Find(&models.FeatureGraphic{Record: models.Record{ID: testFeatureGraphicNotToUpdate.ID}})
		require.NoError(t, err)
		require.False(t, foundFeatureGraphic.Uploaded)
	})

	t.Run("when filesize is too big", func(t *testing.T) {
		testAppVersion := createTestAppVersion(t, &models.AppVersion{
			AppID:    uuid.NewV4(),
			Platform: "iOS",
		})
		testFeatureGraphicToUpdate := *createTestFeatureGraphic(t, &models.FeatureGraphic{
			UploadableObject: models.UploadableObject{Filename: "screenshot1.png", Filesize: 1234},
			AppVersion:       *testAppVersion,
		})
		testFeatureGraphicToUpdate.Filesize = models.MaxFeatureGraphicFileByteSize + 1
		verrs, err := featureGraphicService.Update(testFeatureGraphicToUpdate, []string{"Filesize"})
		require.Equal(t, 1, len(verrs))
		require.Equal(t, "filesize: Must be smaller than 10 megabytes", verrs[0].Error())
		require.NoError(t, err)
	})

	t.Run("when trying to update non-existing field", func(t *testing.T) {
		testFeatureGraphicToUpdate := *createTestFeatureGraphic(t, &models.FeatureGraphic{
			UploadableObject: models.UploadableObject{
				Filename: "screenshot1.png",
				Filesize: 1234,
			},
		})
		verrs, err := featureGraphicService.Update(testFeatureGraphicToUpdate, []string{"NonExistingField"})
		require.EqualError(t, err, "Attribute name doesn't exist in the model")
		require.Equal(t, 0, len(verrs))
	})
}
