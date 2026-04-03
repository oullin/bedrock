package fortify

import (
	"fmt"
	"time"

	auth "github.com/gollin/packages/auth"
	"github.com/gollin/packages/auth/fortify/contracts"
	securitycrypto "github.com/gollin/packages/security/crypto"
	securityotp "github.com/gollin/packages/security/otp"
)

// DefaultTwoFactorProvider is the default Fortify two-factor provider.
type DefaultTwoFactorProvider struct{}

// GenerateSecret creates a TOTP secret.

// GenerateRecoveryCodes creates recovery codes.

// Validate validates a TOTP code.

// OTPAuthURL returns the otpauth URL.

// StaticTwoFactorRedirect is a simple two-factor redirect implementation.
type StaticTwoFactorRedirect struct {
	Path string
}

func (DefaultTwoFactorProvider) GenerateSecret() (string, error) {
	return securityotp.GenerateSecret()
}

func (DefaultTwoFactorProvider) GenerateRecoveryCodes() ([]string, error) {
	codes := make([]string, 0, 8)

	for range 8 {
		value, err := securitycrypto.RandomString(8)

		if err != nil {
			return nil, err
		}

		if len(value) > 10 {
			value = value[:10]
		}

		codes = append(codes, value)
	}

	return codes, nil
}

func (DefaultTwoFactorProvider) Validate(secret string, code string, now time.Time, allowedSkew int) bool {
	return securityotp.Validate(secret, code, now, allowedSkew)
}

func (DefaultTwoFactorProvider) OTPAuthURL(issuer string, account string, secret string) string {
	return securityotp.OTPAuthURL(issuer, account, secret)
}

var _ contracts.TwoFactorAuthenticationProvider = DefaultTwoFactorProvider{}

// Redirect returns the configured path.
func (r StaticTwoFactorRedirect) Redirect(_ auth.Authenticatable) string {
	if r.Path == "" {
		return "/two-factor-challenge"
	}

	return r.Path
}

func parseLimiter(spec string) (int, time.Duration, error) {
	switch spec {
	case "", "login":
		return 5, time.Minute, nil
	case "two-factor":
		return 5, time.Minute, nil
	case "verification":
		return 6, time.Minute, nil
	}

	var limit int

	var minutes int

	if _, err := fmt.Sscanf(spec, "%d,%d", &limit, &minutes); err != nil {
		return 0, 0, fmt.Errorf("fortify: invalid limiter spec %q", spec)
	}

	return limit, time.Duration(minutes) * time.Minute, nil
}
