package reverb

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

// SignChannel computes the HMAC-SHA256 auth signature for a channel subscription.
//
// The signed string is:
//   - private channels:  socket_id + ":" + channel
//   - presence channels: socket_id + ":" + channel + ":" + channelData
func SignChannel(secret, socketID, channel, channelData string) string {
	parts := []string{socketID, channel}
	if channelData != "" {
		parts = append(parts, channelData)
	}
	return sign(secret, strings.Join(parts, ":"))
}

// VerifyChannelAuth validates the client-provided auth token for a channel.
//
// The auth token format is "app_key:hex_signature". Verification uses a
// constant-time comparison to prevent timing attacks.
func VerifyChannelAuth(secret, socketID, channel, channelData, auth string) bool {
	parts := strings.SplitN(auth, ":", 2)
	if len(parts) != 2 {
		return false
	}
	provided, err := hex.DecodeString(parts[1])
	if err != nil {
		return false
	}
	expected := signBytes(secret, buildChannelString(socketID, channel, channelData))
	return hmac.Equal(expected, provided)
}

// SignHTTPRequest computes the auth signature for an HTTP API request.
//
// The signed string is:
//
//	METHOD + "\n" + PATH + "\n" + url-encoded sorted query parameters
//	(excluding auth_signature)
func SignHTTPRequest(secret, method, path string, params map[string]string) string {
	return sign(secret, buildHTTPString(method, path, params))
}

// VerifyHTTPRequest validates the HMAC signature on an incoming HTTP API request.
func VerifyHTTPRequest(secret, method, path string, params map[string]string, signature string) bool {
	provided, err := hex.DecodeString(signature)
	if err != nil {
		return false
	}
	expected := signBytes(secret, buildHTTPString(method, path, params))
	return hmac.Equal(expected, provided)
}

// sign returns the hex-encoded HMAC-SHA256 of message using key.
func sign(key, message string) string {
	return hex.EncodeToString(signBytes(key, message))
}

// signBytes returns the raw HMAC-SHA256 digest.
func signBytes(key, message string) []byte {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(message))
	return mac.Sum(nil)
}

// buildChannelString assembles the signed string for channel auth.
func buildChannelString(socketID, channel, channelData string) string {
	parts := []string{socketID, channel}
	if channelData != "" {
		parts = append(parts, channelData)
	}
	return strings.Join(parts, ":")
}

// buildHTTPString assembles the signed string for HTTP API auth.
func buildHTTPString(method, path string, params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		if k != "auth_signature" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	pairs := make([]string, 0, len(keys))
	for _, k := range keys {
		pairs = append(pairs, fmt.Sprintf("%s=%s", k, params[k]))
	}
	queryString := strings.Join(pairs, "&")
	return fmt.Sprintf("%s\n%s\n%s", strings.ToUpper(method), path, queryString)
}
