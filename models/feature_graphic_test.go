package models_test

import (
	"testing"

	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/c2fo/testify/require"
	uuid "github.com/satori/go.uuid"
)

func Test_FeatureGraphic_AWSPath(t *testing.T) {
	testFeatureGraphic := models.FeatureGraphic{
		Record:           models.Record{ID: uuid.FromStringOrNil("33c7223f-2203-4109-b439-6026e7a374c9")},
		UploadableObject: models.UploadableObject{Filename: "feature_graphic.png"},
		AppVersion: models.AppVersion{
			Record: models.Record{
				ID: uuid.FromStringOrNil("de438ddc-98e5-4226-a5f4-fd2d53474879"),
			},
			App: models.App{AppSlug: "test-app-slug"},
		},
	}

	require.Equal(t, "test-app-slug/de438ddc-98e5-4226-a5f4-fd2d53474879/33c7223f-2203-4109-b439-6026e7a374c9.png", testFeatureGraphic.AWSPath())
}
