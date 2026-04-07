package spark

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// VerifyWebhookSignature returns an HTTP middleware that validates the
// provider webhook signature header.
//
// The signature header format is: ts=<timestamp>;h1=<hash1>,<hash2>,...
// The HMAC is computed as: HMAC-SHA256("{timestamp}:{raw_body}", secret).
func VerifyWebhookSignature(secret string, maxDrift time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sig := r.Header.Get("Paddle-Signature")
			if sig == "" {
				http.Error(w, "missing signature", http.StatusForbidden)

				return
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "bad request", http.StatusBadRequest)

				return
			}

			ts, hashes := parseSignature(sig)
			if ts == "" || len(hashes) == 0 {
				http.Error(w, "invalid signature format", http.StatusForbidden)

				return
			}

			tsInt, err := strconv.ParseInt(ts, 10, 64)
			if err != nil {
				http.Error(w, "invalid timestamp", http.StatusForbidden)

				return
			}

			drift := math.Abs(float64(time.Now().Unix() - tsInt))
			if drift > maxDrift.Seconds() {
				http.Error(w, "timestamp drift too large", http.StatusForbidden)

				return
			}

			expected := computeHMAC(secret, ts, body)

			verified := false
			for _, h := range hashes {
				if hmac.Equal([]byte(expected), []byte(h)) {
					verified = true

					break
				}
			}

			if !verified {
				http.Error(w, "signature mismatch", http.StatusForbidden)

				return
			}

			// Re-populate the body for downstream handlers.
			r.Body = io.NopCloser(strings.NewReader(string(body)))
			next.ServeHTTP(w, r)
		})
	}
}

func parseSignature(sig string) (string, []string) {
	var ts string
	var hashes []string

	for _, part := range strings.Split(sig, ";") {
		if strings.HasPrefix(part, "ts=") {
			ts = strings.TrimPrefix(part, "ts=")
		} else if strings.HasPrefix(part, "h1=") {
			raw := strings.TrimPrefix(part, "h1=")
			hashes = strings.Split(raw, ",")
		}
	}

	return ts, hashes
}

func computeHMAC(secret, ts string, body []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(fmt.Sprintf("%s:%s", ts, body)))

	return hex.EncodeToString(mac.Sum(nil))
}
