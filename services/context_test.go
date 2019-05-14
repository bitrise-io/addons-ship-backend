package services_test

import (
	"context"
	"testing"

	"github.com/bitrise-io/addons-ship-backend/services"
	ctxpkg "github.com/bitrise-io/api-utils/context"
	"github.com/c2fo/testify/require"
	uuid "github.com/satori/go.uuid"
)

func Test_GetAuthorizedAppIDFromContext(t *testing.T) {
	testUUID := uuid.NewV4()

	t.Run("ok", func(t *testing.T) {
		appID, err := services.GetAuthorizedAppIDFromContext(context.WithValue(context.Background(), services.ContextKeyAuthorizedAppID, testUUID))
		require.NoError(t, err)
		require.Equal(t, testUUID, appID)
	})

	t.Run("error - value is not an UUID", func(t *testing.T) {
		appID, err := services.GetAuthorizedAppIDFromContext(context.WithValue(context.Background(), services.ContextKeyAuthorizedAppID, "17"))
		require.Equal(t, "Authorized App ID not found in Context", err.Error())
		require.Equal(t, uuid.UUID{}, appID)
	})

	t.Run("error - wrong key", func(t *testing.T) {
		appID, err := services.GetAuthorizedAppIDFromContext(context.WithValue(context.Background(), ctxpkg.RequestContextKey("WrongKey"), testUUID))
		require.Equal(t, "Authorized App ID not found in Context", err.Error())
		require.Equal(t, uuid.UUID{}, appID)
	})
}

func Test_ContextWithAuthorizedAppID(t *testing.T) {
	testUUID := uuid.NewV4()
	t.Run("ok", func(t *testing.T) {
		contextWithValue := services.ContextWithAuthorizedAppID(context.Background(), testUUID)
		expectedContext := context.WithValue(context.Background(), services.ContextKeyAuthorizedAppID, testUUID)
		require.Equal(t, expectedContext, contextWithValue)
	})

	t.Run("ok - the last set value is the valid", func(t *testing.T) {
		anotherTestUUID := uuid.NewV4()
		previousContext := context.WithValue(context.Background(), services.ContextKeyAuthorizedAppID, testUUID)
		contextWithValue := services.ContextWithAuthorizedAppID(previousContext, anotherTestUUID)
		require.Equal(t, anotherTestUUID, contextWithValue.Value(services.ContextKeyAuthorizedAppID))
	})
}
