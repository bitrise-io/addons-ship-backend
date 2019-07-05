package dataservices

import (
	"time"

	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/satori/go.uuid"
)

// WorkerService ...
type WorkerService interface {
	EnqueueStoreLogToAWS(publishTaskExternalID uuid.UUID, numberOfLogChunks int64, awsPath string, secondsFromNow time.Duration) error
	EnqueueStoreLogChunkToRedis(publishTaskExternalID string, logChunk models.LogChunk, secondsFromNow time.Duration) error
}
