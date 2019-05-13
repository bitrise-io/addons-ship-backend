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

// AutorizedAppMiddleware ...
func AutorizedAppMiddleware(env *env.AppEnv) alice.Chain {
	return middleware.CommonMiddleware().Append(
		createAuthorizeForAppAccessMiddleware(env),
	)
}
