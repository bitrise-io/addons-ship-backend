package services

import (
	"net/http"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/bitrise-io/api-utils/providers"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// AppEventData ...
type AppEventData struct {
	models.AppEvent
	LogDownloadURL string `json:"log_download_url"`
}

// AppEventsGetResponse ...
type AppEventsGetResponse struct {
	Data []AppEventData `json:"data"`
}

// AppEventsGetHandler ...
func AppEventsGetHandler(env *env.AppEnv, w http.ResponseWriter, r *http.Request) error {
	authorizedAppID, err := GetAuthorizedAppIDFromContext(r.Context())
	if err != nil {
		return errors.WithStack(err)
	}

	if env.AppEventService == nil {
		return errors.New("No App Event Service defined for handler")
	}

	appEvents, err := env.AppEventService.FindAll(&models.App{Record: models.Record{ID: authorizedAppID}})
	switch {
	case errors.Cause(err) == gorm.ErrRecordNotFound:
		return httpresponse.RespondWithNotFoundError(w)
	case err != nil:
		return errors.Wrap(err, "SQL Error")
	}

	if env.AWS == nil {
		return errors.New("No AWS Provider defined for handler")
	}

	responseData, err := newAppEventsGetResponse(appEvents, env.AWS)
	if err != nil {
		return errors.WithStack(err)
	}

	return httpresponse.RespondWithSuccess(w, AppEventsGetResponse{
		Data: responseData,
	})
}

func newAppEventsGetResponse(appEvents []models.AppEvent, awsProvider providers.AWSInterface) ([]AppEventData, error) {
	data := []AppEventData{}
	for _, appEvent := range appEvents {
		awsPath, err := appEvent.LogAWSPath()
		if err != nil {
			return []AppEventData{}, errors.WithStack(err)
		}
		presignedURL, err := awsProvider.GeneratePresignedGETURL(awsPath, presignedURLExpirationInterval)
		if err != nil {
			return []AppEventData{}, errors.WithStack(err)
		}
		data = append(data, AppEventData{AppEvent: appEvent, LogDownloadURL: presignedURL})
	}
	return data, nil
}
