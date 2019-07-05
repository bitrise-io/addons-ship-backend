package services

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/api-utils/httprequest"
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

// StatusData ...
type StatusData struct {
	NewStatus     string    `json:"new_status"`
	ExitCode      int       `json:"exit_code"`
	LogChunkCount int64     `json:"generated_log_chunk_count"`
	FinishedAt    time.Time `json:"finished_at"`
}

// LogChunkData ...
type LogChunkData struct {
	Position int    `json:"position"`
	Chunk    string `json:"chunk"`
}

// WebhookPayload ...
type WebhookPayload struct {
	TypeID    string      `json:"type_id"`
	Timestamp int64       `json:"timestamp"`
	TaskID    uuid.UUID   `json:"task_id"`
	Data      interface{} `json:"data"`
}

// WebhookPostHandler ...
func WebhookPostHandler(env *env.AppEnv, w http.ResponseWriter, r *http.Request) error {
	authorizedAppVersionID, err := GetAuthorizedAppVersionIDFromContext(r.Context())
	if err != nil {
		return errors.WithStack(err)
	}

	if env.AppVersionService == nil {
		return errors.New("No App Version Service provided")
	}
	if env.AppVersionEventService == nil {
		return errors.New("No App Version Event Service provided")
	}

	var params WebhookPayload
	defer httprequest.BodyCloseWithErrorLog(r)
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return httpresponse.RespondWithBadRequestError(w, "Invalid request body, JSON decode failed")
	}

	appVersion, err := env.AppVersionService.Find(&models.AppVersion{Record: models.Record{ID: authorizedAppVersionID}})
	if err != nil {
		return errors.WithStack(err)
	}
	switch params.TypeID {
	case "log":
		return WebhookPostLogHelper(env, w, r, params)
	case "status":
		return WebhookPostStatusHelper(env, w, r, params, appVersion)
	default:
		return errors.Errorf("Invalid type of webhook: %s", params.TypeID)
	}
}
