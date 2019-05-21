package services_test

import "github.com/bitrise-io/addons-ship-backend/bitrise"

type testBitriseAPI struct {
	getArtifactMetadataFn func(string, string, string) (*bitrise.ArtifactMeta, error)
}

func (a *testBitriseAPI) GetArtifactMetadata(authToken, appSlug, buildSlug string) (*bitrise.ArtifactMeta, error) {
	if a.getArtifactMetadataFn == nil {
		panic("You have to override GetArtifactMetadata function in tests")
	}
	return a.getArtifactMetadataFn(authToken, appSlug, buildSlug)
}
