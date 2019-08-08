package services_test

import (
	"github.com/bitrise-io/addons-ship-backend/models"
	uuid "github.com/satori/go.uuid"
)

type testWorkerService struct {
	enqueueStoreLogToAWSFn                  func(uuid.UUID, int64, string, int64) error
	enqueueStoreLogChunkToRedisFn           func(string, models.LogChunk, int64) error
	enqueueCopyUploadablesToNewAppVersionFn func(appVersionFromCopyID, appVersionToCopyID string) error
}

func (s *testWorkerService) EnqueueStoreLogToAWS(publishTaskExternalID uuid.UUID, numberOfLogChunks int64, awsPath string, secondsFromNow int64) error {
	if s.enqueueStoreLogToAWSFn == nil {
		panic("You have to override EnqueueStoreLogToAWS function in tests")
	}
	return s.enqueueStoreLogToAWSFn(publishTaskExternalID, numberOfLogChunks, awsPath, secondsFromNow)
}

func (s *testWorkerService) EnqueueStoreLogChunkToRedis(publishTaskExternalID string, logChunk models.LogChunk, secondsFromNow int64) error {
	if s.enqueueStoreLogChunkToRedisFn == nil {
		panic("You have to override EnqueueStoreLogChunkToRedis function in tests")
	}
	return s.enqueueStoreLogChunkToRedisFn(publishTaskExternalID, logChunk, secondsFromNow)
}

func (s *testWorkerService) EnqueueCopyUploadablesToNewAppVersion(appVersionFromCopyID, appVersionToCopyID string) error {
	if s.enqueueCopyUploadablesToNewAppVersionFn == nil {
		panic("You have to override EnqueueCopyUploadablesToNewAppVersion function in tests")
	}
	return s.enqueueCopyUploadablesToNewAppVersionFn(appVersionFromCopyID, appVersionToCopyID)
}
