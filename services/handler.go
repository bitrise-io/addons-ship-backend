package services

import (
	"fmt"
	"net/http"
	"time"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	presignedURLExpirationInterval = 10 * time.Minute
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
			defer func() {
				err := h.Env.Logger.Sync()
				if err != nil {
					fmt.Printf("Failed to sync logger: %#v", err)
				}
			}()
		}
		httpresponse.RespondWithInternalServerError(w, errors.WithStack(err))
	}
}
