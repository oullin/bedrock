package passport

import (
	"github.com/bedrock/packages/config"
)

// PassportConfig holds the Passport configuration loaded from a config.Repository.
// Keys follow "passport.*" dot-notation, matching Laravel's config/passport.php.
type PassportConfig struct {
	// PrivateKey is the RSA private key PEM string used to sign tokens.
	// Corresponds to PASSPORT_PRIVATE_KEY env / passport.private_key config key.
	PrivateKey string

	// PublicKey is the RSA public key PEM string used to verify tokens.
	// Corresponds to PASSPORT_PUBLIC_KEY env / passport.public_key config key.
	PublicKey string

	// Guard is the authentication guard that Passport uses (default: "web").
	// Corresponds to passport.guard config key.
	Guard string

	// Connection is the database connection name for Passport tables.
	// Corresponds to passport.connection config key.
	Connection string

	// PersonalAccessClientID is the client ID used for personal access tokens.
	// Corresponds to PASSPORT_PERSONAL_ACCESS_CLIENT_ID / passport.personal_access_client.id.
	PersonalAccessClientID string

	// PersonalAccessClientSecret is the client secret for personal access tokens.
	// Corresponds to PASSPORT_PERSONAL_ACCESS_CLIENT_SECRET / passport.personal_access_client.secret.
	PersonalAccessClientSecret string

	// StorageDatabase is the named database connection for OAuth storage.
	// Corresponds to passport.storage.database config key.
	StorageDatabase string
}

// LoadConfig reads Passport configuration from a config.Repository.
// It uses dot-notation keys under the "passport" namespace, matching
// Laravel's config/passport.php structure.
func LoadConfig(repo *config.Repository) *PassportConfig {
	guard, _ := repo.String("passport.guard", "web")
	privateKey, _ := repo.String("passport.private_key")
	publicKey, _ := repo.String("passport.public_key")
	connection, _ := repo.String("passport.connection")
	patClientID, _ := repo.String("passport.personal_access_client.id")
	patClientSecret, _ := repo.String("passport.personal_access_client.secret")
	storageDB, _ := repo.String("passport.storage.database")

	return &PassportConfig{
		PrivateKey:                 privateKey,
		PublicKey:                  publicKey,
		Guard:                      guard,
		Connection:                 connection,
		PersonalAccessClientID:     patClientID,
		PersonalAccessClientSecret: patClientSecret,
		StorageDatabase:            storageDB,
	}
}
