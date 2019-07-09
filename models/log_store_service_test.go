package models_test

import (
	"testing"

	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/addons-ship-backend/redis"
	"github.com/c2fo/testify/require"
	"github.com/pkg/errors"
)

func Test_Get(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		logChunkStr := `{"content":"Some content"}`
		logStore := models.LogStoreService{
			Redis: &redis.Mock{
				GetStringFn: func(key string) (string, error) {
					require.Equal(t, "TEST_KEY", key)
					return logChunkStr, nil
				},
			},
		}
		foundChunk, err := logStore.Get("TEST_KEY")
		require.NoError(t, err)
		require.Equal(t, models.LogChunk{Content: "Some content"}, foundChunk)
	})

	t.Run("when format of value in Redis is invalid", func(t *testing.T) {
		logChunkStr := `invalid JSON`
		logStore := models.LogStoreService{
			Redis: &redis.Mock{
				GetStringFn: func(key string) (string, error) {
					require.Equal(t, "TEST_KEY", key)
					return logChunkStr, nil
				},
			},
		}
		foundChunk, err := logStore.Get("TEST_KEY")
		require.EqualError(t, err, "invalid character 'i' looking for beginning of value")
		require.Equal(t, models.LogChunk{}, foundChunk)
	})

	t.Run("when error happens in Redis", func(t *testing.T) {
		logStore := models.LogStoreService{
			Redis: &redis.Mock{
				GetStringFn: func(key string) (string, error) {
					require.Equal(t, "TEST_KEY", key)
					return "", errors.New("SOME-REDIS-ERROR")
				},
			},
		}
		foundChunk, err := logStore.Get("TEST_KEY")
		require.EqualError(t, err, "SOME-REDIS-ERROR")
		require.Equal(t, models.LogChunk{}, foundChunk)
	})
}

func Test_Set(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		logChunkStr := `{"id":"00000000-0000-0000-0000-000000000000","task_id":"00000000-0000-0000-0000-000000000000","pos":0,"content":"Some content","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}`
		logStore := models.LogStoreService{
			Redis: &redis.Mock{
				SetFn: func(key string, value interface{}, ttl int) error {
					require.Equal(t, "TEST_KEY", key)
					require.Equal(t, value, logChunkStr)
					return nil
				},
			},
		}
		err := logStore.Set("TEST_KEY", models.LogChunk{Content: "Some content"})
		require.NoError(t, err)
	})

	t.Run("when error happens in Redis", func(t *testing.T) {
		logChunkStr := `{"id":"00000000-0000-0000-0000-000000000000","task_id":"00000000-0000-0000-0000-000000000000","pos":0,"content":"Some content","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}`
		logStore := models.LogStoreService{
			Redis: &redis.Mock{
				SetFn: func(key string, value interface{}, ttl int) error {
					require.Equal(t, "TEST_KEY", key)
					require.Equal(t, value, logChunkStr)
					return errors.New("SOME-REDIS-ERROR")
				},
			},
		}
		err := logStore.Set("TEST_KEY", models.LogChunk{Content: "Some content"})
		require.EqualError(t, err, "SOME-REDIS-ERROR")
	})
}
