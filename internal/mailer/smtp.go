package mailer

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/smtp"

	"github.com/fayleenpc/tj-jeans/internal/config"
)

const (
	smtpServer = "sandbox.smtp.mailtrap.io"
	smtpPort   = "2525"
)

var EmailVerificationTokens = map[string]string{} // Temporary storage for tokens

func GenerateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func SendVerificationEmail(user string, token string) error {
	from := config.Envs.SMTP_User
	subject := "Email Verification From TJ Jeans"
	body := fmt.Sprintf("Please verify your email by clicking the following link: %s/verify?token=%s", config.Envs.PublicHost, token)

	msg := []byte("To: " + user + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body)

	auth := smtp.PlainAuth("", config.Envs.SMTP_User, config.Envs.SMTP_Password, smtpServer)
	return smtp.SendMail(smtpServer+":"+smtpPort, auth, from, []string{user}, msg)
}
