package bitrise_test

import (
	"fmt"
	"testing"

	"github.com/bitrise-io/addons-ship-backend/bitrise"
	"github.com/stretchr/testify/require"
)

func Test_ArtifactListElementResponseModel_IsXCodeArchive(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		require.True(t, bitrise.ArtifactListElementResponseModel{Title: "app.xcarchive.zip"}.IsXCodeArchive())
	})
	t.Run("when it's not xcarchive.zip", func(t *testing.T) {
		require.False(t, bitrise.ArtifactListElementResponseModel{Title: "export_options.plist"}.IsXCodeArchive())
	})
}

func Test_ArtifactListElementResponseModel_IsIPA(t *testing.T) {
	require.True(t, bitrise.ArtifactListElementResponseModel{Title: "app.ipa"}.IsIPA())
	require.False(t, bitrise.ArtifactListElementResponseModel{Title: "app.apk"}.IsIPA())
}

func Test_ArtifactListElementResponseModel_IsAAB(t *testing.T) {
	require.True(t, bitrise.ArtifactListElementResponseModel{
		Title:        "app.aab",
		ArtifactMeta: &bitrise.ArtifactMeta{Aab: "/somewhere/over/the-rainbow/app.aab"}}.IsAAB())
	require.False(t, bitrise.ArtifactListElementResponseModel{Title: "app.apk"}.IsAAB())
}

func Test_ArtifactListElementResponseModel_IsStandaloneAPK(t *testing.T) {
	require.True(t, bitrise.ArtifactListElementResponseModel{
		Title:        "app.apk",
		ArtifactMeta: &bitrise.ArtifactMeta{Apk: "/somewhere/over/the-rainbow/app.apk"}}.IsStandaloneAPK())
	require.False(t, bitrise.ArtifactListElementResponseModel{Title: "app.apk"}.IsStandaloneAPK())
	require.False(t, bitrise.ArtifactListElementResponseModel{Title: "app.exe"}.IsStandaloneAPK())
}

func Test_ArtifactListElementResponseModel_IsUniversalAPK(t *testing.T) {
	require.True(t, bitrise.ArtifactListElementResponseModel{
		Title:        "app.universal.apk",
		ArtifactMeta: &bitrise.ArtifactMeta{Universal: "/somewhere/over/the-rainbow/app.universal.apk"}}.IsUniversalAPK())
	require.False(t, bitrise.ArtifactListElementResponseModel{Title: "app.apk"}.IsUniversalAPK())
}

func Test_HasDebugIPAExportMethod(t *testing.T) {
	for _, exportMethod := range []string{"development", "ad-hoc", "enterprise"} {
		t.Run(fmt.Sprintf("when distribution type is %s", exportMethod), func(t *testing.T) {
			require.True(t, bitrise.ArtifactListElementResponseModel{
				ArtifactMeta: &bitrise.ArtifactMeta{
					ProvisioningInfo: bitrise.ProvisioningInfo{
						IPAExportMethod: exportMethod,
					},
				},
			}.HasDebugIPAExportMethod())
		})
	}
}
