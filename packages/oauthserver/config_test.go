package oauthserver_test

import (
	"testing"

	"github.com/bedrock/packages/config"
	"github.com/bedrock/packages/oauthserver"
)

func TestLoadConfigReadsPrivateKey(t *testing.T) {
	repo := config.New(map[string]any{
		"oauthserver.private_key": "-----BEGIN RSA PRIVATE KEY-----\nMIIEow==\n-----END RSA PRIVATE KEY-----",
	})

	cfg := oauthserver.LoadConfig(repo)

	if cfg.PrivateKey != "-----BEGIN RSA PRIVATE KEY-----\nMIIEow==\n-----END RSA PRIVATE KEY-----" {
		t.Errorf("PrivateKey = %q, want PEM string", cfg.PrivateKey)
	}
}

func TestLoadConfigReadsPublicKey(t *testing.T) {
	repo := config.New(map[string]any{
		"oauthserver.public_key": "-----BEGIN PUBLIC KEY-----\nMIIBIjAN\n-----END PUBLIC KEY-----",
	})

	cfg := oauthserver.LoadConfig(repo)

	if cfg.PublicKey != "-----BEGIN PUBLIC KEY-----\nMIIBIjAN\n-----END PUBLIC KEY-----" {
		t.Errorf("PublicKey = %q, want PEM string", cfg.PublicKey)
	}
}

func TestLoadConfigReadsGuard(t *testing.T) {
	repo := config.New(map[string]any{
		"oauthserver.guard": "api",
	})

	cfg := oauthserver.LoadConfig(repo)

	if cfg.Guard != "api" {
		t.Errorf("Guard = %q, want %q", cfg.Guard, "api")
	}
}

func TestLoadConfigDefaultGuard(t *testing.T) {
	repo := config.New(map[string]any{})

	cfg := oauthserver.LoadConfig(repo)

	if cfg.Guard != "web" {
		t.Errorf("Guard = %q, want %q (default)", cfg.Guard, "web")
	}
}

func TestLoadConfigReadsPersonalAccessClient(t *testing.T) {
	repo := config.New(map[string]any{
		"oauthserver.personal_access_client.id":     "client-uuid-123",
		"oauthserver.personal_access_client.secret": "super-secret-value",
	})

	cfg := oauthserver.LoadConfig(repo)

	if cfg.PersonalAccessClientID != "client-uuid-123" {
		t.Errorf("PersonalAccessClientID = %q, want %q", cfg.PersonalAccessClientID, "client-uuid-123")
	}

	if cfg.PersonalAccessClientSecret != "super-secret-value" {
		t.Errorf("PersonalAccessClientSecret = %q, want %q", cfg.PersonalAccessClientSecret, "super-secret-value")
	}
}

func TestLoadConfigReadsConnection(t *testing.T) {
	repo := config.New(map[string]any{
		"oauthserver.connection": "mysql",
	})

	cfg := oauthserver.LoadConfig(repo)

	if cfg.Connection != "mysql" {
		t.Errorf("Connection = %q, want %q", cfg.Connection, "mysql")
	}
}

func TestLoadConfigReadsStorageDatabase(t *testing.T) {
	repo := config.New(map[string]any{
		"oauthserver.storage.database": "oauthserver_db",
	})

	cfg := oauthserver.LoadConfig(repo)

	if cfg.StorageDatabase != "oauthserver_db" {
		t.Errorf("StorageDatabase = %q, want %q", cfg.StorageDatabase, "oauthserver_db")
	}
}

func TestLoadConfigReturnsEmptyStringsForMissingKeys(t *testing.T) {
	repo := config.New(map[string]any{})

	cfg := oauthserver.LoadConfig(repo)

	if cfg.PrivateKey != "" {
		t.Errorf("PrivateKey = %q, want empty", cfg.PrivateKey)
	}

	if cfg.PublicKey != "" {
		t.Errorf("PublicKey = %q, want empty", cfg.PublicKey)
	}

	if cfg.Connection != "" {
		t.Errorf("Connection = %q, want empty", cfg.Connection)
	}

	if cfg.PersonalAccessClientID != "" {
		t.Errorf("PersonalAccessClientID = %q, want empty", cfg.PersonalAccessClientID)
	}

	if cfg.PersonalAccessClientSecret != "" {
		t.Errorf("PersonalAccessClientSecret = %q, want empty", cfg.PersonalAccessClientSecret)
	}

	if cfg.StorageDatabase != "" {
		t.Errorf("StorageDatabase = %q, want empty", cfg.StorageDatabase)
	}
}
