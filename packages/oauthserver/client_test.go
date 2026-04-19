package passport_test

import (
	"testing"

	"github.com/bedrock/packages/oauthserver"
)

func TestClientFirstPartyWhenPersonalAccessClient(t *testing.T) {
	c := &oauthserver.Client{PersonalAccessClient: true}

	if !c.FirstParty() {
		t.Error("FirstParty should be true for a personal access client")
	}
}

func TestClientFirstPartyWhenPasswordClient(t *testing.T) {
	c := &oauthserver.Client{PasswordClient: true}

	if !c.FirstParty() {
		t.Error("FirstParty should be true for a password client")
	}
}

func TestClientFirstPartyReturnsFalseForThirdParty(t *testing.T) {
	c := &oauthserver.Client{}

	if c.FirstParty() {
		t.Error("FirstParty should be false for a regular third-party client")
	}
}

func TestClientConfidentialWithSecret(t *testing.T) {
	c := &oauthserver.Client{Secret: "supersecret"}

	if !c.Confidential() {
		t.Error("Confidential should be true when client has a secret")
	}
}

func TestClientConfidentialWithoutSecret(t *testing.T) {
	c := &oauthserver.Client{Secret: ""}

	if c.Confidential() {
		t.Error("Confidential should be false when client has no secret")
	}
}

func TestClientHasGrantType(t *testing.T) {
	c := &oauthserver.Client{GrantTypes: []string{"authorization_code", "refresh_token"}}

	if !c.HasGrantType("authorization_code") {
		t.Error("HasGrantType should return true for a registered grant type")
	}

	if c.HasGrantType("password") {
		t.Error("HasGrantType should return false for an unregistered grant type")
	}
}

func TestClientHasScope(t *testing.T) {
	c := &oauthserver.Client{Scopes: []string{"user:read", "orders:read"}}

	if !c.HasScope("user:read") {
		t.Error("HasScope should return true for a registered scope")
	}

	if c.HasScope("admin") {
		t.Error("HasScope should return false for an unregistered scope")
	}
}

func TestClientHasScopeEmptyScopeList(t *testing.T) {
	c := &oauthserver.Client{}

	if c.HasScope("anything") {
		t.Error("HasScope should return false for a client with no scopes")
	}
}
