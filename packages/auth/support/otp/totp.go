package otp

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gollin/packages/auth/support/crypto"
)

const (
	period = 30
	digits = 6
)

// GenerateSecret creates a base32 encoded TOTP secret.
func GenerateSecret() (string, error) {
	raw, err := crypto.RandomString(20)
	if err != nil {
		return "", err
	}

	encoder := base32.StdEncoding.WithPadding(base32.NoPadding)
	return encoder.EncodeToString([]byte(raw))[:32], nil
}

// Validate checks a TOTP code across the supplied skew window.
func Validate(secret string, code string, now time.Time, allowedSkew int) bool {
	for offset := -allowedSkew; offset <= allowedSkew; offset++ {
		candidate, err := Code(secret, now.Add(time.Duration(offset*period)*time.Second))
		if err != nil {
			return false
		}
		if candidate == strings.TrimSpace(code) {
			return true
		}
	}
	return false
}

// Code returns the current TOTP code.
func Code(secret string, now time.Time) (string, error) {
	decoder := base32.StdEncoding.WithPadding(base32.NoPadding)
	key, err := decoder.DecodeString(strings.ToUpper(secret))
	if err != nil {
		return "", fmt.Errorf("decode secret: %w", err)
	}

	counter := uint64(now.Unix() / period)
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, counter)

	mac := hmac.New(sha1.New, key)
	_, _ = mac.Write(buf)
	sum := mac.Sum(nil)

	offset := sum[len(sum)-1] & 0x0f
	binaryCode := (int(sum[offset]&0x7f) << 24) |
		(int(sum[offset+1]) << 16) |
		(int(sum[offset+2]) << 8) |
		int(sum[offset+3])

	value := binaryCode % 1_000_000
	return fmt.Sprintf("%06d", value), nil
}

// OTPAuthURL returns the otpauth enrollment URL.
func OTPAuthURL(issuer string, account string, secret string) string {
	label := url.PathEscape(fmt.Sprintf("%s:%s", issuer, account))
	query := url.Values{
		"secret": []string{secret},
		"issuer": []string{issuer},
		"digits": []string{strconv.Itoa(digits)},
		"period": []string{strconv.Itoa(period)},
	}

	return fmt.Sprintf("otpauth://totp/%s?%s", label, query.Encode())
}
