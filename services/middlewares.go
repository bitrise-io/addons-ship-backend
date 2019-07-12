package services

import (
	"net/http"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/api-utils/middleware"
	"github.com/justinas/alice"
)

func createAuthorizeForAppAccessMiddleware(env *env.AppEnv) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return AuthorizeForAppAccessHandlerFunc(env, h)
	}
}

func createAuthorizeForAppVersionAccessMiddleware(env *env.AppEnv) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return AuthorizeForAppVersionAccessHandlerFunc(env, h)
	}
}

func createAuthorizeForAppVersionScreenshotAccessMiddleware(env *env.AppEnv) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return AuthorizeForAppVersionScreenshotAccessHandlerFunc(env, h)
	}
}

func createAuthenticateWithAddonAccessTokenMiddleware(env *env.AppEnv) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return AuthenticateWithAddonAccessTokenHandlerFunc(env, h)
	}
}

func createAuthorizeForAppDeprovisioningMiddleware(env *env.AppEnv) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return AuthorizeForAppDeprovisioningHandlerFunc(env, h)
	}
}

func createAuthorizeForWebhookHandlingMiddleware(env *env.AppEnv) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return AuthorizeForWebhookHandlerFunc(env, h)
	}
}

func createAuthenticateForWebhookHandlingMiddleware(env *env.AppEnv) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return AuthenticateWithDENSecretHandlerFunc(env, h)
	}
}

func createAuthorizeForAppContactEmailConfirmationMiddleware(env *env.AppEnv) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return AuthorizeForAppContactEmailConfirmationHandlerFunc(env, h)
	}
}

func createAuthorizeForAppContactAccessMiddleware(env *env.AppEnv) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return AuthorizeForAppContactAccessHandlerFunc(env, h)
	}
}

// CommonMiddleware ...
func CommonMiddleware(appEnv *env.AppEnv) alice.Chain {
	baseMiddleware := middleware.CommonMiddleware()

	if appEnv.Environment == env.ServerEnvProduction {
		baseMiddleware = baseMiddleware.Append(
			middleware.CreateRedirectToHTTPSMiddleware(),
		)
	}
	return baseMiddleware.Append(middleware.CreateOptionsRequestTerminatorMiddleware())
}

// AuthenticateForProvisioning ...
func AuthenticateForProvisioning(appEnv *env.AppEnv) alice.Chain {
	return CommonMiddleware(appEnv).Append(createAuthenticateWithAddonAccessTokenMiddleware(appEnv))
}

// AuthenticateForDeprovisioning ...
func AuthenticateForDeprovisioning(appEnv *env.AppEnv) alice.Chain {
	return AuthenticateForProvisioning(appEnv).Append(
		createAuthorizeForAppDeprovisioningMiddleware(appEnv),
	)
}

// AuthorizedAppMiddleware ...
func AuthorizedAppMiddleware(appEnv *env.AppEnv) alice.Chain {
	return CommonMiddleware(appEnv).Append(
		createAuthorizeForAppAccessMiddleware(appEnv),
	)
}

// AuthorizedAppVersionMiddleware ...
func AuthorizedAppVersionMiddleware(appEnv *env.AppEnv) alice.Chain {
	return AuthorizedAppMiddleware(appEnv).Append(
		createAuthorizeForAppVersionAccessMiddleware(appEnv),
	)
}

// AuthorizedAppVersionScreenshotMiddleware ...
func AuthorizedAppVersionScreenshotMiddleware(appEnv *env.AppEnv) alice.Chain {
	return AuthorizedAppVersionMiddleware(appEnv).Append(
		createAuthorizeForAppVersionScreenshotAccessMiddleware(appEnv),
	)
}

// AuthorizeForWebhookHandling ...
func AuthorizeForWebhookHandling(appEnv *env.AppEnv) alice.Chain {
	return CommonMiddleware(appEnv).Append(
		createAuthenticateForWebhookHandlingMiddleware(appEnv),
		createAuthorizeForWebhookHandlingMiddleware(appEnv),
	)
}

// AuthorizeForAppContactEmailConfirmationHandling ...
func AuthorizeForAppContactEmailConfirmationHandling(appEnv *env.AppEnv) alice.Chain {
	return CommonMiddleware(appEnv).Append(
		createAuthorizeForAppContactEmailConfirmationMiddleware(appEnv),
	)
}

// AuthorizedAppContactMiddleware ...
func AuthorizedAppContactMiddleware(appEnv *env.AppEnv) alice.Chain {
	return AuthorizedAppMiddleware(appEnv).Append(
		createAuthorizeForAppContactAccessMiddleware(appEnv),
	)
}
