package services

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/bitrise-io/addons-ship-backend/bitrise"
	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/pkg/errors"
)

func prepareAppVersionForAndroidPlatform(env *env.AppEnv, w http.ResponseWriter,
	r *http.Request, apiToken, appSlug, buildSlug string) (*models.AppVersion, error) {
	artifacts, err := env.BitriseAPI.GetArtifacts(apiToken, appSlug, buildSlug)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	selectedArtifact, _, _, _ := selectAndroidArtifact(artifacts)
	if selectedArtifact == nil {
		return nil, errors.New("No artifact found")
	}

	if selectedArtifact.ArtifactMeta == nil {
		return nil, errors.New("No artifact meta data found for artifact")
	}

	if reflect.DeepEqual(selectedArtifact.ArtifactMeta.AppInfo, bitrise.AppInfo{}) {
		return nil, errors.New("No artifact app info found for artifact")
	}

	artifactInfo := models.ArtifactInfo{
		MinimumSDK:  selectedArtifact.ArtifactMeta.AppInfo.MinimumSDKVersion,
		PackageName: selectedArtifact.ArtifactMeta.AppInfo.PackageName,
		Version:     selectedArtifact.ArtifactMeta.AppInfo.VersionName,
	}
	artifactInfoData, err := json.Marshal(artifactInfo)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &models.AppVersion{
		Platform:         "android",
		BuildSlug:        buildSlug,
		ArtifactInfoData: artifactInfoData,
	}, nil
}

func selectAndroidArtifact(artifacts []bitrise.ArtifactListElementResponseModel) (*bitrise.ArtifactListElementResponseModel, bool, bool, string) {
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
	return selectedArtifact, publishEnabled, publicInstallPageEnabled, publicInstallPageArtifactSlug
}
