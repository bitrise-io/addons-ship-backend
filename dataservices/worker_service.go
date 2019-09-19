package dataservices

import (
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/satori/go.uuid"
)

// WorkerService ...
type WorkerService interface {
	EnqueueStoreLogToAWS(appVersionEventID, publishTaskExternalID uuid.UUID, numberOfLogChunks int64, awsPath string, secondsFromNow int64) error
	EnqueueStoreLogChunkToRedis(publishTaskExternalID string, logChunk models.LogChunk, secondsFromNow int64) error
	EnqueueCopyUploadablesToNewAppVersion(appVersionFromCopyID, appVersionToCopyID string) error
}
