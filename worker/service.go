package worker

import (
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/gocraft/work"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

// Service ...
type Service struct{}

// EnqueueStoreLogToAWS ...
func (*Service) EnqueueStoreLogToAWS(appVersionEventID, publishTaskExternalID uuid.UUID, numberOfLogChunks int64, awsPath string, secondsFromNow int64) error {
	enqueuer := work.NewEnqueuer(namespace, redisPool)
	var err error
	jobParams := work.Q{
		"app_version_event_id": appVersionEventID,
		"den_task_id":          publishTaskExternalID.String(),
		"aws_path":             awsPath,
		"number_of_log_chunks": numberOfLogChunks,
	}
	if secondsFromNow == 0 {
		_, err = enqueuer.EnqueueUnique(storeLogToAWS, jobParams)
	} else {
		_, err = enqueuer.EnqueueUniqueIn(storeLogToAWS, secondsFromNow, jobParams)
	}
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// EnqueueStoreLogChunkToRedis ...
func (*Service) EnqueueStoreLogChunkToRedis(publishTaskExternalID string, logChunk models.LogChunk, secondsFromNow int64) error {
	enqueuer := work.NewEnqueuer(namespace, redisPool)
	var err error
	jobParams := work.Q{
		"task_id":   publishTaskExternalID,
		"log_chunk": logChunk,
	}
	if secondsFromNow == 0 {
		_, err = enqueuer.EnqueueUnique(storeLogChunkToRedis, jobParams)
	} else {
		_, err = enqueuer.EnqueueUniqueIn(storeLogChunkToRedis, secondsFromNow, jobParams)
	}
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// EnqueueCopyUploadablesToNewAppVersion ...
func (*Service) EnqueueCopyUploadablesToNewAppVersion(appVersionFromCopyID, appVersionToCopyID string) error {
	enqueuer := work.NewEnqueuer(namespace, redisPool)
	var err error
	jobParams := work.Q{
		"from_id": appVersionFromCopyID,
		"to_id":   appVersionToCopyID,
	}

	_, err = enqueuer.EnqueueUnique(copyUploadablesToNewAppVersion, jobParams)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
