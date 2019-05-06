package services

import (
	"net/http"

	"github.com/bitrise-io/api-utils/httpresponse"
)

// RootHandler ...
func RootHandler(w http.ResponseWriter, r *http.Request) error {
	return httpresponse.RespondWithSuccess(w, map[string]string{"message": "Welcome to Bitrise Ship Addon!"})
}
