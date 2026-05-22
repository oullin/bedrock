package oauthserver

import (
	"github.com/bedrock/packages/config"
)

// OAuthServerConfig holds the OAuthServer configuration loaded from a config.Repository.
// Keys follow "oauthserver.*" dot-notation, matching the upstream config/oauthserver.php.
type OAuthServerConfig struct {
	// PrivateKey is the RSA private key PEM string used to sign tokens.
	// Corresponds to PASSPORT_PRIVATE_KEY env / oauthserver.private_key config key.
	PrivateKey string

	// PublicKey is the RSA public key PEM string used to verify tokens.
	// Corresponds to PASSPORT_PUBLIC_KEY env / oauthserver.public_key config key.
	PublicKey string

	// Guard is the authentication guard that OAuthServer uses (default: "web").
	// Corresponds to oauthserver.guard config key.
	Guard string

	// Connection is the database connection name for OAuthServer tables.
	// Corresponds to oauthserver.connection config key.
	Connection string

	// PersonalAccessClientID is the client ID used for personal access tokens.
	// Corresponds to PASSPORT_PERSONAL_ACCESS_CLIENT_ID / oauthserver.personal_access_client.id.
	PersonalAccessClientID string

	// PersonalAccessClientSecret is the client secret for personal access tokens.
	// Corresponds to PASSPORT_PERSONAL_ACCESS_CLIENT_SECRET / oauthserver.personal_access_client.secret.
	PersonalAccessClientSecret string

	// StorageDatabase is the named database connection for OAuth storage.
	// Corresponds to oauthserver.storage.database config key.
	StorageDatabase string
}

// LoadConfig reads OAuthServer configuration from a config.Repository.
// It uses dot-notation keys under the "oauthserver" namespace, matching
// the upstream config/oauthserver.php structure.
func LoadConfig(repo *config.Repository) *OAuthServerConfig {
	guard, _ := repo.String("oauthserver.guard", "web")
	privateKey, _ := repo.String("oauthserver.private_key")
	publicKey, _ := repo.String("oauthserver.public_key")
	connection, _ := repo.String("oauthserver.connection")
	patClientID, _ := repo.String("oauthserver.personal_access_client.id")
	patClientSecret, _ := repo.String("oauthserver.personal_access_client.secret")
	storageDB, _ := repo.String("oauthserver.storage.database")

	return &OAuthServerConfig{
		PrivateKey:                 privateKey,
		PublicKey:                  publicKey,
		Guard:                      guard,
		Connection:                 connection,
		PersonalAccessClientID:     patClientID,
		PersonalAccessClientSecret: patClientSecret,
		StorageDatabase:            storageDB,
	}
}
