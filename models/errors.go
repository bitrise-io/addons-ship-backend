package models

import (
	"strings"

	"github.com/pkg/errors"
)

const (
	validationErrorPrefix = "VERR:"
)

// NewValidationError ...
func NewValidationError(err string) error {
	return errors.New(validationErrorPrefix + err)
}

// ValidationErrors ...
func ValidationErrors(errs []error) []error {
	verrs := []error{}
	for _, err := range errs {
		if strings.HasPrefix(err.Error(), validationErrorPrefix) {
			verrs = append(verrs, errors.New(strings.TrimPrefix(err.Error(), validationErrorPrefix)))
		}
	}
	return verrs
}
