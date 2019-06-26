package services_test

import "github.com/bitrise-io/addons-ship-backend/bitrise"

type testBitriseAPI struct {
	getArtifactDataFn          func(string, string, string) (*bitrise.ArtifactData, error)
	getArtifactPublicPageURLFn func(string, string, string, string) (string, error)
	getAppDetailsFn            func(string, string) (*bitrise.AppDetails, error)
	getProvisioningProfilesFn  func(string, string) ([]bitrise.ProvisioningProfile, error)
	getCodeSigningIdentitiesFn func(string, string) ([]bitrise.CodeSigningIdentity, error)
	getAndroidKeystoreFilesFn  func(string, string) ([]bitrise.AndroidKeystoreFile, error)
	getServiceAccountFilesFn   func(string, string) ([]bitrise.GenericProjectFile, error)
	triggerDENTaskFn           func(params bitrise.TaskParams) (*bitrise.TriggerResponse, error)
}

func (a *testBitriseAPI) GetArtifactData(authToken, appSlug, buildSlug string) (*bitrise.ArtifactData, error) {
	if a.getArtifactDataFn == nil {
		panic("You have to override GetArtifactData function in tests")
	}
	return a.getArtifactDataFn(authToken, appSlug, buildSlug)
}

func (a *testBitriseAPI) GetArtifactPublicInstallPageURL(authToken, appSlug, buildSlug, artifactSlug string) (string, error) {
	if a.getArtifactPublicPageURLFn == nil {
		panic("You have to override GetArtifactPublicInstallPageURL function in tests")
	}
	return a.getArtifactPublicPageURLFn(authToken, appSlug, buildSlug, artifactSlug)
}

func (a *testBitriseAPI) GetAppDetails(authToken, appSlug string) (*bitrise.AppDetails, error) {
	if a.getAppDetailsFn == nil {
		panic("You have to override GetAppDetails function in tests")
	}
	return a.getAppDetailsFn(authToken, appSlug)
}

func (a *testBitriseAPI) GetProvisioningProfiles(authToken, appSlug string) ([]bitrise.ProvisioningProfile, error) {
	if a.getProvisioningProfilesFn == nil {
		panic("You have to override GetProvisioningProfiles function in tests")
	}
	return a.getProvisioningProfilesFn(authToken, appSlug)
}

func (a *testBitriseAPI) GetCodeSigningIdentities(authToken, appSlug string) ([]bitrise.CodeSigningIdentity, error) {
	if a.getCodeSigningIdentitiesFn == nil {
		panic("You have to override GetCodeSigningIdentities function in tests")
	}
	return a.getCodeSigningIdentitiesFn(authToken, appSlug)
}

func (a *testBitriseAPI) GetAndroidKeystoreFiles(authToken, appSlug string) ([]bitrise.AndroidKeystoreFile, error) {
	if a.getAndroidKeystoreFilesFn == nil {
		panic("You have to override GetAndroidKeystoreFiles function in tests")
	}
	return a.getAndroidKeystoreFilesFn(authToken, appSlug)
}

func (a *testBitriseAPI) GetServiceAccountFiles(authToken, appSlug string) ([]bitrise.GenericProjectFile, error) {
	if a.getServiceAccountFilesFn == nil {
		panic("You have to override GetServiceAccountFiles function in tests")
	}
	return a.getServiceAccountFilesFn(authToken, appSlug)
}

func (a *testBitriseAPI) TriggerDENTask(params bitrise.TaskParams) (*bitrise.TriggerResponse, error) {
	if a.triggerDENTaskFn == nil {
		panic("You have to override TriggerDENTask function in tests")
	}
	return a.triggerDENTaskFn(params)
}
