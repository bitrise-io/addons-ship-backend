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
