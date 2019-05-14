package services_test

import (
	"net/http"

	ctxpkg "github.com/bitrise-io/api-utils/context"
	"github.com/bitrise-io/api-utils/httpresponse"
)

type testAuthHandler struct {
	ContextElementList map[string]ctxpkg.RequestContextKey
}

func (h *testAuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{}
	for respKey, ctxKey := range h.ContextElementList {
		response[respKey] = r.Context().Value(ctxKey)
	}

	httpresponse.RespondWithSuccessNoErr(w, response)
}
