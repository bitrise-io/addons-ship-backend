package services

import (
	"fmt"
	"net/http"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/addons-ship-backend/worker"
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/satori/go.uuid"
)

// RootHandler ...
func RootHandler(env *env.AppEnv, w http.ResponseWriter, r *http.Request) error {
	eventID := uuid.NewV4()
	logChunk1 := models.LogChunk{Content: "Something\n"}
	logChunk2 := models.LogChunk{Content: "Another thing\n"}
	env.LogStoreService.Set(eventID.String()+"1", logChunk1)
	fmt.Println(env.LogStoreService.Get(eventID.String() + "1"))
	env.LogStoreService.Set(eventID.String()+"2", logChunk2)
	fmt.Println(env.LogStoreService.Get(eventID.String() + "2"))
	worker.EnqueueStoreLogToAWS(eventID, 2, "test_path/somethings.log")
	return httpresponse.RespondWithSuccess(w, map[string]string{"message": "Welcome to Bitrise Ship Addon!"})
}
