package services_test

import (
	"time"

	"github.com/bitrise-io/addons-ship-backend/models"
	uuid "github.com/satori/go.uuid"
)

type testWorkerService struct {
	enqueueStoreLogToAWSFn        func(uuid.UUID, int64, string, time.Duration) error
	enqueueStoreLogChunkToRedisFn func(string, models.LogChunk, time.Duration) error
}

func (s *testWorkerService) EnqueueStoreLogToAWS(publishTaskExternalID uuid.UUID, numberOfLogChunks int64, awsPath string, secondsFromNow time.Duration) error {
	if s.enqueueStoreLogToAWSFn == nil {
		panic("You have to override EnqueueStoreLogToAWS function in tests")
	}
	return s.enqueueStoreLogToAWSFn(publishTaskExternalID, numberOfLogChunks, awsPath, secondsFromNow)
}

func (s *testWorkerService) EnqueueStoreLogChunkToRedis(publishTaskExternalID string, logChunk models.LogChunk, secondsFromNow time.Duration) error {
	if s.enqueueStoreLogChunkToRedisFn == nil {
		panic("You have to override EnqueueStoreLogChunkToRedis function in tests")
	}
	return s.enqueueStoreLogChunkToRedisFn(publishTaskExternalID, logChunk, secondsFromNow)
}
