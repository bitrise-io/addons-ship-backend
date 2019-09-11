package bitrise_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/bitrise-io/addons-ship-backend/bitrise"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/stretchr/testify/require"
)

func compareAppVersions(t *testing.T, expected, actual models.AppVersion) {
	expected.LastUpdate = time.Time{}
	actual.LastUpdate = time.Time{}

	require.Equal(t, expected, actual)
}

func compareAppVersionArrays(t *testing.T, expecteds, actuals []models.AppVersion) {
	require.Len(t, actuals, len(expecteds))
	for i, expected := range expecteds {
		compareAppVersions(t, expected, actuals[i])
	}
}

func Test_ArtifactSelector_PrepareAndroidAppVersions(t *testing.T) {
	testBuildSlug := "test-build-slug"
	testBuildNumber := "test-build-number"
	testCommitMessage := "Some meaningful string"

	t.Run("ok", func(t *testing.T) {
		testArtifacts := []bitrise.ArtifactListElementResponseModel{
			bitrise.ArtifactListElementResponseModel{
				ArtifactMeta: &bitrise.ArtifactMeta{
					ProductFlavour: "salty",
				},
			},
			bitrise.ArtifactListElementResponseModel{
				ArtifactMeta: &bitrise.ArtifactMeta{
					ProductFlavour: "salty",
				},
			},
			bitrise.ArtifactListElementResponseModel{
				ArtifactMeta: &bitrise.ArtifactMeta{
					ProductFlavour: "sweet",
				},
			},
			bitrise.ArtifactListElementResponseModel{
				ArtifactMeta: &bitrise.ArtifactMeta{
					ProductFlavour: "sweet",
				},
			},
		}
		artifactSelector := bitrise.NewArtifactSelector(testArtifacts)

		expectedArtifactInfo := `{"version":"","version_code":"","minimum_os":"","minimum_sdk":"","size":0,"bundle_id":"","supported_device_types":null,"package_name":"","expire_date":"0001-01-01T00:00:00Z","ipa_export_method":"","module":"","build_type":""}`
		appVersions, settingsErr, err := artifactSelector.PrepareAndroidAppVersions(testBuildSlug, testBuildNumber, testCommitMessage, "")
		require.NoError(t, err)
		require.NoError(t, settingsErr)
		compareAppVersionArrays(t, []models.AppVersion{
			models.AppVersion{
				Platform:         "android",
				BuildSlug:        testBuildSlug,
				BuildNumber:      testBuildNumber,
				CommitMessage:    testCommitMessage,
				ArtifactInfoData: json.RawMessage(expectedArtifactInfo),
				ProductFlavour:   "salty",
			},
			models.AppVersion{
				Platform:         "android",
				BuildSlug:        testBuildSlug,
				BuildNumber:      testBuildNumber,
				CommitMessage:    testCommitMessage,
				ArtifactInfoData: json.RawMessage(expectedArtifactInfo),
				ProductFlavour:   "sweet",
			},
		}, appVersions)
	})

	t.Run("ok - multiple build type", func(t *testing.T) {
		testArtifacts := []bitrise.ArtifactListElementResponseModel{
			bitrise.ArtifactListElementResponseModel{
				ArtifactMeta: &bitrise.ArtifactMeta{
					ProductFlavour: "salty",
					BuildType:      "release",
				},
			},
			bitrise.ArtifactListElementResponseModel{
				ArtifactMeta: &bitrise.ArtifactMeta{
					ProductFlavour: "salty",
					BuildType:      "release",
				},
			},
			bitrise.ArtifactListElementResponseModel{
				ArtifactMeta: &bitrise.ArtifactMeta{
					ProductFlavour: "sweet",
					BuildType:      "release",
				},
			},
			bitrise.ArtifactListElementResponseModel{
				ArtifactMeta: &bitrise.ArtifactMeta{
					ProductFlavour: "sweet",
					BuildType:      "debug",
				},
			},
		}
		artifactSelector := bitrise.NewArtifactSelector(testArtifacts)

		expectedArtifactInfo := `{"version":"","version_code":"","minimum_os":"","minimum_sdk":"","size":0,"bundle_id":"","supported_device_types":null,"package_name":"","expire_date":"0001-01-01T00:00:00Z","ipa_export_method":"","module":"","build_type":"%s"}`
		appVersions, settingsErr, err := artifactSelector.PrepareAndroidAppVersions(testBuildSlug, testBuildNumber, testCommitMessage, "")
		require.NoError(t, err)
		require.NoError(t, settingsErr)
		compareAppVersionArrays(t, []models.AppVersion{
			models.AppVersion{
				Platform:         "android",
				BuildSlug:        testBuildSlug,
				BuildNumber:      testBuildNumber,
				CommitMessage:    testCommitMessage,
				ArtifactInfoData: json.RawMessage(fmt.Sprintf(expectedArtifactInfo, "release")),
				ProductFlavour:   "salty",
			},
			models.AppVersion{
				Platform:         "android",
				BuildSlug:        testBuildSlug,
				BuildNumber:      testBuildNumber,
				CommitMessage:    testCommitMessage,
				ArtifactInfoData: json.RawMessage(fmt.Sprintf(expectedArtifactInfo, "debug, release")),
				ProductFlavour:   "sweet",
			},
		}, appVersions)
	})

	t.Run("ok - multiple module", func(t *testing.T) {
		testArtifacts := []bitrise.ArtifactListElementResponseModel{
			bitrise.ArtifactListElementResponseModel{
				ArtifactMeta: &bitrise.ArtifactMeta{
					ProductFlavour: "salty",
					Module:         "test-module-1",
				},
			},
			bitrise.ArtifactListElementResponseModel{
				ArtifactMeta: &bitrise.ArtifactMeta{
					ProductFlavour: "salty",
					Module:         "test-module-2",
				},
			},
			bitrise.ArtifactListElementResponseModel{
				ArtifactMeta: &bitrise.ArtifactMeta{
					ProductFlavour: "sweet",
					Module:         "test-module-1",
				},
			},
			bitrise.ArtifactListElementResponseModel{
				ArtifactMeta: &bitrise.ArtifactMeta{
					ProductFlavour: "sweet",
					Module:         "test-module-2",
				},
			},
		}
		artifactSelector := bitrise.NewArtifactSelector(testArtifacts)

		expectedArtifactInfo := `{"version":"","version_code":"","minimum_os":"","minimum_sdk":"","size":0,"bundle_id":"","supported_device_types":null,"package_name":"","expire_date":"0001-01-01T00:00:00Z","ipa_export_method":"","module":"test-module-1","build_type":""}`
		appVersions, settingsErr, err := artifactSelector.PrepareAndroidAppVersions(testBuildSlug, testBuildNumber, testCommitMessage, "test-module-1")
		require.NoError(t, err)
		require.NoError(t, settingsErr)
		compareAppVersionArrays(t, []models.AppVersion{
			models.AppVersion{
				Platform:         "android",
				BuildSlug:        testBuildSlug,
				BuildNumber:      testBuildNumber,
				CommitMessage:    testCommitMessage,
				ArtifactInfoData: json.RawMessage(expectedArtifactInfo),
				ProductFlavour:   "salty",
			},
			models.AppVersion{
				Platform:         "android",
				BuildSlug:        testBuildSlug,
				BuildNumber:      testBuildNumber,
				CommitMessage:    testCommitMessage,
				ArtifactInfoData: json.RawMessage(expectedArtifactInfo),
				ProductFlavour:   "sweet",
			},
		}, appVersions)
	})

	t.Run("error - multiple module - no module set in settings", func(t *testing.T) {
		testArtifacts := []bitrise.ArtifactListElementResponseModel{
			bitrise.ArtifactListElementResponseModel{
				ArtifactMeta: &bitrise.ArtifactMeta{
					ProductFlavour: "salty",
					Module:         "test-module-1",
				},
			},
			bitrise.ArtifactListElementResponseModel{
				ArtifactMeta: &bitrise.ArtifactMeta{
					ProductFlavour: "salty",
					Module:         "test-module-2",
				},
			},
			bitrise.ArtifactListElementResponseModel{
				ArtifactMeta: &bitrise.ArtifactMeta{
					ProductFlavour: "sweet",
					Module:         "test-module-1",
				},
			},
			bitrise.ArtifactListElementResponseModel{
				ArtifactMeta: &bitrise.ArtifactMeta{
					ProductFlavour: "sweet",
					Module:         "test-module-2",
				},
			},
		}
		artifactSelector := bitrise.NewArtifactSelector(testArtifacts)
		appVersions, settingsErr, err := artifactSelector.PrepareAndroidAppVersions(testBuildSlug, testBuildNumber, testCommitMessage, "")
		require.NoError(t, err)
		require.EqualError(t, settingsErr, "No module setting found")
		require.Nil(t, appVersions)
	})
}
