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

// AuthorizedAppMiddleware ...
func AuthorizedAppMiddleware(appEnv *env.AppEnv) alice.Chain {
	commonMiddleware := middleware.CommonMiddleware()

	if appEnv.Environment == env.ServerEnvProduction {
		commonMiddleware = commonMiddleware.Append(
			middleware.CreateRedirectToHTTPSMiddleware(),
		)
	}

	return middleware.CommonMiddleware().Append(
		middleware.CreateOptionsRequestTerminatorMiddleware(),
		createAuthorizeForAppAccessMiddleware(appEnv),
	)
}
