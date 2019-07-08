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

func webhookPostStatusHelper(env *env.AppEnv, w http.ResponseWriter, r *http.Request, params WebhookPayload, appVersion *models.AppVersion) error {
	data, err := parseStatusData(params.Data)
	if err != nil {
		return httpresponse.RespondWithBadRequestError(w, "Invalid format of status data")
	}
	switch data.NewStatus {
	case "started":
		_, err := env.AppVersionEventService.Create(&models.AppVersionEvent{
			Status:       "in_progress",
			Text:         "Publishing has started",
			AppVersionID: appVersion.ID,
		})
		if err != nil {
			return errors.Wrap(err, "SQL Error")
		}
		err = env.Redis.Set(fmt.Sprintf("%s_chunk_count", params.TaskID.String()), 0, env.RedisExpirationTime)
		if err != nil {
			return errors.WithStack(err)
		}
		return httpresponse.RespondWithSuccess(w, httpresponse.StandardErrorRespModel{Message: "ok"})
	case "finished":
		var eventText, eventStatus string
		if data.ExitCode != 0 {
			eventStatus = "failed"
			eventText = "Failed to publish"
		} else {
			eventStatus = "success"
			eventText = "Successfully published"
		}
		event, err := env.AppVersionEventService.Create(&models.AppVersionEvent{
			Status:       eventStatus,
			Text:         eventText,
			AppVersionID: appVersion.ID,
		})
		if err != nil {
			return errors.Wrap(err, "SQL Error")
		}
		logAWSPath, err := event.LogAWSPath()
		if err != nil {
			return errors.WithStack(err)
		}
		err = env.WorkerService.EnqueueStoreLogToAWS(params.TaskID, data.LogChunkCount, logAWSPath, 30)
		if err != nil {
			return errors.Wrap(err, "Worker error")
		}
		return httpresponse.RespondWithSuccess(w, httpresponse.StandardErrorRespModel{Message: "ok"})
	default:
		return errors.Errorf("Invalid status of incoming webhook: %s", data.NewStatus)
	}
}

func parseStatusData(data interface{}) (StatusData, error) {
	var statusData StatusData
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return StatusData{}, err
	}
	err = json.Unmarshal(dataBytes, &statusData)
	if err != nil {
		return StatusData{}, err
	}
	return statusData, nil
}