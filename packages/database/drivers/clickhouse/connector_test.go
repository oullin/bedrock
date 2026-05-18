package clickhouse

import (
	"net/url"
	"strings"
	"testing"

	"github.com/bedrock/packages/database"
)

func TestBuildDSN_Defaults(t *testing.T) {
	t.Parallel()

	dsn := buildDSN(database.ConnectionConfig{
		Driver:   DriverName,
		Host:     "localhost",
		Database: "default",
	})

	u, err := url.Parse(dsn)

	if err != nil {
		t.Fatalf("invalid DSN %q: %v", dsn, err)
	}

	if u.Scheme != DriverName {
		t.Fatalf("expected scheme %s, got %s", DriverName, u.Scheme)
	}

	if u.Hostname() != "localhost" {
		t.Fatalf("expected host localhost, got %s", u.Hostname())
	}

	if u.Port() != "9000" {
		t.Fatalf("expected default port 9000, got %s", u.Port())
	}

	if strings.TrimPrefix(u.Path, "/") != "default" {
		t.Fatalf("expected database default, got %s", u.Path)
	}
}

func TestBuildDSN_WithCredentialsAndOptions(t *testing.T) {
	t.Parallel()

	dsn := buildDSN(database.ConnectionConfig{
		Driver:   DriverName,
		Host:     "ch.internal",
		Port:     19000,
		Database: "metrics",
		Username: "writer",
		Password: "s3cret",
		Options: map[string]any{
			"compression": "lz4",
			"secure":      true,
		},
	})

	u, err := url.Parse(dsn)

	if err != nil {
		t.Fatalf("invalid DSN %q: %v", dsn, err)
	}

	if u.Port() != "19000" {
		t.Fatalf("expected port 19000, got %s", u.Port())
	}

	if user := u.User.Username(); user != "writer" {
		t.Fatalf("expected username writer, got %s", user)
	}

	if pw, ok := u.User.Password(); !ok || pw != "s3cret" {
		t.Fatalf("expected password s3cret, got %q (set=%v)", pw, ok)
	}

	q := u.Query()

	if q.Get("compression") != "lz4" {
		t.Fatalf("expected compression=lz4, got %q", q.Get("compression"))
	}

	if q.Get("secure") != "true" {
		t.Fatalf("expected secure=true, got %q", q.Get("secure"))
	}
}

func TestBuildDSN_HTTPProtocol(t *testing.T) {
	t.Parallel()

	dsn := buildDSN(database.ConnectionConfig{
		Driver:   DriverName,
		Host:     "ch.internal",
		Port:     8123,
		Database: "default",
		Options:  map[string]any{"protocol": "http"},
	})

	u, err := url.Parse(dsn)

	if err != nil {
		t.Fatalf("invalid DSN %q: %v", dsn, err)
	}

	if u.Scheme != "http" {
		t.Fatalf("expected scheme http, got %s", u.Scheme)
	}

	// protocol option must not leak into the query string.
	if u.Query().Get("protocol") != "" {
		t.Fatalf("protocol option leaked into query: %s", u.RawQuery)
	}
}

func TestNewConnectorFactory_ReturnsFactory(t *testing.T) {
	t.Parallel()

	factory := NewConnectorFactory()

	if factory == nil {
		t.Fatal("expected non-nil factory")
	}
	// Cannot exercise Connect() in a unit test without a server; covered by
	// the integration suite. Asserting the factory exists is enough here.
}
