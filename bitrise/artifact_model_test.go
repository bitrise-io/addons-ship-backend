package bitrise_test

import (
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
