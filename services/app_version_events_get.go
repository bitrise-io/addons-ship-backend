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

// AppVersionEventData ...
type AppVersionEventData struct {
	models.AppVersionEvent
	LogDownloadURL string `json:"log_download_url"`
}

// AppVersionEventsGetResponse ...
type AppVersionEventsGetResponse struct {
	Data []AppVersionEventData `json:"data"`
}

// AppVersionEventsGetHandler ...
func AppVersionEventsGetHandler(env *env.AppEnv, w http.ResponseWriter, r *http.Request) error {
	authorizedAppVersionID, err := GetAuthorizedAppVersionIDFromContext(r.Context())
	if err != nil {
		return errors.WithStack(err)
	}

	if env.AppVersionEventService == nil {
		return errors.New("No App Version Event Service defined for handler")
	}

	appVersionEvents, err := env.AppVersionEventService.FindAll(&models.AppVersion{Record: models.Record{ID: authorizedAppVersionID}})
	switch {
	case errors.Cause(err) == gorm.ErrRecordNotFound:
		return httpresponse.RespondWithNotFoundError(w)
	case err != nil:
		return errors.Wrap(err, "SQL Error")
	}

	if env.AWS == nil {
		return errors.New("No AWS Provider defined for handler")
	}

	responseData, err := newAppVersionEventsGetResponse(appVersionEvents, env.AWS)
	if err != nil {
		return errors.WithStack(err)
	}

	return httpresponse.RespondWithSuccess(w, AppVersionEventsGetResponse{
		Data: responseData,
	})
}

func newAppVersionEventsGetResponse(appVersionEvents []models.AppVersionEvent, awsProvider providers.AWSInterface) ([]AppVersionEventData, error) {
	data := []AppVersionEventData{}
	for _, appVersionEvent := range appVersionEvents {
		awsPath, err := appVersionEvent.LogAWSPath()
		if err != nil {
			return []AppVersionEventData{}, errors.WithStack(err)
		}
		presignedURL, err := awsProvider.GeneratePresignedGETURL(awsPath, presignedURLExpirationInterval)
		if err != nil {
			return []AppVersionEventData{}, errors.WithStack(err)
		}
		data = append(data, AppVersionEventData{AppVersionEvent: appVersionEvent, LogDownloadURL: presignedURL})
	}
	return data, nil
}
