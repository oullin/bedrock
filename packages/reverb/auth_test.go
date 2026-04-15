package reverb_test

import (
	"testing"

	"github.com/bedrock/packages/reverb"
)

const (
	testSecret   = "test-secret"
	testSocketID = "12345.67890"
	testChannel  = "private-test"
	testAppKey   = "test-key"
)

func TestSignChannel_Private(t *testing.T) {
	t.Parallel()

	sig := reverb.SignChannel(testSecret, testSocketID, testChannel, "")
	auth := testAppKey + ":" + sig

	if !reverb.VerifyChannelAuth(testSecret, testSocketID, testChannel, "", auth) {
		t.Error("VerifyChannelAuth should return true for a valid signature")
	}
}

func TestSignChannel_Presence(t *testing.T) {
	t.Parallel()

	channelData := `{"user_id":"42"}`
	sigWithData := reverb.SignChannel(testSecret, testSocketID, testChannel, channelData)
	sigWithout := reverb.SignChannel(testSecret, testSocketID, testChannel, "")

	if sigWithData == sigWithout {
		t.Error("signature with channelData should differ from signature without channelData")
	}
}

func TestVerifyChannelAuth_Valid(t *testing.T) {
	t.Parallel()

	sig := reverb.SignChannel(testSecret, testSocketID, testChannel, "")
	auth := testAppKey + ":" + sig

	if !reverb.VerifyChannelAuth(testSecret, testSocketID, testChannel, "", auth) {
		t.Error("expected VerifyChannelAuth to return true for correct secret")
	}
}

func TestVerifyChannelAuth_Invalid(t *testing.T) {
	t.Parallel()

	sig := reverb.SignChannel(testSecret, testSocketID, testChannel, "")
	auth := testAppKey + ":" + sig

	if reverb.VerifyChannelAuth("wrong-secret", testSocketID, testChannel, "", auth) {
		t.Error("expected VerifyChannelAuth to return false for wrong secret")
	}
}

func TestVerifyChannelAuth_WrongChannel(t *testing.T) {
	t.Parallel()

	sig := reverb.SignChannel(testSecret, testSocketID, testChannel, "")
	auth := testAppKey + ":" + sig

	if reverb.VerifyChannelAuth(testSecret, testSocketID, "private-other", "", auth) {
		t.Error("expected VerifyChannelAuth to return false for different channel")
	}
}

func TestVerifyHTTPRequest_Valid(t *testing.T) {
	t.Parallel()

	params := map[string]string{
		"auth_key":       testAppKey,
		"auth_timestamp": "1234567890",
		"auth_version":   "1.0",
		"body_md5":       "abc123",
	}

	sig := reverb.SignHTTPRequest(testSecret, "POST", "/apps/123/events", params)

	if !reverb.VerifyHTTPRequest(testSecret, "POST", "/apps/123/events", params, sig) {
		t.Error("expected VerifyHTTPRequest to return true for correct signature")
	}
}

func TestVerifyHTTPRequest_Invalid(t *testing.T) {
	t.Parallel()

	params := map[string]string{
		"auth_key":       testAppKey,
		"auth_timestamp": "1234567890",
	}

	if reverb.VerifyHTTPRequest(testSecret, "POST", "/apps/123/events", params, "badsignature00") {
		t.Error("expected VerifyHTTPRequest to return false for invalid signature")
	}
}

func TestBuildHTTPString_SortedParams(t *testing.T) {
	t.Parallel()

	// Verify that the order of params passed to SignHTTPRequest does not affect the result.
	params1 := map[string]string{
		"auth_key":       testAppKey,
		"auth_timestamp": "1234567890",
		"body_md5":       "abc",
	}
	params2 := map[string]string{
		"body_md5":       "abc",
		"auth_timestamp": "1234567890",
		"auth_key":       testAppKey,
	}

	sig1 := reverb.SignHTTPRequest(testSecret, "GET", "/apps/1/channels", params1)
	sig2 := reverb.SignHTTPRequest(testSecret, "GET", "/apps/1/channels", params2)

	if sig1 != sig2 {
		t.Errorf("expected identical signatures regardless of map iteration order: %q vs %q", sig1, sig2)
	}
}
