package services

import (
	"context"
	"errors"

	ctxpkg "github.com/bitrise-io/api-utils/context"
	uuid "github.com/satori/go.uuid"
)

const (
	contextKeyAuthorizedAppID ctxpkg.RequestContextKey = "ctx-authorized-app-id"
)

// GetAuthorizedAppIDFromContextErr ...
func GetAuthorizedAppIDFromContextErr(ctx context.Context) (uuid.UUID, error) {
	id, ok := ctx.Value(contextKeyAuthorizedAppID).(uuid.UUID)
	if !ok {
		return uuid.UUID{}, errors.New("Authenticated User ID not found in Context")
	}
	return id, nil
}

// ContextWithAuthorizedAppID ...
func ContextWithAuthorizedAppID(ctx context.Context, appID uuid.UUID) context.Context {
	return context.WithValue(ctx, contextKeyAuthorizedAppID, appID)
}
