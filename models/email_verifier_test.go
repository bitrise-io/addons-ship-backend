package models_test

import (
	"fmt"
	"testing"

	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/c2fo/testify/require"
)

func Test_EmailVerifier_Verify(t *testing.T) {
	for _, tc := range []struct {
		email string
		valid bool
	}{
		{email: "no a valid email", valid: false},
		{email: "#@%^%#$@#$@#.com", valid: false},
		{email: "@domain.com", valid: false},
		{email: "Joe Smith <email@domain.com>", valid: false},
		{email: "email.domain.com", valid: false},
		{email: "email@domain@domain.com", valid: false},
		{email: "email@domain.com (Joe Smith)", valid: false},
		{email: "email@-domain.com", valid: false},
		{email: "email.@domain.com", valid: true},
		{email: "firstname.lastname@domain.com", valid: true},
		{email: "email@subdomain.domain.com", valid: true},
		{email: "firstname+lastname@domain.com", valid: true},
		{email: "email@123.123.123.123", valid: true},
		{email: "1234567890@domain.com", valid: true},
	} {
		t.Run(fmt.Sprintf("%s is valid: %t", tc.email, tc.valid), func(t *testing.T) {
			require.Equal(t, tc.valid, models.EmailVerifier{Email: tc.email}.Verify())
		})
	}
}
