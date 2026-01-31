package main

import (
	"fmt"
	"os"

	veilmail "github.com/Resonia-Health/veilmail-go"
)

func getMailClient() *veilmail.Client {
	return veilmail.NewClient(os.Getenv("VEILMAIL_API_KEY"))
}

func fromEmail() string {
	e := os.Getenv("VEILMAIL_FROM_EMAIL")
	if e == "" {
		return "noreply@veilmail.xyz"
	}
	return e
}

func appURL() string {
	u := os.Getenv("APP_URL")
	if u == "" {
		return "http://localhost:8080"
	}
	return u
}

func sendVerificationEmail(email, token string) {
	url := fmt.Sprintf("%s/auth/verify-email?token=%s", appURL(), token)
	getMailClient().Emails.Send(veilmail.SendEmailParams{
		From:    fromEmail(),
		To:      email,
		Subject: "Verify your email address",
		Html:    fmt.Sprintf(`<p>Click <a href="%s">here</a> to verify your email. This link expires in 1 hour.</p>`, url),
		Tags:    []string{"auth", "verification"},
	})
}

func sendPasswordResetEmail(email, token string) {
	url := fmt.Sprintf("%s/auth/reset-password?token=%s", appURL(), token)
	getMailClient().Emails.Send(veilmail.SendEmailParams{
		From:    fromEmail(),
		To:      email,
		Subject: "Reset your password",
		Html:    fmt.Sprintf(`<p>Click <a href="%s">here</a> to reset your password. This link expires in 1 hour.</p>`, url),
		Tags:    []string{"auth", "password-reset"},
	})
}

func sendTwoFactorCode(email, code string) {
	getMailClient().Emails.Send(veilmail.SendEmailParams{
		From:    fromEmail(),
		To:      email,
		Subject: fmt.Sprintf("%s is your verification code", code),
		Html:    fmt.Sprintf(`<p>Your code: <strong>%s</strong></p><p>Expires in 5 minutes.</p>`, code),
		Tags:    []string{"auth", "2fa"},
	})
}

func sendWelcomeEmail(email, name string) {
	getMailClient().Emails.Send(veilmail.SendEmailParams{
		From:    fromEmail(),
		To:      email,
		Subject: "Welcome!",
		Html:    fmt.Sprintf(`<p>Welcome, %s! Your account is active.</p>`, name),
		Tags:    []string{"auth", "welcome"},
	})
}

func sendPasswordChangedEmail(email string) {
	getMailClient().Emails.Send(veilmail.SendEmailParams{
		From:    fromEmail(),
		To:      email,
		Subject: "Your password was changed",
		Html:    `<p>Your password was changed. If you didn't do this, reset it immediately.</p>`,
		Tags:    []string{"auth", "security"},
	})
}

func send2FAToggledEmail(email string, enabled bool) {
	status := "enabled"
	if !enabled {
		status = "disabled"
	}
	getMailClient().Emails.Send(veilmail.SendEmailParams{
		From:    fromEmail(),
		To:      email,
		Subject: fmt.Sprintf("Two-factor authentication %s", status),
		Html:    fmt.Sprintf(`<p>2FA has been %s on your account.</p>`, status),
		Tags:    []string{"auth", "2fa", "security"},
	})
}
