package router

import (
	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/services"
	"github.com/bitrise-io/api-utils/handlers"
	"github.com/bitrise-io/api-utils/middleware"
	"gopkg.in/DataDog/dd-trace-go.v1/contrib/gorilla/mux"
)

// New ...
func New(env *env.AppEnv) *mux.Router {
	// StrictSlash: allow "trim slash"; /x/ REDIRECTS to /x
	r := mux.NewRouter(mux.WithServiceName("addons-ship-mux")).StrictSlash(true)

	r.Handle("/", middleware.CommonMiddleware().Then(
		services.Handler{Env: env, H: services.RootHandler})).Methods("GET", "OPTIONS")
	r.Handle("/apps/{app-slug}", services.AutorizedAppMiddleware(env).Then(
		services.Handler{Env: env, H: services.AppGetHandler})).Methods("GET", "OPTIONS")
	// r.Handle("/apps/{app-slug}/versions", services.AutorizedAppMiddleware(env).Then(
	// 	services.Handler{Env: env, H: services.AppsGetHandler})).Methods("GET", "OPTIONS")
	r.NotFoundHandler = middleware.CommonMiddleware().Then(&handlers.NotFoundHandler{})
	return r
}
