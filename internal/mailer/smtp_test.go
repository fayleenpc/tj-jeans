package mailer

import (
	"testing"

	"github.com/fayleenpc/tj-jeans/internal/config"
)

func TestGenerateToken(t *testing.T) {
	token, err := GenerateToken()
	if err != nil {
		t.Errorf("error creating verification code: %v", err)
	}

	if token == "" {
		t.Error("expected token to be not empty")
	}
	t.Logf("verification code : %v", token)
}

func TestSendVerificationEmail(t *testing.T) {
	token, err := GenerateToken()
	if err != nil {
		t.Errorf("error creating JWT: %v", err)
	}

	if token == "" {
		t.Error("expected token to be not empty")
	}

	if err := SendVerificationEmail(config.Envs.SMTP_User, token); err != nil {
		t.Errorf("error send verification email, %v", err)
	}

	t.Logf("verification code : %v", token)
}
