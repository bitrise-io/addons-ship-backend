package worker

import (
	"encoding/json"
	"fmt"

	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/addons-ship-backend/redis"
	"github.com/gocraft/work"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var storeLogChunkToRedis = "store_chunk_to_redis"

// StoreLogChunkToRedis ...
func (c *Context) StoreLogChunkToRedis(job *work.Job) error {
	c.env.Logger.Info("[i] Job StoreLogChunkToRedis started")
	taskID := job.ArgString("task_id")
	if taskID == "" {
		c.env.Logger.Error("Failed to get task_id", zap.String("task_id", taskID))
		return errors.New("Failed to get task_id")
	}
	chunkCountRedisKey := fmt.Sprintf("%s_chunk_count", taskID)
	redisClient := redis.New()
	latestChunkIndex, err := redisClient.GetInt64(chunkCountRedisKey)
	if err != nil {
		c.env.Logger.Error("Failed to get chunk count", zap.Error(err))
		return errors.WithStack(err)
	}

	logChunk, err := convertToLogChunk(job.Args["log_chunk"])
	if err != nil {
		c.env.Logger.Error("Failed to get Log Chunk", zap.Error(err), zap.Any("log_chunk", job.Args["log_chunk"]))
		return errors.New("Failed to get Log Chunk")
	}

	chunkRedisKey := fmt.Sprintf("%s%d", taskID, latestChunkIndex+1)
	logStoreService := &models.LogStoreService{Redis: redisClient, Expiration: c.env.RedisExpirationTime}
	err = logStoreService.Set(chunkRedisKey, logChunk)
	if err != nil {
		c.env.Logger.Error("Failed to store Log Chunk in Redis", zap.Error(err))
		return errors.New("Failed to store Log Chunk in Redis")
	}

	err = redisClient.Set(chunkCountRedisKey, latestChunkIndex+1, 0)
	if err != nil {
		c.env.Logger.Error("Failed to set new chunk count", zap.Error(err))
		return errors.WithStack(err)
	}
	c.env.Logger.Info("[i] Job StoreLogChunkToRedis finished")
	return nil
}

func convertToLogChunk(data interface{}) (models.LogChunk, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return models.LogChunk{}, err
	}
	var logChunk models.LogChunk
	err = json.Unmarshal(dataBytes, &logChunk)
	if err != nil {
		return models.LogChunk{}, err
	}
	return logChunk, nil
}
