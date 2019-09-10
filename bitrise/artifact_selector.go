package bitrise

import (
	"encoding/json"
	"strings"
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
		} else {
			return nil, errors.New("No artifact meta data found for artifact")
		}
	}
	for _, group := range flavorGroups {
		buildTypeGroups := groupByBuildType(group)
		keys := []string{}
		for key := range buildTypeGroups {
			keys = append(keys, key)
		}
		artifactInfo := models.ArtifactInfo{
			MinimumSDK:  group[0].ArtifactMeta.AppInfo.MinimumSDKVersion,
			PackageName: group[0].ArtifactMeta.AppInfo.PackageName,
			Version:     group[0].ArtifactMeta.AppInfo.VersionName,
			VersionCode: group[0].ArtifactMeta.AppInfo.VersionCode,
			Module:      group[0].ArtifactMeta.Module,
			BuildType:   strings.Join(keys, ", "),
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
			ProductFlavour:   group[0].ArtifactMeta.ProductFlavour,
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
