package bitrise

import (
	"time"

	"github.com/bitrise-io/go-utils/pointers"
)

// APIDev ...
type APIDev struct{}

// GetArtifactData ...
func (a *APIDev) GetArtifactData(authToken, appSlug, buildSlug string) (*ArtifactData, error) {
	now := time.Now()
	expirationDate := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), time.UTC)
	return &ArtifactData{
		Meta: ArtifactMeta{
			AppInfo: AppInfo{
				MinimumOS:         "11.1",
				MinimumSDKVersion: "15.2",
				BundleID:          "test.bundle.id",
				DeviceFamilyList:  []int{1, 2},
				PackageName:       "test_package_name",
			},
			ProvisioningInfo: ProvisioningInfo{
				ExpireDate:       expirationDate,
				DistributionType: "development",
			},
		},
		Slug: "test-app-slug",
	}, nil
}

// GetArtifactPublicInstallPageURL ...
func (a *APIDev) GetArtifactPublicInstallPageURL(authToken, appSlug, buildSlug, artifactSlug string) (string, error) {
	return "http://don.t.go.there", nil
}

// GetAppDetails ...
func (a *APIDev) GetAppDetails(authToken, appSlug string) (*AppDetails, error) {
	return &AppDetails{
		Title:     "The Adventures of Stealy",
		AvatarURL: pointers.NewStringPtr("https://bit.ly/1LixVJu"),
	}, nil
}
