package models_test

import (
	"testing"

	"github.com/bitrise-io/addons-ship-backend/dataservices"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/c2fo/testify/require"
)

func Test_FeatureGraphicService_Create(t *testing.T) {
	dbCloseCallbackMethod := prepareDB(t)
	defer dbCloseCallbackMethod()

	featureGraphicService := models.FeatureGraphicService{DB: dataservices.GetDB()}

	t.Run("ok", func(t *testing.T) {
		testFeatureGraphic := &models.FeatureGraphic{
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
		testFeatureGraphic := &models.FeatureGraphic{
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
}
