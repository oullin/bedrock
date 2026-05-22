package oauthserver_test

import (
	"encoding/json"
	"testing"

	"github.com/bedrock/packages/oauthserver"
)

func TestScopeFields(t *testing.T) {
	s := oauthserver.Scope{ID: "orders:read", Description: "Read orders"}

	if s.ID != "orders:read" {
		t.Errorf("ID = %q, want %q", s.ID, "orders:read")
	}

	if s.Description != "Read orders" {
		t.Errorf("Description = %q, want %q", s.Description, "Read orders")
	}
}

func TestScopeToArray(t *testing.T) {
	s := oauthserver.Scope{ID: "user:read", Description: "Read user data"}
	got := s.ToArray()

	if got["id"] != "user:read" {
		t.Errorf("ToArray id = %q, want %q", got["id"], "user:read")
	}

	if got["description"] != "Read user data" {
		t.Errorf("ToArray description = %q, want %q", got["description"], "Read user data")
	}
}

func TestScopeToJSON(t *testing.T) {
	s := oauthserver.Scope{ID: "admin", Description: "Admin access"}

	b, err := json.Marshal(s)

	if err != nil {
		t.Fatal(err)
	}

	var got map[string]string

	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatal(err)
	}

	if got["id"] != "admin" {
		t.Errorf("JSON id = %q, want %q", got["id"], "admin")
	}

	if got["description"] != "Admin access" {
		t.Errorf("JSON description = %q, want %q", got["description"], "Admin access")
	}
}

func TestScopeEmptyDescription(t *testing.T) {
	s := oauthserver.Scope{ID: "write"}
	got := s.ToArray()

	if got["description"] != "" {
		t.Errorf("expected empty description, got %q", got["description"])
	}
}
