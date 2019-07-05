package worker

import (
	"time"

	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/gocraft/work"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

// Service ...
type Service struct{}

// EnqueueStoreLogToAWS ...
func (*Service) EnqueueStoreLogToAWS(publishTaskExternalID uuid.UUID, numberOfLogChunks int64, awsPath string, secondsFromNow time.Duration) error {
	enqueuer := work.NewEnqueuer(namespace, redisPool)
	var err error
	jobParams := work.Q{
		"den_task_id":          publishTaskExternalID.String(),
		"aws_path":             awsPath,
		"number_of_log_chunks": numberOfLogChunks,
	}
	if secondsFromNow == 0 {
		_, err = enqueuer.EnqueueUnique(storeLogToAWS, jobParams)
	} else {
		_, err = enqueuer.EnqueueUniqueIn(storeLogToAWS, int64(secondsFromNow), jobParams)
	}
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// EnqueueStoreLogChunkToRedis ...
func (*Service) EnqueueStoreLogChunkToRedis(publishTaskExternalID string, logChunk models.LogChunk, secondsFromNow time.Duration) error {
	enqueuer := work.NewEnqueuer(namespace, redisPool)
	var err error
	jobParams := work.Q{
		"task_id":   publishTaskExternalID,
		"log_chunk": logChunk,
	}
	if secondsFromNow == 0 {
		_, err = enqueuer.EnqueueUnique(storeLogChunkToRedis, jobParams)
	} else {
		_, err = enqueuer.EnqueueUniqueIn(storeLogChunkToRedis, int64(secondsFromNow), jobParams)
	}
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
