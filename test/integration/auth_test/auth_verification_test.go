package auth_test

import (
	"github.com/BrunoMartins11/mid-crowdsensor/cmd/auth"
	"github.com/BrunoMartins11/mid-crowdsensor/test"
	"testing"
)



func TestIsValidToken(t *testing.T) {
	token := test.SetupToken(t)

	testVar := auth.IsValidToken(token)

	if !testVar {
		t.Errorf("Token validation failed")
	}
}
