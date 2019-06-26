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

// GetProvisioningProfiles ...
func (a *APIDev) GetProvisioningProfiles(authToken, appSlug string) ([]ProvisioningProfile, error) {
	return []ProvisioningProfile{
		ProvisioningProfile{
			Filename: "prov-profile-1.provisionprofile",
			Slug:     "prov-profile-1-slug",
		},
		ProvisioningProfile{
			Filename: "prov-profile-2.provisionprofile",
			Slug:     "prov-profile-2-slug",
		},
	}, nil
}

// GetCodeSigningIdentities ...
func (a *APIDev) GetCodeSigningIdentities(authToken, appSlug string) ([]CodeSigningIdentity, error) {
	return []CodeSigningIdentity{
		CodeSigningIdentity{
			Filename: "build-certificate-1.cert",
			Slug:     "build-certificate-1-slug",
		},
		CodeSigningIdentity{
			Filename: "build-certificate-2.cert",
			Slug:     "build-certificate-2-slug",
		},
	}, nil
}

// GetAndroidKeystoreFiles ...
func (a *APIDev) GetAndroidKeystoreFiles(authToken, appSlug string) ([]AndroidKeystoreFile, error) {
	return []AndroidKeystoreFile{
		AndroidKeystoreFile{
			Filename: "android-keystore-1.keystore",
			Slug:     "android-keystore-1-slug",
		},
		AndroidKeystoreFile{
			Filename: "android-keystore-2.keystore",
			Slug:     "android-keystore-2-slug",
		},
	}, nil
}

// GetServiceAccountFiles ...
func (a *APIDev) GetServiceAccountFiles(authToken, appSlug string) ([]GenericProjectFile, error) {
	return []GenericProjectFile{
		GenericProjectFile{
			Filename: "service-account-1.json",
			Slug:     "generic-file-1-slug",
		},
		GenericProjectFile{
			Filename: "service-account-2.json",
			Slug:     "generic-file-2-slug",
		},
		GenericProjectFile{
			Filename: "package.json",
			Slug:     "generic-file-3-slug",
		},
	}, nil
}

// TriggerDENTask ...
func (a *APIDev) TriggerDENTask(params TaskParams) (*TriggerResponse, error) {
	return &TriggerResponse{}, nil
}
