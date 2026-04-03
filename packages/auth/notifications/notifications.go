package notifications

import auth "github.com/gollin/packages/auth"

// ResetPassword returns a reset-password notification payload.
func ResetPassword(to string, url string, token string) auth.MailMessage {
	return auth.MailMessage{
		To:      to,
		Subject: "Reset your password",
		Body:    "Reset your password by visiting " + url,
		Metadata: map[string]string{
			"url":   url,
			"token": token,
		},
	}
}

// VerifyEmail returns an email-verification notification payload.
func VerifyEmail(to string, url string) auth.MailMessage {
	return auth.MailMessage{
		To:      to,
		Subject: "Verify your email address",
		Body:    "Verify your email by visiting " + url,
		Metadata: map[string]string{
			"url": url,
		},
	}
}
