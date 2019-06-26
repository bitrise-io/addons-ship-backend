package services

import (
	"encoding/json"
	"net/http"

	rice "github.com/GeertJohan/go.rice"
	"github.com/bitrise-io/addons-ship-backend/bitrise"
	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/bitrise-io/api-utils/structs"
	"github.com/go-yaml/yaml"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// AppVersionPublishResponse ...
type AppVersionPublishResponse struct {
	Data *bitrise.TriggerResponse `json:"data"`
}

// AppVersionPublishPostHandler ...
func AppVersionPublishPostHandler(env *env.AppEnv, w http.ResponseWriter, r *http.Request) error {
	authorizedAppVersionID, err := GetAuthorizedAppVersionIDFromContext(r.Context())
	if err != nil {
		return errors.WithStack(err)
	}

	if env.AppVersionService == nil {
		return errors.New("No App Version Service defined for handler")
	}

	appVersion, err := env.AppVersionService.Find(
		&models.AppVersion{Record: models.Record{ID: authorizedAppVersionID}},
	)
	switch {
	case errors.Cause(err) == gorm.ErrRecordNotFound:
		return httpresponse.RespondWithNotFoundError(w)
	case err != nil:
		return errors.Wrap(err, "SQL Error")
	}
	var workflowToTrigger string
	if appVersion.Platform == "ios" {
		workflowToTrigger = "resign_archive_app_store"
	}

	if env.BitriseAPI == nil {
		return errors.New("No Bitrise API Service defined for handler")
	}

	config, err := getConfigJSON()
	if err != nil {
		return errors.WithStack(err)
	}

	artifactData, err := env.BitriseAPI.GetArtifactData(
		appVersion.App.BitriseAPIToken,
		appVersion.App.AppSlug,
		appVersion.BuildSlug,
	)
	if err != nil {
		return errors.WithStack(err)
	}

	inlineEnvs := map[string]string{
		"BITRISE_APP_SLUG":      appVersion.App.AppSlug,
		"BITRISE_BUILD_SLUG":    appVersion.BuildSlug,
		"BITRISE_ARTIFACT_SLUG": artifactData.Slug,
	}
	inlineEnvsBytes, err := json.Marshal(inlineEnvs)
	if err != nil {
		return errors.WithStack(err)
	}
	secrets := map[string]string{"BITRISE_ACCESS_TOKEN": appVersion.App.BitriseAPIToken}
	secretsBytes, err := json.Marshal(secrets)
	if err != nil {
		return errors.WithStack(err)
	}

	response, err := env.BitriseAPI.TriggerDENTask(bitrise.TaskParams{
		Workflow:    workflowToTrigger,
		BuildConfig: config,
		InlineEnvs:  string(inlineEnvsBytes),
		Secrets:     string(secretsBytes),
		WebhookURL:  env.AddonHostURL + "/webhook",
	})
	if err != nil {
		return errors.WithStack(err)
	}
	return httpresponse.RespondWithSuccess(w, AppVersionPublishResponse{
		Data: response,
	})
}

func getConfigJSON() (string, error) {
	templateBox, err := rice.FindBox("../utility")
	if err != nil {
		return "", errors.WithStack(err)
	}
	tmpContent, err := templateBox.String("workflows.yml")
	if err != nil {
		return "", errors.WithStack(err)
	}

	var config interface{}
	err = yaml.Unmarshal([]byte(tmpContent), &config)
	if err != nil {
		return "", err
	}
	config = structs.ConvertMapIToMapS(config)

	jsonBytes, err := json.Marshal(config)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}
