package models_test

import (
	"testing"
	"time"

	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/c2fo/testify/require"
)

func Test_UpdatableModelService_UpdateData(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		testModel := struct {
			JSONTaggedField string    `json:"json_tagged_field"`
			UpdatedAt       time.Time `json:"updated_at"`
		}{
			JSONTaggedField: "some-test-value",
		}
		testService := models.UpdatableModelService{}

		updateData, err := testService.UpdateData(testModel, []string{"JSONTaggedField"})
		require.NoError(t, err)
		require.Equal(t, map[string]interface{}{
			"json_tagged_field": "some-test-value",
			"updated_at":        time.Time{},
		}, updateData)
	})

	t.Run("when whitelist is empty", func(t *testing.T) {
		testModel := struct {
			JSONTaggedField string `json:"json_tagged_field"`
		}{
			JSONTaggedField: "some-test-value",
		}
		testService := models.UpdatableModelService{}

		updateData, err := testService.UpdateData(testModel, []string{})
		require.EqualError(t, err, "No attributes to update")
		require.Nil(t, updateData)
	})

	t.Run("when struct contains no json tagged field", func(t *testing.T) {
		testModel := struct {
			DBTaggedField string `db:"json_tagged_field"`
		}{
			DBTaggedField: "some-test-value",
		}
		testService := models.UpdatableModelService{}

		updateData, err := testService.UpdateData(testModel, []string{"DBTaggedField"})
		require.EqualError(t, err, "Attribute doesn't have 'json' tag")
		require.Nil(t, updateData)
	})

	t.Run("when struct field has json tag with value '-' and does not have db tag", func(t *testing.T) {
		testModel := struct {
			DBTaggedField string `json:"-" yaml:"json_tagged_field"`
		}{
			DBTaggedField: "some-test-value",
		}
		testService := models.UpdatableModelService{}

		updateData, err := testService.UpdateData(testModel, []string{"DBTaggedField"})
		require.EqualError(t, err, "Attribute doesn't have 'db' tag")
		require.Nil(t, updateData)
	})

	t.Run("when struct does not have field listed in the whitelist", func(t *testing.T) {
		testModel := struct {
			JSONTaggedField string `json:"json_tagged_field"`
		}{
			JSONTaggedField: "some-test-value",
		}
		testService := models.UpdatableModelService{}

		updateData, err := testService.UpdateData(testModel, []string{"NonExistingField"})
		require.EqualError(t, err, "Attribute name doesn't exist in the model")
		require.Nil(t, updateData)
	})
}
