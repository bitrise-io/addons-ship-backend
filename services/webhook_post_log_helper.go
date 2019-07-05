package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/pkg/errors"
)

// WebhookPostLogHelper ...
func WebhookPostLogHelper(env *env.AppEnv, w http.ResponseWriter, r *http.Request, params WebhookPayload) error {
	data, err := parseLogChunkData(params.Data)
	if err != nil {
		return errors.WithStack(err)
	}
	fmt.Printf("%#v\n", data)
	chunkCountRedisKey := fmt.Sprintf("%s_chunk_count", params.TaskID.String())
	latestChunkIndex, err := env.Redis.GetInt64(chunkCountRedisKey)
	if err != nil {
		fmt.Println(err)
		return errors.WithStack(err)
	}
	chunkRedisKey := fmt.Sprintf("%s%d", params.TaskID.String(), latestChunkIndex+1)
	err = env.LogStoreService.Set(chunkRedisKey, models.LogChunk{
		TaskID:  params.TaskID,
		Pos:     data.Position,
		Content: data.Chunk,
	})
	if err != nil {
		return errors.WithStack(err)
	}
	err = env.Redis.Set(chunkCountRedisKey, latestChunkIndex+1, env.RedisExpirationTime)
	if err != nil {
		fmt.Println(err)
		return errors.WithStack(err)
	}
	return httpresponse.RespondWithSuccess(w, httpresponse.StandardErrorRespModel{Message: "ok"})
}

func parseLogChunkData(data interface{}) (LogChunkData, error) {
	var logChunkData LogChunkData
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return LogChunkData{}, err
	}
	err = json.Unmarshal(dataBytes, &logChunkData)
	if err != nil {
		return LogChunkData{}, err
	}
	return logChunkData, nil
}
