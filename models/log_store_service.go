package models

import (
	"encoding/json"

	"github.com/bitrise-io/addons-ship-backend/redis"
)

// LogStoreService ...
type LogStoreService struct {
	Redis      *redis.Client
	Expiration int
}

// Get ...
func (s *LogStoreService) Get(key string) (LogChunk, error) {
	chunkStr, err := s.Redis.Get(key)
	if err != nil {
		return LogChunk{}, err
	}
	var chunk LogChunk
	err = json.Unmarshal([]byte(chunkStr), &chunk)
	if err != nil {
		return LogChunk{}, err
	}
	return chunk, nil
}

// Set ...
func (s *LogStoreService) Set(key string, chunk LogChunk) error {
	chunkBytes, err := json.Marshal(chunk)
	if err != nil {
		return err
	}
	err = s.Redis.Set(key, string(chunkBytes), s.Expiration)
	if err != nil {
		return err
	}
	return nil
}
