package router

import (
	"net/http"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/services"
	"github.com/bitrise-io/api-utils/handlers"
	"github.com/bitrise-io/api-utils/middleware"
	"github.com/justinas/alice"
	"gopkg.in/DataDog/dd-trace-go.v1/contrib/gorilla/mux"
)

// New ...
func New(appEnv *env.AppEnv) *mux.Router {
	// StrictSlash: allow "trim slash"; /x/ REDIRECTS to /x
	r := mux.NewRouter(mux.WithServiceName("addons-ship-mux")).StrictSlash(true)

	r.Handle("/", middleware.CommonMiddleware().Then(
		services.Handler{Env: appEnv, H: services.RootHandler})).Methods("GET", "OPTIONS")
	r.Handle("/apps/{app-slug}", services.AuthorizedAppMiddleware(appEnv).Then(
		services.Handler{Env: appEnv, H: services.AppGetHandler})).Methods("GET", "OPTIONS")
	r.Handle("/apps/{app-slug}/versions", services.AuthorizedAppMiddleware(appEnv).Then(
		services.Handler{Env: appEnv, H: services.AppVersionsGetHandler})).Methods("GET", "OPTIONS")
	r.NotFoundHandler = middleware.CommonMiddleware().Then(&handlers.NotFoundHandler{})

	for _, route := range []struct {
		path           string
		middleware     alice.Chain
		handler        func(e *env.AppEnv, w http.ResponseWriter, r *http.Request) error
		allowedMethods []string
	}{
		{
			path: "/", middleware: middleware.CommonMiddleware(),
			handler: services.RootHandler, allowedMethods: []string{"GET", "OPTIONS"},
		},
		{
			path: "/apps/{app-slug}", middleware: services.AuthorizedAppMiddleware(appEnv),
			handler: services.AppGetHandler, allowedMethods: []string{"GET", "OPTIONS"},
		},
		{
			path: "/apps/{app-slug}/version", middleware: services.AuthorizedAppMiddleware(appEnv),
			handler: services.AppVersionsGetHandler, allowedMethods: []string{"GET", "OPTIONS"},
		},
	} {
		r.Handle(route.path, route.middleware.Then(services.Handler{Env: appEnv, H: route.handler})).
			Methods(route.allowedMethods...)
	}
	return r
}
