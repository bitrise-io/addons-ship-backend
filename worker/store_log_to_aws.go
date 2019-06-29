package worker

import (
	"fmt"

	"github.com/gocraft/work"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

var storeLogToAWS = "store_log_to_aws"

// StoreLogToAWS ...
func (c *Context) StoreLogToAWS(job *work.Job) error {
	eventID := job.ArgString("event_id")
	if eventID == (uuid.UUID{}).String() {
		return errors.New("Failed to get App Event ID")
	}
	awsPath := job.ArgString("aws_path")
	if awsPath == "" {
		return errors.New("Failed to get AWS path")
	}
	numberOfChunks := job.ArgInt64("number_of_log_chunks")
	content := []byte{}
	for i := int64(1); i <= numberOfChunks; i++ {
		chunk, err := c.env.LogStoreService.Get(fmt.Sprintf("%s%d", eventID, i))
		if err != nil {
			return errors.WithStack(err)
		}
		fmt.Println(chunk.Content)
		content = append(content, []byte(chunk.Content)...)
	}
	err := c.env.AWS.PutObject(awsPath, content)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// EnqueueStoreLogToAWS ...
func EnqueueStoreLogToAWS(appEventID uuid.UUID, numberOfLogChunks int64, awsPath string) error {
	enqueuer := work.NewEnqueuer(namespace, redisPool)
	_, err := enqueuer.EnqueueUnique(storeLogToAWS, work.Q{
		"event_id":             appEventID.String(),
		"aws_path":             awsPath,
		"number_of_log_chunks": numberOfLogChunks,
	})
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
