package auth

import (
	"testing"
)

func TestCreateJWT(t *testing.T) {
	secret := []byte("secret")
	access := []byte("access")
	accessToken, secretToken, err := CreateJWT(access, secret, 1, "customer", "dummy")
	if err != nil {
		t.Errorf("error creating JWT: %v", err)
	}

	if accessToken == "" {
		t.Error("expected token to be not empty")
	}
	if secretToken == "" {
		t.Error("expected token to be not empty")
	}
	t.Logf("accessToken : %v", accessToken)
	t.Logf("secretToken : %v", secretToken)
}
