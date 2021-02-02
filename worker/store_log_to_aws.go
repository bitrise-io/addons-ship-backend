package worker

import (
	"fmt"
	"sort"

	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/gocraft/work"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

var storeLogToAWS = "store_log_to_aws"

// StoreLogToAWS ...
func (c *Context) StoreLogToAWS(job *work.Job) error {
	c.env.Logger.Info("[i] Job StoreLogToAWS started")
	appVersionEventID := job.ArgString("app_version_event_id")
	if appVersionEventID == (uuid.UUID{}).String() {
		c.env.Logger.Error("Failed to get App Version Event ID", zap.String("app_version_event_id", appVersionEventID))
		return errors.New("Failed to get App Version Event ID")
	}
	denTaskID := job.ArgString("den_task_id")
	if denTaskID == (uuid.UUID{}).String() {
		c.env.Logger.Error("Failed to get App Event ID", zap.String("den_task_id", denTaskID))
		return errors.New("Failed to get App Event ID")
	}
	awsPath := job.ArgString("aws_path")
	if awsPath == "" {
		c.env.Logger.Error("Failed to get AWS path", zap.String("aws_path", awsPath))
		return errors.New("Failed to get AWS path")
	}

	numberOfChunks := job.ArgInt64("number_of_log_chunks")
	chunks := []models.LogChunk{}
	for i := int64(1); i <= numberOfChunks; i++ {
		chunk, err := c.env.LogStoreService.Get(fmt.Sprintf("%s%d", denTaskID, i))
		if err != nil {
			c.env.Logger.Error("Failed to get log chunk", zap.String("redis_key", fmt.Sprintf("%s%d", denTaskID, i)), zap.Error(err))
			continue
		}
		chunks = append(chunks, chunk)
	}
	sort.Slice(chunks, func(i, j int) bool {
		return chunks[i].Pos < chunks[j].Pos
	})

	content := []byte{}
	for _, chunk := range chunks {
		content = append(content, []byte(chunk.Content)...)
	}
	err := c.env.AWS.PutObject(awsPath, content)
	if err != nil {
		c.env.Logger.Error("Failed to save object to AWS", zap.String("aws_path", awsPath), zap.String("content", string(content)))
		return errors.WithStack(err)
	}

	appVersionEvent, err := c.env.AppVersionEventService.Find(&models.AppVersionEvent{Record: models.Record{ID: uuid.FromStringOrNil(appVersionEventID)}})
	switch {
	case errors.Cause(err) == gorm.ErrRecordNotFound:
		c.env.Logger.Error("App Version Event not found", zap.String("app_version_event_id", appVersionEventID), zap.Error(err))
		return errors.New("App Version Event not found")
	case err != nil:
		c.env.Logger.Error("SQL Error", zap.String("app_version_event_id", appVersionEventID), zap.Error(err))
		return errors.Wrap(err, "SQL Error")
	}

	appVersionEvent.IsLogAvailable = true
	verr, err := c.env.AppVersionEventService.Update(appVersionEvent, []string{"IsLogAvailable"})
	if len(verr) > 0 {
		c.env.Logger.Error("Failed to update App Version Event", zap.String("app_version_event_id", appVersionEventID), zap.Any("validation_errors", verr))
		return errors.New("Failed to update App Version Event")
	}
	if err != nil {
		c.env.Logger.Error("SQL Error", zap.String("app_version_event_id", appVersionEventID), zap.Error(err))
		return errors.Wrap(err, "SQL Error")
	}

	c.env.Logger.Info("[i] Job StoreLogToAWS finished")

	return nil
}
