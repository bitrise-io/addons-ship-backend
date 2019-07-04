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
			path: "/provision", middleware: services.AuthenticateForProvisioning(appEnv),
			handler: services.ProvisionHandler, allowedMethods: []string{"POST", "OPTIONS"},
		},
		{
			path: "/provision/{app-slug}", middleware: services.AuthenticateForDeprovisioning(appEnv),
			handler: services.DeprovisionHandler, allowedMethods: []string{"DELETE", "OPTIONS"},
		},
		{
			path: "/apps/{app-slug}", middleware: services.AuthorizedAppMiddleware(appEnv),
			handler: services.AppGetHandler, allowedMethods: []string{"GET", "OPTIONS"},
		},
		{
			path: "/apps/{app-slug}/versions", middleware: services.AuthorizedAppMiddleware(appEnv),
			handler: services.AppVersionsGetHandler, allowedMethods: []string{"GET", "OPTIONS"},
		},
		{
			path: "/apps/{app-slug}/versions/{version-id}", middleware: services.AuthorizedAppVersionMiddleware(appEnv),
			handler: services.AppVersionGetHandler, allowedMethods: []string{"GET", "OPTIONS"},
		},
		{
			path: "/apps/{app-slug}/versions/{version-id}", middleware: services.AuthorizedAppVersionMiddleware(appEnv),
			handler: services.AppVersionPutHandler, allowedMethods: []string{"PUT", "OPTIONS"},
		},
		{
			path: "/apps/{app-slug}/versions/{version-id}/publish", middleware: services.AuthorizedAppVersionMiddleware(appEnv),
			handler: services.AppVersionPublishPostHandler, allowedMethods: []string{"POST", "OPTIONS"},
		},
		{
			path: "/apps/{app-slug}/versions/{version-id}/screenshots", middleware: services.AuthorizedAppVersionMiddleware(appEnv),
			handler: services.ScreenshotsGetHandler, allowedMethods: []string{"GET", "OPTIONS"},
		},
		{
			path: "/apps/{app-slug}/versions/{version-id}/screenshots", middleware: services.AuthorizedAppVersionMiddleware(appEnv),
			handler: services.ScreenshotsPostHandler, allowedMethods: []string{"POST", "OPTIONS"},
		},
		{
			path: "/apps/{app-slug}/versions/{version-id}/screenshots/uploaded", middleware: services.AuthorizedAppVersionMiddleware(appEnv),
			handler: services.ScreenshotsUploadedPatchHandler, allowedMethods: []string{"PATCH", "OPTIONS"},
		},
		{
			path: "/apps/{app-slug}/versions/{version-id}/screenshots/{screenshot-slug}", middleware: services.AuthorizedAppVersionScreenshotMiddleware(appEnv),
			handler: services.ScreenshotDeleteHandler, allowedMethods: []string{"DELETE", "OPTIONS"},
		},
		{
			path: "/apps/{app-slug}/versions/{version-id}/feature-graphic", middleware: services.AuthorizedAppVersionMiddleware(appEnv),
			handler: services.FeatureGraphicGetHandler, allowedMethods: []string{"GET", "OPTIONS"},
		},
		{
			path: "/apps/{app-slug}/versions/{version-id}/feature-graphic", middleware: services.AuthorizedAppVersionMiddleware(appEnv),
			handler: services.FeatureGraphicPostHandler, allowedMethods: []string{"POST", "OPTIONS"},
		},
		{
			path: "/apps/{app-slug}/versions/{version-id}/feature-graphic/uploaded", middleware: services.AuthorizedAppVersionMiddleware(appEnv),
			handler: services.FeatureGraphicUploadedPatchHandler, allowedMethods: []string{"PATCH", "OPTIONS"},
		},
		{
			path: "/apps/{app-slug}/versions/{version-id}/feature-graphic", middleware: services.AuthorizedAppVersionMiddleware(appEnv),
			handler: services.FeatureGraphicDeleteHandler, allowedMethods: []string{"DELETE", "OPTIONS"},
		},
		{
			path: "/apps/{app-slug}/settings", middleware: services.AuthorizedAppMiddleware(appEnv),
			handler: services.AppSettingsGetHandler, allowedMethods: []string{"GET", "OPTIONS"},
		},
		{
			path: "/apps/{app-slug}/settings", middleware: services.AuthorizedAppMiddleware(appEnv),
			handler: services.AppSettingsPatchHandler, allowedMethods: []string{"PATCH", "OPTIONS"},
		},
		{
			path: "/apps/{app-slug}/versions/{version-id}/events", middleware: services.AuthorizedAppVersionMiddleware(appEnv),
			handler: services.AppVersionEventsGetHandler, allowedMethods: []string{"GET", "OPTIONS"},
		},
		{
			path: "/webhook", middleware: services.AuthenticateForWebhookHandling(appEnv),
			handler: services.WebhookPostHandler, allowedMethods: []string{"POST", "OPTIONS"},
		},
	} {
		r.Handle(route.path, route.middleware.Then(services.Handler{Env: appEnv, H: route.handler})).
			Methods(route.allowedMethods...)
	}

	r.NotFoundHandler = middleware.CommonMiddleware().Then(&handlers.NotFoundHandler{})
	return r
}
