package services_test

import "github.com/bitrise-io/addons-ship-backend/bitrise"

type testBitriseAPI struct {
	getArtifactDataFn          func(string, string, string) (*bitrise.ArtifactData, error)
	getArtifactsFn             func(string, string, string) ([]bitrise.ArtifactListElementResponseModel, error)
	getArtifactFn              func(string, string, string, string) (*bitrise.ArtifactShowResponseItemModel, error)
	getArtifactPublicPageURLFn func(string, string, string, string) (string, error)
	getAppDetailsFn            func(string, string) (*bitrise.AppDetails, error)
	getProvisioningProfilesFn  func(string, string) ([]bitrise.ProvisioningProfile, error)
	getProvisioningProfileFn   func(string, string, string) (*bitrise.ProvisioningProfile, error)
	getCodeSigningIdentitiesFn func(string, string) ([]bitrise.CodeSigningIdentity, error)
	getCodeSigningIdentityFn   func(string, string, string) (*bitrise.CodeSigningIdentity, error)
	getAndroidKeystoreFilesFn  func(string, string) ([]bitrise.AndroidKeystoreFile, error)
	getAndroidKeystoreFileFn   func(string, string, string) (*bitrise.AndroidKeystoreFile, error)
	getServiceAccountFilesFn   func(string, string) ([]bitrise.GenericProjectFile, error)
	getServiceAccountFileFn    func(string, string, string) (*bitrise.GenericProjectFile, error)
	triggerDENTaskFn           func(params bitrise.TaskParams) (*bitrise.TriggerResponse, error)
	registerWebhookFn          func(string, string, string, string) error
}

func (a *testBitriseAPI) GetArtifactData(authToken, appSlug, buildSlug string) (*bitrise.ArtifactData, error) {
	if a.getArtifactDataFn == nil {
		panic("You have to override GetArtifactData function in tests")
	}
	return a.getArtifactDataFn(authToken, appSlug, buildSlug)
}

func (a *testBitriseAPI) GetArtifacts(authToken, appSlug, buildSlug string) ([]bitrise.ArtifactListElementResponseModel, error) {
	if a.getArtifactsFn == nil {
		panic("You have to override GetArtifacts function in tests")
	}
	return a.getArtifactsFn(authToken, appSlug, buildSlug)
}

func (a *testBitriseAPI) GetArtifact(authToken, appSlug, buildSlug, artifactSlug string) (*bitrise.ArtifactShowResponseItemModel, error) {
	if a.getArtifactFn == nil {
		panic("You have to override BitriseAPI.GetArtifact function in tests")
	}
	return a.getArtifactFn(authToken, appSlug, buildSlug, artifactSlug)
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

func (a *testBitriseAPI) GetProvisioningProfile(authToken, appSlug, provProfileSlug string) (*bitrise.ProvisioningProfile, error) {
	if a.getProvisioningProfileFn == nil {
		panic("You have to override GetProvisioningProfile function in tests")
	}
	return a.getProvisioningProfileFn(authToken, appSlug, provProfileSlug)
}

func (a *testBitriseAPI) GetCodeSigningIdentities(authToken, appSlug string) ([]bitrise.CodeSigningIdentity, error) {
	if a.getCodeSigningIdentitiesFn == nil {
		panic("You have to override GetCodeSigningIdentities function in tests")
	}
	return a.getCodeSigningIdentitiesFn(authToken, appSlug)
}

func (a *testBitriseAPI) GetCodeSigningIdentity(authToken, appSlug, codeSigningSlug string) (*bitrise.CodeSigningIdentity, error) {
	if a.getCodeSigningIdentityFn == nil {
		panic("You have to override GetCodeSigningIdentity function in tests")
	}
	return a.getCodeSigningIdentityFn(authToken, appSlug, codeSigningSlug)
}

func (a *testBitriseAPI) GetAndroidKeystoreFiles(authToken, appSlug string) ([]bitrise.AndroidKeystoreFile, error) {
	if a.getAndroidKeystoreFilesFn == nil {
		panic("You have to override GetAndroidKeystoreFiles function in tests")
	}
	return a.getAndroidKeystoreFilesFn(authToken, appSlug)
}

func (a *testBitriseAPI) GetAndroidKeystoreFile(authToken, appSlug, keystoreSlug string) (*bitrise.AndroidKeystoreFile, error) {
	if a.getAndroidKeystoreFileFn == nil {
		panic("You have to override GetAndroidKeystoreFile function in tests")
	}
	return a.getAndroidKeystoreFileFn(authToken, appSlug, keystoreSlug)
}

func (a *testBitriseAPI) GetServiceAccountFiles(authToken, appSlug string) ([]bitrise.GenericProjectFile, error) {
	if a.getServiceAccountFilesFn == nil {
		panic("You have to override GetServiceAccountFiles function in tests")
	}
	return a.getServiceAccountFilesFn(authToken, appSlug)
}

func (a *testBitriseAPI) GetServiceAccountFile(authToken, appSlug, serviceJSONSLug string) (*bitrise.GenericProjectFile, error) {
	if a.getServiceAccountFileFn == nil {
		panic("You have to override GetServiceAccountFile function in tests")
	}
	return a.getServiceAccountFileFn(authToken, appSlug, serviceJSONSLug)
}

func (a *testBitriseAPI) TriggerDENTask(params bitrise.TaskParams) (*bitrise.TriggerResponse, error) {
	if a.triggerDENTaskFn == nil {
		panic("You have to override TriggerDENTask function in tests")
	}
	return a.triggerDENTaskFn(params)
}

func (a *testBitriseAPI) RegisterWebhook(authToken, appSlug, secret, callbackURL string) error {
	if a.registerWebhookFn == nil {
		panic("You have to override RegisterWebhook function in tests")
	}
	return a.registerWebhookFn(authToken, appSlug, secret, callbackURL)
}
