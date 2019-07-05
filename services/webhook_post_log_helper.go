package services

import (
	"encoding/json"
	"net/http"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/pkg/errors"
)

func webhookPostLogHelper(env *env.AppEnv, w http.ResponseWriter, r *http.Request, params WebhookPayload) error {
	data, err := parseLogChunkData(params.Data)
	if err != nil {
		return errors.WithStack(err)
	}

	err = env.WorkerService.EnqueueStoreLogChunkToRedis(params.TaskID.String(), models.LogChunk{
		TaskID:  params.TaskID,
		Pos:     data.Position,
		Content: data.Chunk,
	}, 5)
	if err != nil {
		return errors.Wrap(err, "Worker error")
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
