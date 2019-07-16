package models_test

import (
	"os"
	"testing"

	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/go-crypto/crypto"
	"github.com/bitrise-io/go-utils/envutil"
	"github.com/c2fo/testify/require"
)

func Test_App_Secret(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		secret := "my-super-secret"
		encryptKey := "06042e86a7bd421c642c8c3e4ab13840"
		revokeFn, err := envutil.RevokableSetenv("APP_WEBHOOK_SECRET_ENCRYPT_KEY", encryptKey)
		require.NoError(t, err)

		iv, err := crypto.GenerateIV()
		require.NoError(t, err)
		encryptedSecret, err := crypto.AES256GCMCipher(secret, iv, encryptKey)
		require.NoError(t, err)

		testApp := models.App{EncryptedSecret: encryptedSecret, EncryptedSecretIV: iv}

		calculatedSecret, err := testApp.Secret()
		require.NoError(t, err)
		require.Equal(t, secret, calculatedSecret)
		require.NoError(t, revokeFn())
	})

	t.Run("when no encrypt key set in env var", func(t *testing.T) {
		os.Clearenv()
		testApp := models.App{}
		calculatedSecret, err := testApp.Secret()
		require.EqualError(t, err, "No encrypt key provided")
		require.Empty(t, calculatedSecret)
	})

	t.Run("when no encrypt key set in env var", func(t *testing.T) {
		revokeFn, err := envutil.RevokableSetenv("APP_WEBHOOK_SECRET_ENCRYPT_KEY", "06042e86a7bd421c642c8c3e4ab13840")
		require.NoError(t, err)
		iv, err := crypto.GenerateIV()
		require.NoError(t, err)
		testApp := models.App{EncryptedSecretIV: iv}
		calculatedSecret, err := testApp.Secret()
		require.EqualError(t, err, "cipher: message authentication failed")
		require.Empty(t, calculatedSecret)
		require.NoError(t, revokeFn())
	})
}
