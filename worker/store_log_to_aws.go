package worker

import (
	"fmt"

	"github.com/gocraft/work"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
)

var storeLogToAWS = "store_log_to_aws"

// StoreLogToAWS ...
func (c *Context) StoreLogToAWS(job *work.Job) error {
	c.env.Logger.Info("[i] Job StoreLogToAWS started")
	eventID := job.ArgString("den_task_id")
	if eventID == (uuid.UUID{}).String() {
		c.env.Logger.Error("Failed to get App Event ID", zap.String("den_task_id", eventID))
		return errors.New("Failed to get App Event ID")
	}
	awsPath := job.ArgString("aws_path")
	if awsPath == "" {
		c.env.Logger.Error("Failed to get AWS path", zap.String("aws_path", awsPath))
		return errors.New("Failed to get AWS path")
	}

	numberOfChunks := job.ArgInt64("number_of_log_chunks")
	content := []byte{}
	for i := int64(1); i <= numberOfChunks; i++ {
		chunk, err := c.env.LogStoreService.Get(fmt.Sprintf("%s%d", eventID, i))
		if err != nil {
			c.env.Logger.Warn("Failed to get log chunk", zap.Error(err))
			continue
		}
		content = append(content, []byte(chunk.Content)...)
	}
	err := c.env.AWS.PutObject(awsPath, content)
	if err != nil {
		c.env.Logger.Error("Failed to save object to AWS", zap.String("aws_path", awsPath), zap.String("content", string(content)))
		return errors.WithStack(err)
	}
	c.env.Logger.Info("[i] Job StoreLogToAWS finished")
	return nil
}

// EnqueueStoreLogToAWS ...
func EnqueueStoreLogToAWS(publishTaskExternalID uuid.UUID, numberOfLogChunks int64, awsPath string) error {
	enqueuer := work.NewEnqueuer(namespace, redisPool)
	_, err := enqueuer.EnqueueUnique(storeLogToAWS, work.Q{
		"den_task_id":          publishTaskExternalID.String(),
		"aws_path":             awsPath,
		"number_of_log_chunks": numberOfLogChunks,
	})
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
