package services_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/services"
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/c2fo/testify/require"
	"github.com/pkg/errors"
)

func Test_Handler_ServeHTTP(t *testing.T) {
	t.Run("when handler responds with non-5xx error", func(t *testing.T) {
		handler := services.Handler{
			Env: &env.AppEnv{},
			H: func(env *env.AppEnv, w http.ResponseWriter, r *http.Request) error {
				return httpresponse.RespondWithSuccess(w, "ok")
			},
		}
		r, err := http.NewRequest("GET", "...", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, r)
		require.Equal(t, `"ok"`+"\n", rr.Body.String())
	})

	t.Run("when handler responds with 5xx error", func(t *testing.T) {
		handler := services.Handler{
			Env: &env.AppEnv{},
			H: func(env *env.AppEnv, w http.ResponseWriter, r *http.Request) error {
				return errors.New("Some internal error")
			},
		}
		r, err := http.NewRequest("GET", "...", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, r)
		require.Equal(t, `{"message":"Internal Server Error"}`+"\n", rr.Body.String())
	})
}
