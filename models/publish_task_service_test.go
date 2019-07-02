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
)

func Test_PublishTaskService_Create(t *testing.T) {
	dbCloseCallbackMethod := prepareDB(t)
	defer dbCloseCallbackMethod()

	appEventService := models.PublishTaskService{DB: dataservices.GetDB()}
	testPublishTask := &models.PublishTask{TaskID: "abcd-efgh-1234"}

	createdPublishTask, err := appEventService.Create(testPublishTask)
	require.NoError(t, err)
	require.False(t, createdPublishTask.ID.String() == "")
	require.False(t, createdPublishTask.CreatedAt.String() == "")
	require.False(t, createdPublishTask.UpdatedAt.String() == "")
}

func Test_PublishTaskService_Find(t *testing.T) {
	dbCloseCallbackMethod := prepareDB(t)
	defer dbCloseCallbackMethod()

	publishTaskService := models.PublishTaskService{DB: dataservices.GetDB()}
	testAppVersion := createTestAppVersion(t, &models.AppVersion{Version: "v1.0", Platform: "ios"})
	testPublishTask1 := createTestPublishTask(t, &models.PublishTask{TaskID: "abcd-efcg-1234", AppVersion: *testAppVersion})
	testPublishTask2 := createTestPublishTask(t, &models.PublishTask{TaskID: "abcd-efcg-5678", AppVersion: *testAppVersion})

	t.Run("when querying publish task that belongs to an app", func(t *testing.T) {
		foundPublishTask, err := publishTaskService.Find(&models.PublishTask{TaskID: "abcd-efcg-1234"})
		require.NoError(t, err)
		reflect.DeepEqual(testPublishTask1, foundPublishTask)

		foundPublishTask, err = publishTaskService.Find(&models.PublishTask{TaskID: "abcd-efcg-5678"})
		require.NoError(t, err)
		reflect.DeepEqual(testPublishTask2, foundPublishTask)
	})

	t.Run("error - when app event is not found", func(t *testing.T) {
		foundPublishTask, err := publishTaskService.Find(&models.PublishTask{TaskID: "invalid-task-id"})
		require.Equal(t, errors.Cause(err), gorm.ErrRecordNotFound)
		require.Nil(t, foundPublishTask)
	})
}
