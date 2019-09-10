package bitrise

import (
	"encoding/json"
	"time"

	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/pkg/errors"
)

// ArtifactSelector ...
type ArtifactSelector struct {
	artifacts []ArtifactListElementResponseModel
}

// NewArtifactSelector ...
func NewArtifactSelector(artifacts []ArtifactListElementResponseModel) ArtifactSelector {
	return ArtifactSelector{
		artifacts: artifacts,
	}
}

// PrepareAndroidAppVersions ...
func (s *ArtifactSelector) PrepareAndroidAppVersions(buildSlug, buildNumber, buildCommitMessage string) ([]models.AppVersion, error) {
	appVersions := []models.AppVersion{}
	flavorGroups := map[string][]ArtifactListElementResponseModel{}
	for _, artifact := range s.artifacts {
		if artifact.ArtifactMeta != nil {
			flavorGroups[artifact.ArtifactMeta.ProductFlavour] = append(flavorGroups[artifact.ArtifactMeta.ProductFlavour], artifact)
		}
	}
	for _, group := range flavorGroups {
		var buildType string
		buildTypeGroups := groupByBuildType(group)
		if len(buildTypeGroups) == 1 {
			buildType = group[0].ArtifactMeta.BuildType
		}
		artifactInfo := models.ArtifactInfo{
			MinimumSDK:     group[0].ArtifactMeta.AppInfo.MinimumSDKVersion,
			PackageName:    group[0].ArtifactMeta.AppInfo.PackageName,
			Version:        group[0].ArtifactMeta.AppInfo.VersionName,
			VersionCode:    group[0].ArtifactMeta.AppInfo.VersionCode,
			Module:         group[0].ArtifactMeta.Module,
			ProductFlavour: group[0].ArtifactMeta.ProductFlavour,
			BuildType:      buildType,
		}
		artifactInfoData, err := json.Marshal(artifactInfo)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		appVersion := models.AppVersion{
			Platform:         "android",
			BuildSlug:        buildSlug,
			BuildNumber:      buildNumber,
			ArtifactInfoData: artifactInfoData,
			LastUpdate:       time.Now(),
			CommitMessage:    buildCommitMessage,
		}
		appVersions = append(appVersions, appVersion)
	}
	return appVersions, nil
}

func groupByBuildType(artifacts []ArtifactListElementResponseModel) map[string][]ArtifactListElementResponseModel {
	buildTypeGroups := map[string][]ArtifactListElementResponseModel{}
	for _, artifact := range artifacts {
		buildTypeGroups[artifact.ArtifactMeta.BuildType] = append(buildTypeGroups[artifact.ArtifactMeta.BuildType], artifact)
	}
	return buildTypeGroups
}
