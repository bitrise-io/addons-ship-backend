package services

import (
	"fmt"
	"net/http"
	"os"

	rice "github.com/GeertJohan/go.rice"
	"github.com/bitrise-io/addons-ship-backend/bitrise"
	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/bitrise-io/api-utils/structs"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
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

	if env.BitriseAPI == nil {
		return errors.New("No Bitrise API Service defined for handler")
	}

	config, err := getConfigJSON()
	if err != nil {
		return errors.WithStack(err)
	}

	artifactList, err := env.BitriseAPI.GetArtifacts(
		appVersion.App.BitriseAPIToken,
		appVersion.App.AppSlug,
		appVersion.BuildSlug,
	)
	if err != nil {
		return errors.WithStack(err)
	}

	type taskSecret map[string]string

	var workflowToTrigger, stackIDForTrigger string
	var inlineEnvs map[string]string
	var secrets map[string]interface{}
	switch appVersion.Platform {
	case "ios":
		artifactData, _, _, _ := selectIosArtifact(artifactList)
		workflowToTrigger = "resign_archive_app_store"
		stackIDForTrigger = "osx-vs4mac-stable"
		inlineEnvs = map[string]string{
			"BITRISE_APP_SLUG":      appVersion.App.AppSlug,
			"BITRISE_BUILD_SLUG":    appVersion.BuildSlug,
			"BITRISE_ARTIFACT_SLUG": artifactData.Slug,
			"CONFIG_JSON_URL":       fmt.Sprintf("%s/apps/%s/versions/%s/ios-config", env.AddonHostURL, appVersion.App.AppSlug, authorizedAppVersionID),
		}
		// {
		// 	key_str.to_s => value_str.to_s,
		// 	'opts' => {
		// 	  'is_expand' => is_expand == true
		// 	}
		//   }
		secrets = map[string]interface{}{"envs": []taskSecret{taskSecret{"BITRISE_ACCESS_TOKEN": appVersion.App.BitriseAPIToken}, taskSecret{"SHIP_ADDON_ACCESS_TOKEN": appVersion.App.APIToken}}}
	case "android":
		workflowToTrigger = "resign_android"
		stackIDForTrigger = "osx-vs4mac-stable"
		cloneUser := os.Getenv("ANDROID_PUBLISH_WF_GIT_CLONE_USER")
		clonePwd := os.Getenv("ANDROID_PUBLISH_WF_GIT_CLONE_PWD")
		inlineEnvs = map[string]string{
			"CONFIG_JSON_URL":    fmt.Sprintf("%s/apps/%s/versions/%s/android-config", env.AddonHostURL, appVersion.App.AppSlug, authorizedAppVersionID),
			"GIT_REPOSITORY_URL": fmt.Sprintf("https://%s:%s@github.com/bitrise-io/addons-ship-bg-worker-task-android", cloneUser, clonePwd),
		}
		secrets = map[string]interface{}{"envs": []taskSecret{
			taskSecret{"ADDON_SHIP_ACCESS_TOKEN": env.AddonAccessToken},
			taskSecret{"SHIP_ADDON_ACCESS_TOKEN": appVersion.App.APIToken},
			taskSecret{"BITRISE_ACCESS_TOKEN": appVersion.App.BitriseAPIToken},
		}}
	}

	if env.PublishTaskService == nil {
		return errors.New("No Publish Task Service defined for handler")
	}
	response, err := env.BitriseAPI.TriggerDENTask(bitrise.TaskParams{
		StackID:     stackIDForTrigger,
		Workflow:    workflowToTrigger,
		BuildConfig: config,
		InlineEnvs:  inlineEnvs,
		Secrets:     secrets,
		WebhookURL:  env.AddonHostURL + "/task-webhook",
	})
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = env.PublishTaskService.Create(&models.PublishTask{
		TaskID:       response.TaskIdentifier,
		AppVersionID: authorizedAppVersionID,
	})
	if err != nil {
		return errors.Wrap(err, "SQL Error")
	}

	return httpresponse.RespondWithSuccess(w, AppVersionPublishResponse{
		Data: response,
	})
}

func getConfigJSON() (interface{}, error) {
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
	return structs.ConvertMapIToMapS(config), nil
}
