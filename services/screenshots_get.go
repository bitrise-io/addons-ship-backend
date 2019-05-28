package services

import (
	"net/http"
	"time"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/bitrise-io/api-utils/providers"
	"github.com/pkg/errors"
)

const (
	presignedURLExpirationInterval = 10 * time.Minute
)

// ScreenshotData ...
type ScreenshotData struct {
	models.Screenshot
	DownloadURL string `json:"download_url,omitempty"`
	UploadURL   string `json:"upload_url,omitempty"`
}

// ScreenshotsGetResponse ...
type ScreenshotsGetResponse struct {
	Data []ScreenshotData `json:"data"`
}

// ScreenshotsGetHandler ...
func ScreenshotsGetHandler(env *env.AppEnv, w http.ResponseWriter, r *http.Request) error {
	authorizedAppVersionID, err := GetAuthorizedAppVersionIDFromContext(r.Context())
	if err != nil {
		return errors.WithStack(err)
	}
	if env.ScreenshotService == nil {
		return errors.New("No Screenshot Service defined for handler")
	}

	screenshots, err := env.ScreenshotService.FindAll(
		&models.AppVersion{Record: models.Record{ID: authorizedAppVersionID}},
	)
	if err != nil {
		return errors.Wrap(err, "SQL Error")
	}

	if env.AWS == nil {
		return errors.New("No AWS Provider defined for handler")
	}
	responseData, err := newScreenshotGetResponseData(screenshots, env.AWS)
	if err != nil {
		return errors.WithStack(err)
	}

	return httpresponse.RespondWithSuccess(w, ScreenshotsGetResponse{
		Data: responseData,
	})
}

func newScreenshotGetResponseData(screenshots []models.Screenshot, awsProvider providers.AWSInterface) ([]ScreenshotData, error) {
	data := []ScreenshotData{}
	for _, screenshot := range screenshots {
		presignedURL, err := awsProvider.GeneratePresignedGETURL(screenshot.AWSPath(), presignedURLExpirationInterval)
		if err != nil {
			return []ScreenshotData{}, errors.WithStack(err)
		}
		data = append(data, ScreenshotData{Screenshot: screenshot, DownloadURL: presignedURL})
	}
	return data, nil
}
