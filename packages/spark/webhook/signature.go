// Package webhook provides webhook handling and signature verification.
package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	signatureHeader = "Paddle-Signature"
	hashAlgorithm   = "h1"
)

// VerifySignature returns middleware that validates Paddle webhook
// signatures using HMAC-SHA256. Mirrors Laravel\Paddle\Http\Middleware\
// VerifyWebhookSignature.
func VerifySignature(secret string, maxDrift time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get(signatureHeader)

			if header == "" {
				http.Error(w, "missing signature", http.StatusForbidden)

				return
			}

			ts, sig, err := parseSignature(header)

			if err != nil {
				http.Error(w, err.Error(), http.StatusForbidden)

				return
			}

			// Check timestamp drift.
			sigTime := time.Unix(ts, 0)

			if time.Since(sigTime).Abs() > maxDrift {
				http.Error(w, "timestamp outside variance window", http.StatusForbidden)

				return
			}

			// Read and verify body.
			body, err := readBody(r)

			if err != nil {
				http.Error(w, "failed to read body", http.StatusBadRequest)

				return
			}

			expected := computeHMAC(secret, ts, body)

			if !hmac.Equal([]byte(sig), []byte(expected)) {
				http.Error(w, "signature mismatch", http.StatusForbidden)

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func parseSignature(header string) (int64, string, error) {
	var ts int64

	var sig string

	parts := strings.Split(header, ";")

	for _, part := range parts {
		kv := strings.SplitN(part, "=", 2)

		if len(kv) != 2 {
			continue
		}

		switch kv[0] {
		case "ts":
			var err error

			ts, err = strconv.ParseInt(kv[1], 10, 64)

			if err != nil {
				return 0, "", fmt.Errorf("malformed timestamp: %w", err)
			}
		case hashAlgorithm:
			sig = kv[1]
		}
	}

	if ts == 0 || sig == "" {
		return 0, "", fmt.Errorf("malformed signature header")
	}

	return ts, sig, nil
}

func computeHMAC(secret string, ts int64, body []byte) string {
	payload := fmt.Sprintf("%d:%s", ts, body)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))

	return hex.EncodeToString(mac.Sum(nil))
}

func readBody(r *http.Request) ([]byte, error) {
	if r.Body == nil {
		return []byte{}, nil
	}

	var buf []byte
	buf, err := readAll(r.Body)

	if err != nil {
		return nil, err
	}

	return buf, nil
}

func readAll(r interface{ Read([]byte) (int, error) }) ([]byte, error) {
	var result []byte
	buf := make([]byte, 4096)

	for {
		n, err := r.Read(buf)

		if n > 0 {
			result = append(result, buf[:n]...)
		}

		if err != nil {
			if err.Error() == "EOF" {
				break
			}

			return nil, err
		}
	}

	return result, nil
}
