package services

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/api-utils/httprequest"
	"github.com/bitrise-io/api-utils/httpresponse"
)

// StatusData ...
type StatusData struct {
	NewStatus     string    `json:"new_status"`
	ExitCode      int       `json:"exit_code"`
	LogChunkCount int       `json:"generated_log_chunk_count"`
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
	TaskID    string      `json:"task_id"`
	Data      interface{} `json:"data"`
}

// WebhookPostHandler ...
func WebhookPostHandler(env *env.AppEnv, w http.ResponseWriter, r *http.Request) error {
	var params WebhookPayload
	defer httprequest.BodyCloseWithErrorLog(r)
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return httpresponse.RespondWithBadRequestError(w, "Invalid request body, JSON decode failed")
	}

	switch params.TypeID {
	case "log":
		data, ok := params.Data.(LogChunkData)
		if !ok {
			return httpresponse.RespondWithBadRequestError(w, "Invalid format of log type webhook data")
		}
	case "status":
		data, ok := params.Data.(StatusData)
		if !ok {
			return httpresponse.RespondWithBadRequestError(w, "Invalid format of status type webhook data")
		}
		switch data.NewStatus {
		case "started":
		case "finidhed":

		}
	}
}
