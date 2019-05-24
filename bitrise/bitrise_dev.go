package bitrise

import (
	"time"
)

// APIDev ...
type APIDev struct{}

// GetArtifactMetadata ...
func (a *APIDev) GetArtifactMetadata(authToken, appSlug, buildSlug string) (*ArtifactMeta, error) {
	now := time.Now()
	expirationDate := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), time.UTC)
	return &ArtifactMeta{
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
	}, nil
}
