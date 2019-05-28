package services_test

import "github.com/bitrise-io/addons-ship-backend/bitrise"

type testBitriseAPI struct {
	getArtifactDataFn          func(string, string, string) (*bitrise.ArtifactData, error)
	getArtifactPublicPageURLFn func(string, string, string, string) (string, error)
	getAppDetailsFn            func(string, string) (*bitrise.AppDetails, error)
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
