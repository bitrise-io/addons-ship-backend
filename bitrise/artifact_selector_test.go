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
					ProductFlavor: "salty",
				},
			},
			bitrise.ArtifactListElementResponseModel{
				ArtifactMeta: &bitrise.ArtifactMeta{
					ProductFlavor: "salty",
				},
			},
			bitrise.ArtifactListElementResponseModel{
				ArtifactMeta: &bitrise.ArtifactMeta{
					ProductFlavor: "sweet",
				},
			},
			bitrise.ArtifactListElementResponseModel{
				ArtifactMeta: &bitrise.ArtifactMeta{
					ProductFlavor: "sweet",
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
				ProductFlavor:    "salty",
			},
			models.AppVersion{
				Platform:         "android",
				BuildSlug:        testBuildSlug,
				BuildNumber:      testBuildNumber,
				CommitMessage:    testCommitMessage,
				ArtifactInfoData: json.RawMessage(expectedArtifactInfo),
				ProductFlavor:    "sweet",
			},
		}, appVersions)
	})

	t.Run("ok - multiple build type", func(t *testing.T) {
		testArtifacts := []bitrise.ArtifactListElementResponseModel{
			bitrise.ArtifactListElementResponseModel{
				ArtifactMeta: &bitrise.ArtifactMeta{
					ProductFlavor: "salty",
					BuildType:     "release",
				},
			},
			bitrise.ArtifactListElementResponseModel{
				ArtifactMeta: &bitrise.ArtifactMeta{
					ProductFlavor: "salty",
					BuildType:     "release",
				},
			},
			bitrise.ArtifactListElementResponseModel{
				ArtifactMeta: &bitrise.ArtifactMeta{
					ProductFlavor: "sweet",
					BuildType:     "release",
				},
			},
			bitrise.ArtifactListElementResponseModel{
				ArtifactMeta: &bitrise.ArtifactMeta{
					ProductFlavor: "sweet",
					BuildType:     "debug",
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
				ProductFlavor:    "salty",
			},
			models.AppVersion{
				Platform:         "android",
				BuildSlug:        testBuildSlug,
				BuildNumber:      testBuildNumber,
				CommitMessage:    testCommitMessage,
				ArtifactInfoData: json.RawMessage(fmt.Sprintf(expectedArtifactInfo, "debug, release")),
				ProductFlavor:    "sweet",
			},
		}, appVersions)
	})

	t.Run("ok - multiple module", func(t *testing.T) {
		testArtifacts := []bitrise.ArtifactListElementResponseModel{
			bitrise.ArtifactListElementResponseModel{
				ArtifactMeta: &bitrise.ArtifactMeta{
					ProductFlavor: "salty",
					Module:        "test-module-1",
				},
			},
			bitrise.ArtifactListElementResponseModel{
				ArtifactMeta: &bitrise.ArtifactMeta{
					ProductFlavor: "salty",
					Module:        "test-module-2",
				},
			},
			bitrise.ArtifactListElementResponseModel{
				ArtifactMeta: &bitrise.ArtifactMeta{
					ProductFlavor: "sweet",
					Module:        "test-module-1",
				},
			},
			bitrise.ArtifactListElementResponseModel{
				ArtifactMeta: &bitrise.ArtifactMeta{
					ProductFlavor: "sweet",
					Module:        "test-module-2",
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
				ProductFlavor:    "salty",
			},
			models.AppVersion{
				Platform:         "android",
				BuildSlug:        testBuildSlug,
				BuildNumber:      testBuildNumber,
				CommitMessage:    testCommitMessage,
				ArtifactInfoData: json.RawMessage(expectedArtifactInfo),
				ProductFlavor:    "sweet",
			},
		}, appVersions)
	})

	t.Run("error - multiple module - no module set in settings", func(t *testing.T) {
		testArtifacts := []bitrise.ArtifactListElementResponseModel{
			bitrise.ArtifactListElementResponseModel{
				ArtifactMeta: &bitrise.ArtifactMeta{
					ProductFlavor: "salty",
					Module:        "test-module-1",
				},
			},
			bitrise.ArtifactListElementResponseModel{
				ArtifactMeta: &bitrise.ArtifactMeta{
					ProductFlavor: "salty",
					Module:        "test-module-2",
				},
			},
			bitrise.ArtifactListElementResponseModel{
				ArtifactMeta: &bitrise.ArtifactMeta{
					ProductFlavor: "sweet",
					Module:        "test-module-1",
				},
			},
			bitrise.ArtifactListElementResponseModel{
				ArtifactMeta: &bitrise.ArtifactMeta{
					ProductFlavor: "sweet",
					Module:        "test-module-2",
				},
			},
		}
		artifactSelector := bitrise.NewArtifactSelector(testArtifacts)
		appVersions, settingsErr, err := artifactSelector.PrepareAndroidAppVersions(testBuildSlug, testBuildNumber, testCommitMessage, "")
		require.NoError(t, err)
		require.EqualError(t, settingsErr, "Multiple modules found, but none was selected in the settings")
		require.Nil(t, appVersions)
	})
}

func Test_ArtifactSelector_Select(t *testing.T) {
	for _, tc := range []struct {
		testName            string
		artifacts           []bitrise.ArtifactListElementResponseModel
		moduleName          string
		flavor              string
		expectedSlugs       []string
		expectedSettingsErr string
	}{
		{
			testName:      "ok - minimal",
			artifacts:     []bitrise.ArtifactListElementResponseModel{},
			expectedSlugs: []string{},
		},
		{
			testName: "ok - release build type - standalone apk",
			artifacts: []bitrise.ArtifactListElementResponseModel{
				bitrise.ArtifactListElementResponseModel{
					Title: "my-awesome.apk",
					Slug:  "test-apk-1",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Apk:           "/bitrise/my-project/my-awesome.apk",
						ProductFlavor: "",
						Module:        "",
						BuildType:     "release",
					},
				},
			},
			expectedSlugs: []string{"test-apk-1"},
		},
		{
			testName: "ok - debug build type - standalone apk",
			artifacts: []bitrise.ArtifactListElementResponseModel{
				bitrise.ArtifactListElementResponseModel{
					Title: "my-awesome.apk",
					Slug:  "test-apk-1",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Apk:           "/bitrise/my-project/my-awesome.apk",
						ProductFlavor: "",
						Module:        "",
						BuildType:     "debug",
					},
				},
			},
			expectedSlugs: []string{},
		},
		{
			testName: "ok - release build type - split apk",
			artifacts: []bitrise.ArtifactListElementResponseModel{
				bitrise.ArtifactListElementResponseModel{
					Title: "app-hdpi.apk",
					Slug:  "test-apk-1",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Apk:           "",
						Split:         []string{"app-hdpi.apk", "app-mdpi.apk", "app-xhdpi.apk", "app-xxhdpi.apk"},
						ProductFlavor: "",
						Module:        "",
						BuildType:     "release",
					},
				},
				bitrise.ArtifactListElementResponseModel{
					Title: "app-mdpi.apk",
					Slug:  "test-apk-2",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Apk:           "",
						Split:         []string{"app-hdpi.apk", "app-mdpi.apk", "app-xhdpi.apk", "app-xxhdpi.apk"},
						ProductFlavor: "",
						Module:        "",
						BuildType:     "release",
					},
				},
				bitrise.ArtifactListElementResponseModel{
					Title: "app-xhdpi.apk",
					Slug:  "test-apk-3",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Apk:           "",
						Split:         []string{"app-hdpi.apk", "app-mdpi.apk", "app-xhdpi.apk", "app-xxhdpi.apk"},
						ProductFlavor: "",
						Module:        "",
						BuildType:     "release",
					},
				},
				bitrise.ArtifactListElementResponseModel{
					Title: "app-xxhdpi.apk",
					Slug:  "test-apk-4",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Apk:           "",
						Split:         []string{"app-hdpi.apk", "app-mdpi.apk", "app-xhdpi.apk", "app-xxhdpi.apk"},
						ProductFlavor: "",
						Module:        "",
						BuildType:     "release",
					},
				},
			},
			expectedSlugs: []string{"test-apk-1", "test-apk-2", "test-apk-3", "test-apk-4"},
		},
		{
			testName: "ok - release build type - split apk with universal apk",
			artifacts: []bitrise.ArtifactListElementResponseModel{
				bitrise.ArtifactListElementResponseModel{
					Title: "app-hdpi.apk",
					Slug:  "test-apk-1",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Apk:           "",
						Split:         []string{"app-hdpi.apk", "app-mdpi.apk", "app-xhdpi.apk", "app-xxhdpi.apk"},
						ProductFlavor: "",
						Module:        "",
						BuildType:     "release",
						Universal:     "/bitrise/my-project/app-hdpi-universal.apk",
					},
				},
				bitrise.ArtifactListElementResponseModel{
					Title: "app-mdpi.apk",
					Slug:  "test-apk-2",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Apk:           "",
						Split:         []string{"app-hdpi.apk", "app-mdpi.apk", "app-xhdpi.apk", "app-xxhdpi.apk"},
						ProductFlavor: "",
						Module:        "",
						BuildType:     "release",
						Universal:     "/bitrise/my-project/app-hdpi-universal.apk",
					},
				},
				bitrise.ArtifactListElementResponseModel{
					Title: "app-xhdpi.apk",
					Slug:  "test-apk-3",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Apk:           "",
						Split:         []string{"app-hdpi.apk", "app-mdpi.apk", "app-xhdpi.apk", "app-xxhdpi.apk"},
						ProductFlavor: "",
						Module:        "",
						BuildType:     "release",
						Universal:     "/bitrise/my-project/app-hdpi-universal.apk",
					},
				},
				bitrise.ArtifactListElementResponseModel{
					Title: "app-xxhdpi.apk",
					Slug:  "test-apk-4",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Apk:           "",
						Split:         []string{"app-hdpi.apk", "app-mdpi.apk", "app-xhdpi.apk", "app-xxhdpi.apk"},
						ProductFlavor: "",
						Module:        "",
						BuildType:     "release",
						Universal:     "/bitrise/my-project/app-hdpi-universal.apk",
					},
				},
				bitrise.ArtifactListElementResponseModel{
					Title: "app-hdpi-universal.apk",
					Slug:  "test-apk-5",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Apk:           "",
						Split:         []string{"app-hdpi.apk", "app-mdpi.apk", "app-xhdpi.apk", "app-xxhdpi.apk"},
						ProductFlavor: "",
						Module:        "",
						BuildType:     "release",
						Universal:     "/bitrise/my-project/app-hdpi-universal.apk",
					},
				},
			},
			expectedSlugs: []string{"test-apk-1", "test-apk-2", "test-apk-3", "test-apk-4", "test-apk-5"},
		},
		{
			testName: "ok - release build type - split apk with aab",
			artifacts: []bitrise.ArtifactListElementResponseModel{
				bitrise.ArtifactListElementResponseModel{
					Title: "app-hdpi.apk",
					Slug:  "test-apk-1",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Aab:           "/bitrise/my-project/app.aab",
						Apk:           "",
						Split:         []string{"app-hdpi.apk", "app-mdpi.apk", "app-xhdpi.apk", "app-xxhdpi.apk"},
						ProductFlavor: "",
						Module:        "",
						BuildType:     "release",
					},
				},
				bitrise.ArtifactListElementResponseModel{
					Title: "app-mdpi.apk",
					Slug:  "test-apk-2",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Aab:           "/bitrise/my-project/app.aab",
						Apk:           "",
						Split:         []string{"app-hdpi.apk", "app-mdpi.apk", "app-xhdpi.apk", "app-xxhdpi.apk"},
						ProductFlavor: "",
						Module:        "",
						BuildType:     "release",
					},
				},
				bitrise.ArtifactListElementResponseModel{
					Title: "app-xhdpi.apk",
					Slug:  "test-apk-3",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Aab:           "/bitrise/my-project/app.aab",
						Apk:           "",
						Split:         []string{"app-hdpi.apk", "app-mdpi.apk", "app-xhdpi.apk", "app-xxhdpi.apk"},
						ProductFlavor: "",
						Module:        "",
						BuildType:     "release",
					},
				},
				bitrise.ArtifactListElementResponseModel{
					Title: "app-xxhdpi.apk",
					Slug:  "test-apk-4",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Aab:           "/bitrise/my-project/app.aab",
						Apk:           "",
						Split:         []string{"app-hdpi.apk", "app-mdpi.apk", "app-xhdpi.apk", "app-xxhdpi.apk"},
						ProductFlavor: "",
						Module:        "",
						BuildType:     "release",
					},
				},
				bitrise.ArtifactListElementResponseModel{
					Title: "app.aab",
					Slug:  "test-apk-5",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Aab:           "/bitrise/my-project/app.aab",
						Apk:           "",
						Split:         []string{"app-hdpi.apk", "app-mdpi.apk", "app-xhdpi.apk", "app-xxhdpi.apk"},
						ProductFlavor: "",
						Module:        "",
						BuildType:     "release",
					},
				},
			},
			expectedSlugs: []string{"test-apk-5"},
		},
		{
			testName: "ok - release build type - multiple flavor - split apk",
			artifacts: []bitrise.ArtifactListElementResponseModel{
				bitrise.ArtifactListElementResponseModel{
					Title: "app-sweet-hdpi.apk",
					Slug:  "test-apk-1",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Apk:           "",
						Split:         []string{"app-sweet-hdpi.apk", "app-sweet-mdpi.apk", "app-sweet-xhdpi.apk", "app-sweet-xxhdpi.apk"},
						ProductFlavor: "sweet",
						Module:        "",
						BuildType:     "release",
					},
				},
				bitrise.ArtifactListElementResponseModel{
					Title: "app-sweet-mdpi.apk",
					Slug:  "test-apk-2",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Apk:           "",
						Split:         []string{"app-sweet-hdpi.apk", "app-sweet-mdpi.apk", "app-sweet-xhdpi.apk", "app-sweet-xxhdpi.apk"},
						ProductFlavor: "sweet",
						Module:        "",
						BuildType:     "release",
					},
				},
				bitrise.ArtifactListElementResponseModel{
					Title: "app-sweet-xhdpi.apk",
					Slug:  "test-apk-3",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Apk:           "",
						Split:         []string{"app-sweet-hdpi.apk", "app-sweet-mdpi.apk", "app-sweet-xhdpi.apk", "app-sweet-xxhdpi.apk"},
						ProductFlavor: "sweet",
						Module:        "",
						BuildType:     "release",
					},
				},
				bitrise.ArtifactListElementResponseModel{
					Title: "app-sweet-xxhdpi.apk",
					Slug:  "test-apk-4",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Apk:           "",
						Split:         []string{"app-sweet-hdpi.apk", "app-sweet-mdpi.apk", "app-sweet-xhdpi.apk", "app-sweet-xxhdpi.apk"},
						ProductFlavor: "sweet",
						Module:        "",
						BuildType:     "release",
					},
				},
				bitrise.ArtifactListElementResponseModel{
					Title: "app-salty-hdpi.apk",
					Slug:  "test-apk-5",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Apk:           "",
						Split:         []string{"app-salty-hdpi.apk", "app-salty-mdpi.apk", "app-salty-xhdpi.apk", "app-salty-xxhdpi.apk"},
						ProductFlavor: "salty",
						Module:        "",
						BuildType:     "release",
					},
				},
				bitrise.ArtifactListElementResponseModel{
					Title: "app-salty-mdpi.apk",
					Slug:  "test-apk-6",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Apk:           "",
						Split:         []string{"app-salty-hdpi.apk", "app-salty-mdpi.apk", "app-salty-xhdpi.apk", "app-salty-xxhdpi.apk"},
						ProductFlavor: "salty",
						Module:        "",
						BuildType:     "release",
					},
				},
				bitrise.ArtifactListElementResponseModel{
					Title: "app-salty-xhdpi.apk",
					Slug:  "test-apk-7",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Apk:           "",
						Split:         []string{"app-salty-hdpi.apk", "app-salty-mdpi.apk", "app-salty-xhdpi.apk", "app-salty-xxhdpi.apk"},
						ProductFlavor: "salty",
						Module:        "",
						BuildType:     "release",
					},
				},
				bitrise.ArtifactListElementResponseModel{
					Title: "app-salty-xxhdpi.apk",
					Slug:  "test-apk-8",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Apk:           "",
						Split:         []string{"app-salty-hdpi.apk", "app-salty-mdpi.apk", "app-salty-xhdpi.apk", "app-salty-xxhdpi.apk"},
						ProductFlavor: "salty",
						Module:        "",
						BuildType:     "release",
					},
				},
			},
			flavor:        "salty",
			expectedSlugs: []string{"test-apk-5", "test-apk-6", "test-apk-7", "test-apk-8"},
		},
		{
			testName: "ok - release build type - multiple flavor, multiple module - split apk",
			artifacts: []bitrise.ArtifactListElementResponseModel{
				bitrise.ArtifactListElementResponseModel{
					Title: "app-sweet-hdpi.apk",
					Slug:  "test-apk-1",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Apk:           "",
						Split:         []string{"app-sweet-hdpi.apk", "app-sweet-mdpi.apk", "app-sweet-xhdpi.apk", "app-sweet-xxhdpi.apk"},
						ProductFlavor: "sweet",
						Module:        "module-1",
						BuildType:     "release",
					},
				},
				bitrise.ArtifactListElementResponseModel{
					Title: "app-sweet-mdpi.apk",
					Slug:  "test-apk-2",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Apk:           "",
						Split:         []string{"app-sweet-hdpi.apk", "app-sweet-mdpi.apk", "app-sweet-xhdpi.apk", "app-sweet-xxhdpi.apk"},
						ProductFlavor: "sweet",
						Module:        "module-1",
						BuildType:     "release",
					},
				},
				bitrise.ArtifactListElementResponseModel{
					Title: "app-sweet-xhdpi.apk",
					Slug:  "test-apk-3",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Apk:           "",
						Split:         []string{"app-sweet-hdpi.apk", "app-sweet-mdpi.apk", "app-sweet-xhdpi.apk", "app-sweet-xxhdpi.apk"},
						ProductFlavor: "sweet",
						Module:        "module-1",
						BuildType:     "release",
					},
				},
				bitrise.ArtifactListElementResponseModel{
					Title: "app-sweet-xxhdpi.apk",
					Slug:  "test-apk-4",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Apk:           "",
						Split:         []string{"app-sweet-hdpi.apk", "app-sweet-mdpi.apk", "app-sweet-xhdpi.apk", "app-sweet-xxhdpi.apk"},
						ProductFlavor: "sweet",
						Module:        "module-1",
						BuildType:     "release",
					},
				},
				bitrise.ArtifactListElementResponseModel{
					Title: "app-salty-hdpi.apk",
					Slug:  "test-apk-5",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Apk:           "",
						Split:         []string{"app-salty-hdpi.apk", "app-salty-mdpi.apk", "app-salty-xhdpi.apk", "app-salty-xxhdpi.apk"},
						ProductFlavor: "salty",
						Module:        "module-1",
						BuildType:     "release",
					},
				},
				bitrise.ArtifactListElementResponseModel{
					Title: "app-salty-mdpi.apk",
					Slug:  "test-apk-6",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Apk:           "",
						Split:         []string{"app-salty-hdpi.apk", "app-salty-mdpi.apk", "app-salty-xhdpi.apk", "app-salty-xxhdpi.apk"},
						ProductFlavor: "salty",
						Module:        "module-1",
						BuildType:     "release",
					},
				},
				bitrise.ArtifactListElementResponseModel{
					Title: "app-salty-xhdpi.apk",
					Slug:  "test-apk-7",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Apk:           "",
						Split:         []string{"app-salty-hdpi.apk", "app-salty-mdpi.apk", "app-salty-xhdpi.apk", "app-salty-xxhdpi.apk"},
						ProductFlavor: "salty",
						Module:        "module-1",
						BuildType:     "release",
					},
				},
				bitrise.ArtifactListElementResponseModel{
					Title: "app-salty-xxhdpi.apk",
					Slug:  "test-apk-8",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Apk:           "",
						Split:         []string{"app-salty-hdpi.apk", "app-salty-mdpi.apk", "app-salty-xhdpi.apk", "app-salty-xxhdpi.apk"},
						ProductFlavor: "salty",
						Module:        "module-1",
						BuildType:     "release",
					},
				},
				bitrise.ArtifactListElementResponseModel{
					Title: "app-sweet-hdpi.apk",
					Slug:  "test-apk-9",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Apk:           "",
						Split:         []string{"app-sweet-hdpi.apk", "app-sweet-mdpi.apk", "app-sweet-xhdpi.apk", "app-sweet-xxhdpi.apk"},
						ProductFlavor: "sweet",
						Module:        "module-2",
						BuildType:     "release",
					},
				},
				bitrise.ArtifactListElementResponseModel{
					Title: "app-sweet-mdpi.apk",
					Slug:  "test-apk-10",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Apk:           "",
						Split:         []string{"app-sweet-hdpi.apk", "app-sweet-mdpi.apk", "app-sweet-xhdpi.apk", "app-sweet-xxhdpi.apk"},
						ProductFlavor: "sweet",
						Module:        "module-2",
						BuildType:     "release",
					},
				},
				bitrise.ArtifactListElementResponseModel{
					Title: "app-sweet-xhdpi.apk",
					Slug:  "test-apk-11",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Apk:           "",
						Split:         []string{"app-sweet-hdpi.apk", "app-sweet-mdpi.apk", "app-sweet-xhdpi.apk", "app-sweet-xxhdpi.apk"},
						ProductFlavor: "sweet",
						Module:        "module-2",
						BuildType:     "release",
					},
				},
				bitrise.ArtifactListElementResponseModel{
					Title: "app-sweet-xxhdpi.apk",
					Slug:  "test-apk-12",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Apk:           "",
						Split:         []string{"app-sweet-hdpi.apk", "app-sweet-mdpi.apk", "app-sweet-xhdpi.apk", "app-sweet-xxhdpi.apk"},
						ProductFlavor: "sweet",
						Module:        "module-2",
						BuildType:     "release",
					},
				},
				bitrise.ArtifactListElementResponseModel{
					Title: "app-salty-hdpi.apk",
					Slug:  "test-apk-13",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Apk:           "",
						Split:         []string{"app-salty-hdpi.apk", "app-salty-mdpi.apk", "app-salty-xhdpi.apk", "app-salty-xxhdpi.apk"},
						ProductFlavor: "salty",
						Module:        "module-2",
						BuildType:     "release",
					},
				},
				bitrise.ArtifactListElementResponseModel{
					Title: "app-salty-mdpi.apk",
					Slug:  "test-apk-14",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Apk:           "",
						Split:         []string{"app-salty-hdpi.apk", "app-salty-mdpi.apk", "app-salty-xhdpi.apk", "app-salty-xxhdpi.apk"},
						ProductFlavor: "salty",
						Module:        "module-2",
						BuildType:     "release",
					},
				},
				bitrise.ArtifactListElementResponseModel{
					Title: "app-salty-xhdpi.apk",
					Slug:  "test-apk-15",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Apk:           "",
						Split:         []string{"app-salty-hdpi.apk", "app-salty-mdpi.apk", "app-salty-xhdpi.apk", "app-salty-xxhdpi.apk"},
						ProductFlavor: "salty",
						Module:        "module-2",
						BuildType:     "release",
					},
				},
				bitrise.ArtifactListElementResponseModel{
					Title: "app-salty-xxhdpi.apk",
					Slug:  "test-apk-16",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Apk:           "",
						Split:         []string{"app-salty-hdpi.apk", "app-salty-mdpi.apk", "app-salty-xhdpi.apk", "app-salty-xxhdpi.apk"},
						ProductFlavor: "salty",
						Module:        "module-2",
						BuildType:     "release",
					},
				},
			},
			moduleName:    "module-1",
			flavor:        "sweet",
			expectedSlugs: []string{"test-apk-1", "test-apk-2", "test-apk-3", "test-apk-4"},
		},
		{
			testName: "error - release build type - multiple flavor, multiple module - split apk, without module settings",
			artifacts: []bitrise.ArtifactListElementResponseModel{
				bitrise.ArtifactListElementResponseModel{
					Title: "app-sweet-hdpi.apk",
					Slug:  "test-apk-1",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Apk:           "",
						Split:         []string{"app-sweet-hdpi.apk", "app-sweet-mdpi.apk", "app-sweet-xhdpi.apk", "app-sweet-xxhdpi.apk"},
						ProductFlavor: "sweet",
						Module:        "module-1",
						BuildType:     "release",
					},
				},
				bitrise.ArtifactListElementResponseModel{
					Title: "app-sweet-mdpi.apk",
					Slug:  "test-apk-2",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Apk:           "",
						Split:         []string{"app-sweet-hdpi.apk", "app-sweet-mdpi.apk", "app-sweet-xhdpi.apk", "app-sweet-xxhdpi.apk"},
						ProductFlavor: "sweet",
						Module:        "module-1",
						BuildType:     "release",
					},
				},
				bitrise.ArtifactListElementResponseModel{
					Title: "app-sweet-xhdpi.apk",
					Slug:  "test-apk-3",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Apk:           "",
						Split:         []string{"app-sweet-hdpi.apk", "app-sweet-mdpi.apk", "app-sweet-xhdpi.apk", "app-sweet-xxhdpi.apk"},
						ProductFlavor: "sweet",
						Module:        "module-1",
						BuildType:     "release",
					},
				},
				bitrise.ArtifactListElementResponseModel{
					Title: "app-sweet-xxhdpi.apk",
					Slug:  "test-apk-4",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Apk:           "",
						Split:         []string{"app-sweet-hdpi.apk", "app-sweet-mdpi.apk", "app-sweet-xhdpi.apk", "app-sweet-xxhdpi.apk"},
						ProductFlavor: "sweet",
						Module:        "module-1",
						BuildType:     "release",
					},
				},
				bitrise.ArtifactListElementResponseModel{
					Title: "app-salty-hdpi.apk",
					Slug:  "test-apk-5",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Apk:           "",
						Split:         []string{"app-salty-hdpi.apk", "app-salty-mdpi.apk", "app-salty-xhdpi.apk", "app-salty-xxhdpi.apk"},
						ProductFlavor: "salty",
						Module:        "module-1",
						BuildType:     "release",
					},
				},
				bitrise.ArtifactListElementResponseModel{
					Title: "app-salty-mdpi.apk",
					Slug:  "test-apk-6",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Apk:           "",
						Split:         []string{"app-salty-hdpi.apk", "app-salty-mdpi.apk", "app-salty-xhdpi.apk", "app-salty-xxhdpi.apk"},
						ProductFlavor: "salty",
						Module:        "module-1",
						BuildType:     "release",
					},
				},
				bitrise.ArtifactListElementResponseModel{
					Title: "app-salty-xhdpi.apk",
					Slug:  "test-apk-7",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Apk:           "",
						Split:         []string{"app-salty-hdpi.apk", "app-salty-mdpi.apk", "app-salty-xhdpi.apk", "app-salty-xxhdpi.apk"},
						ProductFlavor: "salty",
						Module:        "module-1",
						BuildType:     "release",
					},
				},
				bitrise.ArtifactListElementResponseModel{
					Title: "app-salty-xxhdpi.apk",
					Slug:  "test-apk-8",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Apk:           "",
						Split:         []string{"app-salty-hdpi.apk", "app-salty-mdpi.apk", "app-salty-xhdpi.apk", "app-salty-xxhdpi.apk"},
						ProductFlavor: "salty",
						Module:        "module-1",
						BuildType:     "release",
					},
				},
				bitrise.ArtifactListElementResponseModel{
					Title: "app-sweet-hdpi.apk",
					Slug:  "test-apk-9",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Apk:           "",
						Split:         []string{"app-sweet-hdpi.apk", "app-sweet-mdpi.apk", "app-sweet-xhdpi.apk", "app-sweet-xxhdpi.apk"},
						ProductFlavor: "sweet",
						Module:        "module-2",
						BuildType:     "release",
					},
				},
				bitrise.ArtifactListElementResponseModel{
					Title: "app-sweet-mdpi.apk",
					Slug:  "test-apk-10",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Apk:           "",
						Split:         []string{"app-sweet-hdpi.apk", "app-sweet-mdpi.apk", "app-sweet-xhdpi.apk", "app-sweet-xxhdpi.apk"},
						ProductFlavor: "sweet",
						Module:        "module-2",
						BuildType:     "release",
					},
				},
				bitrise.ArtifactListElementResponseModel{
					Title: "app-sweet-xhdpi.apk",
					Slug:  "test-apk-11",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Apk:           "",
						Split:         []string{"app-sweet-hdpi.apk", "app-sweet-mdpi.apk", "app-sweet-xhdpi.apk", "app-sweet-xxhdpi.apk"},
						ProductFlavor: "sweet",
						Module:        "module-2",
						BuildType:     "release",
					},
				},
				bitrise.ArtifactListElementResponseModel{
					Title: "app-sweet-xxhdpi.apk",
					Slug:  "test-apk-12",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Apk:           "",
						Split:         []string{"app-sweet-hdpi.apk", "app-sweet-mdpi.apk", "app-sweet-xhdpi.apk", "app-sweet-xxhdpi.apk"},
						ProductFlavor: "sweet",
						Module:        "module-2",
						BuildType:     "release",
					},
				},
				bitrise.ArtifactListElementResponseModel{
					Title: "app-salty-hdpi.apk",
					Slug:  "test-apk-13",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Apk:           "",
						Split:         []string{"app-salty-hdpi.apk", "app-salty-mdpi.apk", "app-salty-xhdpi.apk", "app-salty-xxhdpi.apk"},
						ProductFlavor: "salty",
						Module:        "module-2",
						BuildType:     "release",
					},
				},
				bitrise.ArtifactListElementResponseModel{
					Title: "app-salty-mdpi.apk",
					Slug:  "test-apk-14",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Apk:           "",
						Split:         []string{"app-salty-hdpi.apk", "app-salty-mdpi.apk", "app-salty-xhdpi.apk", "app-salty-xxhdpi.apk"},
						ProductFlavor: "salty",
						Module:        "module-2",
						BuildType:     "release",
					},
				},
				bitrise.ArtifactListElementResponseModel{
					Title: "app-salty-xhdpi.apk",
					Slug:  "test-apk-15",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Apk:           "",
						Split:         []string{"app-salty-hdpi.apk", "app-salty-mdpi.apk", "app-salty-xhdpi.apk", "app-salty-xxhdpi.apk"},
						ProductFlavor: "salty",
						Module:        "module-2",
						BuildType:     "release",
					},
				},
				bitrise.ArtifactListElementResponseModel{
					Title: "app-salty-xxhdpi.apk",
					Slug:  "test-apk-16",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Apk:           "",
						Split:         []string{"app-salty-hdpi.apk", "app-salty-mdpi.apk", "app-salty-xhdpi.apk", "app-salty-xxhdpi.apk"},
						ProductFlavor: "salty",
						Module:        "module-2",
						BuildType:     "release",
					},
				},
			},
			expectedSettingsErr: "Multiple modules found, but none was selected in the settings",
		},
	} {
		t.Run(tc.testName, func(t *testing.T) {
			selector := bitrise.NewArtifactSelector(tc.artifacts)
			selectedSlugs, err := selector.Select(tc.moduleName, tc.flavor)
			if tc.expectedSettingsErr != "" {
				require.EqualError(t, err, tc.expectedSettingsErr)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.expectedSlugs, selectedSlugs)
		})
	}
}

func Test_ArtifactSelector_PublishAndShareInfo(t *testing.T) {
	for _, tc := range []struct {
		testName                    string
		appVersion                  models.AppVersion
		artifacts                   []bitrise.ArtifactListElementResponseModel
		expectedPublishAndShareInfo bitrise.PublishAndShareInfo
		expectedErr                 string
	}{
		{
			testName: "ok - build type is release",
			appVersion: models.AppVersion{
				ArtifactInfoData: json.RawMessage(`{"build_type":"release"}`),
			},
			expectedPublishAndShareInfo: bitrise.PublishAndShareInfo{
				PublishEnabled: true,
			},
		},
		{
			testName: "ok - build type is debug",
			appVersion: models.AppVersion{
				ArtifactInfoData: json.RawMessage(`{"build_type":"debug"}`),
			},
			expectedPublishAndShareInfo: bitrise.PublishAndShareInfo{
				PublishEnabled: false,
			},
		},
		{
			testName: "ok - split without universal",
			appVersion: models.AppVersion{
				ProductFlavor:    "sweet",
				ArtifactInfoData: json.RawMessage(`{"build_type":"release","module":"module-1"}`),
			},
			artifacts: []bitrise.ArtifactListElementResponseModel{
				bitrise.ArtifactListElementResponseModel{
					Title: "app-hdpi.apk",
					Slug:  "test-apk-1",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Split:         []string{"app-hdpi.apk", "app-mdpi.apk", "app-xhdpi.apk", "app-xxhdpi.apk"},
						ProductFlavor: "sweet",
						Module:        "module-1",
						BuildType:     "release",
					},
				},
				bitrise.ArtifactListElementResponseModel{
					Title: "app-mdpi.apk",
					Slug:  "test-apk-2",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Split:         []string{"app-hdpi.apk", "app-mdpi.apk", "app-xhdpi.apk", "app-xxhdpi.apk"},
						ProductFlavor: "sweet",
						Module:        "module-1",
						BuildType:     "release",
					},
				},
				bitrise.ArtifactListElementResponseModel{
					Title: "app-xhdpi.apk",
					Slug:  "test-apk-3",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Split:         []string{"app-hdpi.apk", "app-mdpi.apk", "app-xhdpi.apk", "app-xxhdpi.apk"},
						ProductFlavor: "sweet",
						Module:        "module-1",
						BuildType:     "release",
					},
				},
				bitrise.ArtifactListElementResponseModel{
					Title: "app-xxhdpi.apk",
					Slug:  "test-apk-4",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Split:         []string{"app-hdpi.apk", "app-mdpi.apk", "app-xhdpi.apk", "app-xxhdpi.apk"},
						ProductFlavor: "sweet",
						Module:        "module-1",
						BuildType:     "release",
					},
				},
			},
			expectedPublishAndShareInfo: bitrise.PublishAndShareInfo{
				PublishEnabled: true,
				Split:          true,
				PublicInstallPageEnabled:      false,
				PublicInstallPageArtifactSlug: "",
				UniversalAvailable:            false,
			},
		},
		{
			testName: "ok - split with universal",
			appVersion: models.AppVersion{
				ProductFlavor:    "sweet",
				ArtifactInfoData: json.RawMessage(`{"build_type":"release","module":"module-1"}`),
			},
			artifacts: []bitrise.ArtifactListElementResponseModel{
				bitrise.ArtifactListElementResponseModel{
					Title: "app-hdpi.apk",
					Slug:  "test-apk-1",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Split:         []string{"app-hdpi.apk", "app-mdpi.apk", "app-xhdpi.apk", "app-xxhdpi.apk"},
						ProductFlavor: "sweet",
						Module:        "module-1",
						BuildType:     "release",
						Universal:     "/bitrise/my-project/app-universal.apk",
					},
				},
				bitrise.ArtifactListElementResponseModel{
					Title: "app-mdpi.apk",
					Slug:  "test-apk-2",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Split:         []string{"app-hdpi.apk", "app-mdpi.apk", "app-xhdpi.apk", "app-xxhdpi.apk"},
						ProductFlavor: "sweet",
						Module:        "module-1",
						BuildType:     "release",
						Universal:     "/bitrise/my-project/app-universal.apk",
					},
				},
				bitrise.ArtifactListElementResponseModel{
					Title: "app-xhdpi.apk",
					Slug:  "test-apk-3",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Split:         []string{"app-hdpi.apk", "app-mdpi.apk", "app-xhdpi.apk", "app-xxhdpi.apk"},
						ProductFlavor: "sweet",
						Module:        "module-1",
						BuildType:     "release",
						Universal:     "/bitrise/my-project/app-universal.apk",
					},
				},
				bitrise.ArtifactListElementResponseModel{
					Title: "app-xxhdpi.apk",
					Slug:  "test-apk-4",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Split:         []string{"app-hdpi.apk", "app-mdpi.apk", "app-xhdpi.apk", "app-xxhdpi.apk"},
						ProductFlavor: "sweet",
						Module:        "module-1",
						BuildType:     "release",
						Universal:     "/bitrise/my-project/app-universal.apk",
					},
				},
				bitrise.ArtifactListElementResponseModel{
					Title:               "app-universal.apk",
					IsPublicPageEnabled: true,
					Slug:                "test-apk-5",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Split:         []string{"app-hdpi.apk", "app-mdpi.apk", "app-xhdpi.apk", "app-xxhdpi.apk"},
						ProductFlavor: "sweet",
						Module:        "module-1",
						BuildType:     "release",
						Universal:     "/bitrise/my-project/app-universal.apk",
					},
				},
			},
			expectedPublishAndShareInfo: bitrise.PublishAndShareInfo{
				PublishEnabled: true,
				Split:          true,
				PublicInstallPageEnabled:      true,
				PublicInstallPageArtifactSlug: "test-apk-5",
				UniversalAvailable:            true,
			},
		},
		{
			testName: "ok - universal with public install page disabled",
			appVersion: models.AppVersion{
				ProductFlavor:    "sweet",
				ArtifactInfoData: json.RawMessage(`{"build_type":"release","module":"module-1"}`),
			},
			artifacts: []bitrise.ArtifactListElementResponseModel{
				bitrise.ArtifactListElementResponseModel{
					Title:               "app-universal.apk",
					IsPublicPageEnabled: false,
					Slug:                "test-apk-1",
					ArtifactMeta: &bitrise.ArtifactMeta{
						ProductFlavor: "sweet",
						Module:        "module-1",
						BuildType:     "release",
						Universal:     "/bitrise/my-project/app-universal.apk",
					},
				},
			},
			expectedPublishAndShareInfo: bitrise.PublishAndShareInfo{
				PublishEnabled: true,
				Split:          false,
				PublicInstallPageEnabled:      false,
				PublicInstallPageArtifactSlug: "",
				UniversalAvailable:            true,
			},
		},
		{
			testName: "ok - standalone apk",
			appVersion: models.AppVersion{
				ProductFlavor:    "sweet",
				ArtifactInfoData: json.RawMessage(`{"build_type":"release","module":"module-1"}`),
			},
			artifacts: []bitrise.ArtifactListElementResponseModel{
				bitrise.ArtifactListElementResponseModel{
					Title:               "app.apk",
					IsPublicPageEnabled: true,
					Slug:                "test-apk-1",
					ArtifactMeta: &bitrise.ArtifactMeta{
						Apk:           "/bitrise/my-project/app.apk",
						ProductFlavor: "sweet",
						Module:        "module-1",
						BuildType:     "release",
						Universal:     "",
					},
				},
			},
			expectedPublishAndShareInfo: bitrise.PublishAndShareInfo{
				PublishEnabled: true,
				Split:          false,
				PublicInstallPageEnabled:      true,
				PublicInstallPageArtifactSlug: "test-apk-1",
				UniversalAvailable:            false,
			},
		},
	} {
		t.Run(tc.testName, func(t *testing.T) {
			selector := bitrise.NewArtifactSelector(tc.artifacts)
			publishAndShareInfo, err := selector.PublishAndShareInfo(&tc.appVersion)
			if tc.expectedErr != "" {
				require.EqualError(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.expectedPublishAndShareInfo, publishAndShareInfo)
		})
	}
}

func Test_ArtifactSelector_HasAndroidArtifact(t *testing.T) {
	t.Run("ok when there's apk", func(t *testing.T) {
		artifactSelector := bitrise.NewArtifactSelector([]bitrise.ArtifactListElementResponseModel{
			bitrise.ArtifactListElementResponseModel{Title: "app.apk"},
		})
		require.True(t, artifactSelector.HasAndroidArtifact())
	})

	t.Run("ok when there's aab", func(t *testing.T) {
		artifactSelector := bitrise.NewArtifactSelector([]bitrise.ArtifactListElementResponseModel{
			bitrise.ArtifactListElementResponseModel{Title: "app.aab"},
		})
		require.True(t, artifactSelector.HasAndroidArtifact())
	})

	t.Run("ok when there's no android artifact", func(t *testing.T) {
		artifactSelector := bitrise.NewArtifactSelector([]bitrise.ArtifactListElementResponseModel{
			bitrise.ArtifactListElementResponseModel{Title: "app.ipa"},
		})
		require.False(t, artifactSelector.HasAndroidArtifact())
	})
}
