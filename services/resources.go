package services

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/api-utils/proxy"
	"github.com/pkg/errors"
)

var bitriseAPIVersion = "v0.1"

// ResourcesHandler ...
func ResourcesHandler(env *env.AppEnv, w http.ResponseWriter, r *http.Request) error {
	authorizedAppID, err := GetAuthorizedAppIDFromContext(r.Context())
	if err != nil {
		return errors.WithStack(err)
	}
	if env.AppService == nil {
		return errors.New("No App Service defined for handler")
	}
	app, err := env.AppService.Find(&models.App{Record: models.Record{ID: authorizedAppID}})
	if err != nil {
		return errors.Wrap(err, "SQL Error")
	}
	path := "/" + bitriseAPIVersion + strings.Replace(r.URL.Path, "/resources", "", -1)
	proxyHandler := proxy.NewSingleEndpointSameHostReverseProxyHandler(&url.URL{
		Scheme: env.BitriseAPIRootURL.Scheme,
		Host:   env.BitriseAPIRootURL.Host,
		Path:   path,
	}, &r.Body, map[string]string{"Content-Type": "application/json", "Bitrise-Addon-Auth-Token": app.BitriseAPIToken})
	proxyHandler.ServeHTTP(w, r)

	return nil
}
