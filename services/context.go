package services

import (
	"context"
	"errors"

	ctxpkg "github.com/bitrise-io/api-utils/context"
	uuid "github.com/satori/go.uuid"
)

const (
	// ContextKeyAuthorizedAppID ...
	ContextKeyAuthorizedAppID ctxpkg.RequestContextKey = "ctx-authorized-app-id"
	// ContextKeyAuthorizedAppVersionID ...
	ContextKeyAuthorizedAppVersionID ctxpkg.RequestContextKey = "ctx-authorized-app-version-id"
	// ContextKeyAuthorizedScreenshotID ...
	ContextKeyAuthorizedScreenshotID ctxpkg.RequestContextKey = "ctx-authorized-screenshot-id"
)

// GetAuthorizedAppIDFromContext ...
func GetAuthorizedAppIDFromContext(ctx context.Context) (uuid.UUID, error) {
	id, ok := ctx.Value(ContextKeyAuthorizedAppID).(uuid.UUID)
	if !ok {
		return uuid.UUID{}, errors.New("Authorized App ID not found in Context")
	}
	return id, nil
}

// ContextWithAuthorizedAppID ...
func ContextWithAuthorizedAppID(ctx context.Context, appID uuid.UUID) context.Context {
	return context.WithValue(ctx, ContextKeyAuthorizedAppID, appID)
}

// GetAuthorizedAppVersionIDFromContext ...
func GetAuthorizedAppVersionIDFromContext(ctx context.Context) (uuid.UUID, error) {
	id, ok := ctx.Value(ContextKeyAuthorizedAppVersionID).(uuid.UUID)
	if !ok {
		return uuid.UUID{}, errors.New("Authorized App Version ID not found in Context")
	}
	return id, nil
}

// ContextWithAuthorizedAppVersionID ...
func ContextWithAuthorizedAppVersionID(ctx context.Context, appVersionID uuid.UUID) context.Context {
	return context.WithValue(ctx, ContextKeyAuthorizedAppVersionID, appVersionID)
}

// GetAuthorizedScreenshotIDFromContext ...
func GetAuthorizedScreenshotIDFromContext(ctx context.Context) (uuid.UUID, error) {
	id, ok := ctx.Value(ContextKeyAuthorizedScreenshotID).(uuid.UUID)
	if !ok {
		return uuid.UUID{}, errors.New("Authorized App Version Screenshot ID not found in Context")
	}
	return id, nil
}

// ContextWithAuthorizedScreenshotID ...
func ContextWithAuthorizedScreenshotID(ctx context.Context, screenshotID uuid.UUID) context.Context {
	return context.WithValue(ctx, ContextKeyAuthorizedScreenshotID, screenshotID)
}
