package services

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/bitrise-io/addons-ship-backend/bitrise"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/pkg/errors"
)

func prepareAppVersionForAndroidPlatform(w http.ResponseWriter, r *http.Request, artifacts []bitrise.ArtifactListElementResponseModel, buildSlug string) (*models.AppVersion, error) {
	selectedArtifact, _, _, _ := selectAndroidArtifact(artifacts)

	if selectedArtifact == nil || reflect.DeepEqual(*selectedArtifact, bitrise.ArtifactListElementResponseModel{}) {
		splitAPKs := checkForSplitAPKs(artifacts)
		if len(splitAPKs) == 0 {
			return nil, errors.New("No Android artifact found")
		}
		selectedArtifact = &splitAPKs[0]
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

func hasAndroidArtifact(artifacts []bitrise.ArtifactListElementResponseModel) bool {
	for _, artifact := range artifacts {
		if artifact.IsAAB() || artifact.IsUniversalAPK() {
			return true
		}
	}

	return false
}

func selectAndroidArtifact(artifacts []bitrise.ArtifactListElementResponseModel) (*bitrise.ArtifactListElementResponseModel, bool, bool, string) {
	publishEnabled := false
	publicInstallPageEnabled := false
	publicInstallPageArtifactSlug := ""
	var selectedArtifact bitrise.ArtifactListElementResponseModel
	for _, artifact := range artifacts {
		if artifact.IsAAB() {
			publishEnabled = true
			selectedArtifact = artifact
		}
		if artifact.IsUniversalAPK() {
			publicInstallPageEnabled = true
			publicInstallPageArtifactSlug = artifact.Slug
			selectedArtifact = artifact
		}
	}
	return &selectedArtifact, publishEnabled, publicInstallPageEnabled, publicInstallPageArtifactSlug
}

func checkForSplitAPKs(artifacts []bitrise.ArtifactListElementResponseModel) []bitrise.ArtifactListElementResponseModel {
	selectedArtifacts := map[string][]bitrise.ArtifactListElementResponseModel{}
	maxNumber := 0
	maxNumberKey := ""
	for _, artifact := range artifacts {
		if artifact.ArtifactMeta != nil {
			key := artifact.ArtifactMeta.AppInfo.AppName + artifact.ArtifactMeta.AppInfo.PackageName + artifact.ArtifactMeta.AppInfo.VersionName
			selectedArtifacts[key] = append(selectedArtifacts[key], artifact)
			if len(selectedArtifacts[key]) > maxNumber {
				maxNumber = len(selectedArtifacts[key])
				maxNumberKey = key
			}
		}
	}
	if maxNumber > 1 {
		return selectedArtifacts[maxNumberKey]
	}
	return []bitrise.ArtifactListElementResponseModel{}
}
