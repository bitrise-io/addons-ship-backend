package models

import (
	"strings"

	"github.com/pkg/errors"
)

const (
	validateionErrorPrefix = "VERR:"
)

// NewValidationError ...
func NewValidationError(err string) error {
	return errors.New(validateionErrorPrefix + err)
}

// ValidationErrors ...
func ValidationErrors(errs []error) []error {
	verrs := []error{}
	for _, err := range errs {
		if strings.HasPrefix(err.Error(), validateionErrorPrefix) {
			verrs = append(verrs, errors.New(strings.TrimPrefix(err.Error(), validateionErrorPrefix)))
		}
	}
	return verrs
}
