package worker

import (
	"fmt"
	"sort"

	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/gocraft/work"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
)

var storeLogToAWS = "store_log_to_aws"

// StoreLogToAWS ...
func (c *Context) StoreLogToAWS(job *work.Job) error {
	c.env.Logger.Info("[i] Job StoreLogToAWS started")
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
		if chunks[i].Pos < chunks[j].Pos {
			return true
		}
		return false
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
	c.env.Logger.Info("[i] Job StoreLogToAWS finished")

	return nil
}
