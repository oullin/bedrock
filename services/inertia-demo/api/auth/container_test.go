package auth

import (
	"database/sql"
	"net/http"
	"strings"
	"testing"

	"github.com/bedrock/packages/inertia/flash"
	"github.com/bedrock/packages/inertia/protocol"
)

func TestContainerValidate_ZeroValue(t *testing.T) {
	err := Container{}.Validate()

	if err == nil {
		t.Fatal("expected error for zero-value container")
	}

	for _, field := range []string{"DB", "CryptoKey", "Render", "Redirect", "RouteURL", "SetFlash"} {
		if !strings.Contains(err.Error(), field) {
			t.Errorf("error should mention %s, got: %s", field, err.Error())
		}
	}
}

func TestContainerValidate_Valid(t *testing.T) {
	renderFn := func(http.ResponseWriter, *http.Request, string, protocol.Props) {}
	redirectFn := func(http.ResponseWriter, *http.Request, string) {}
	routeURLFn := func(string, map[string]string) string { return "" }
	setFlashFn := func(http.ResponseWriter, flash.Message) error { return nil }

	c := Container{
		DB:        &sql.DB{},
		CryptoKey: []byte("secret"),
		Render:    renderFn,
		Redirect:  redirectFn,
		RouteURL:  routeURLFn,
		SetFlash:  setFlashFn,
	}

	if err := c.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestContainerValidate_SecureCookieFalseIsValid(t *testing.T) {
	renderFn := func(http.ResponseWriter, *http.Request, string, protocol.Props) {}
	redirectFn := func(http.ResponseWriter, *http.Request, string) {}
	routeURLFn := func(string, map[string]string) string { return "" }
	setFlashFn := func(http.ResponseWriter, flash.Message) error { return nil }

	c := Container{
		DB:           &sql.DB{},
		CryptoKey:    []byte("secret"),
		Render:       renderFn,
		Redirect:     redirectFn,
		RouteURL:     routeURLFn,
		SetFlash:     setFlashFn,
		SecureCookie: false,
	}

	if err := c.Validate(); err != nil {
		t.Fatalf("SecureCookie=false should be valid, got: %v", err)
	}
}

func TestContainerValidate_EmptyCryptoKey(t *testing.T) {
	renderFn := func(http.ResponseWriter, *http.Request, string, protocol.Props) {}
	redirectFn := func(http.ResponseWriter, *http.Request, string) {}
	routeURLFn := func(string, map[string]string) string { return "" }
	setFlashFn := func(http.ResponseWriter, flash.Message) error { return nil }

	c := Container{
		DB:        &sql.DB{},
		CryptoKey: []byte{},
		Render:    renderFn,
		Redirect:  redirectFn,
		RouteURL:  routeURLFn,
		SetFlash:  setFlashFn,
	}

	err := c.Validate()

	if err == nil {
		t.Fatal("expected error for empty CryptoKey")
	}

	if !strings.Contains(err.Error(), "CryptoKey") {
		t.Errorf("error should mention CryptoKey, got: %s", err.Error())
	}
}
