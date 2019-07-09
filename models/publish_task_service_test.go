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
	"github.com/satori/go.uuid"
)

func Test_PublishTaskService_Create(t *testing.T) {
	dbCloseCallbackMethod := prepareDB(t)
	defer dbCloseCallbackMethod()

	appVersionEventService := models.PublishTaskService{DB: dataservices.GetDB()}
	testPublishTask := &models.PublishTask{TaskID: uuid.FromStringOrNil("600450c0-ca1a-4d01-afee-d184722cc63a")}

	createdPublishTask, err := appVersionEventService.Create(testPublishTask)
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
	testPublishTask1 := createTestPublishTask(t, &models.PublishTask{TaskID: uuid.FromStringOrNil("600450c0-ca1a-4d01-afee-d184722cc63a"), AppVersion: *testAppVersion})
	testPublishTask2 := createTestPublishTask(t, &models.PublishTask{TaskID: uuid.FromStringOrNil("188dfbe5-4505-44d9-ae95-9d7aa84dd0be"), AppVersion: *testAppVersion})

	t.Run("when querying publish task that belongs to an app", func(t *testing.T) {
		foundPublishTask, err := publishTaskService.Find(&models.PublishTask{TaskID: uuid.FromStringOrNil("600450c0-ca1a-4d01-afee-d184722cc63a")})
		require.NoError(t, err)
		reflect.DeepEqual(testPublishTask1, foundPublishTask)

		foundPublishTask, err = publishTaskService.Find(&models.PublishTask{TaskID: uuid.FromStringOrNil("188dfbe5-4505-44d9-ae95-9d7aa84dd0be")})
		require.NoError(t, err)
		reflect.DeepEqual(testPublishTask2, foundPublishTask)
	})

	t.Run("error - when app event is not found", func(t *testing.T) {
		foundPublishTask, err := publishTaskService.Find(&models.PublishTask{TaskID: uuid.FromStringOrNil("841b360b-d467-481c-8813-5937dd5c4d8e")})
		require.Equal(t, errors.Cause(err), gorm.ErrRecordNotFound)
		require.Nil(t, foundPublishTask)
	})
}
