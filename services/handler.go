package services

import (
	"net/http"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// Handler ...
type Handler struct {
	Env *env.AppEnv
	H   func(e *env.AppEnv, w http.ResponseWriter, r *http.Request) error
}

// ServeHTTP ...
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.H(h.Env, w, r)
	if err != nil {
		if h.Env.Logger != nil {
			h.Env.Logger.Error(" [!] Exception: Internal Server Error", zap.Error(err))
			defer h.Env.Logger.Sync()
		}
		httpresponse.RespondWithInternalServerError(w, errors.WithStack(err))
	}
}
