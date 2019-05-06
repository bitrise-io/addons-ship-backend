package services_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bitrise-io/addons-ship-backend/services"
	"github.com/c2fo/testify/require"
)

func Test_RootHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	require.NoError(t, services.RootHandler(rr, req))

	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, `{"message":"Welcome to Bitrise Ship Addon!"}`+"\n", rr.Body.String())
}
