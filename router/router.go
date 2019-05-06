package router

import (
	"github.com/bitrise-io/addons-ship-backend/services"
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/bitrise-io/api-utils/middleware"
	"gopkg.in/DataDog/dd-trace-go.v1/contrib/gorilla/mux"
)

// New ...
func New() *mux.Router {
	// StrictSlash: allow "trim slash"; /x/ REDIRECTS to /x
	r := mux.NewRouter(mux.WithServiceName("addons-ship-mux")).StrictSlash(true)

	r.Handle("/", middleware.CommonMiddleware().Then(
		httpresponse.InternalErrHandlerFuncAdapter(services.RootHandler))).Methods("GET", "OPTIONS")
	return r
}
