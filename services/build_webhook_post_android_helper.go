package services

import (
	"encoding/json"
	"net/http"

	"github.com/bitrise-io/addons-ship-backend/bitrise"
	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/pkg/errors"
)

func prepareAppVersionForAndroidPlatform(env *env.AppEnv, w http.ResponseWriter,
	r *http.Request, apiToken, appSlug, buildSlug string) (*models.AppVersion, error) {
	artifacts, err := env.BitriseAPI.GetArtifacts(apiToken, appSlug, buildSlug)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	publishEnabled := false
	publicInstallPageEnabled := false
	publicInstallPageArtifactSlug := ""
	var selectedArtifact *bitrise.ArtifactListElementResponseModel

	for _, artifact := range artifacts {
		if artifact.IsAAB() {
			publishEnabled = true
			selectedArtifact = &artifact
		}
		if artifact.IsUniversalAPK() {
			publicInstallPageEnabled = true
			publicInstallPageArtifactSlug = artifact.Slug
			selectedArtifact = &artifact
		}
		// TODO: check the split APK condition
	}
	var artifactPublicInstallPageURL string
	if publicInstallPageEnabled {
		var err error
		artifactPublicInstallPageURL, err = env.BitriseAPI.GetArtifactPublicInstallPageURL(
			apiToken, appSlug, buildSlug, publicInstallPageArtifactSlug,
		)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	if selectedArtifact == nil {
		return nil, httpresponse.RespondWithNotFoundError(w)
	}

	appDetails, err := env.BitriseAPI.GetAppDetails(apiToken, appSlug)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	appInfo := models.AppInfo{
		MinimumSDK:  selectedArtifact.ArtifactMeta.AppInfo.MinimumSDKVersion,
		PackageName: selectedArtifact.ArtifactMeta.AppInfo.PackageName,
	}
	appInfoData, err := json.Marshal(appInfo)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &models.AppVersion{
		Platform:    "android",
		Version:     selectedArtifact.ArtifactMeta.AppInfo.VersionName,
		BuildSlug:   buildSlug,
		AppInfoData: appInfoData,
	}, nil
}
