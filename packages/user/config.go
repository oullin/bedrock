package user

import (
	"fmt"
	"strings"
	"time"

	configpkg "github.com/gollin/packages/config"
)

// Config controls default user behavior.
type Config struct {
	DefaultName    string
	NormalizeEmail bool
}

// ConfigFromRepository loads user config from the shared repository.
func ConfigFromRepository(repo *configpkg.Repository) (Config, error) {
	if repo == nil {
		return Config{}, fmt.Errorf("user: config repository is required")
	}

	defaultName, err := repo.String("user.defaults.name")

	if err != nil {
		return Config{}, err
	}

	normalizeEmail, err := repo.Bool("user.defaults.normalize_email")

	if err != nil {
		return Config{}, err
	}

	return Config{
		DefaultName:    strings.TrimSpace(defaultName),
		NormalizeEmail: normalizeEmail,
	}, nil
}

// New creates a configured default user model.
func New(cfg Config, id string, name string, email string, now time.Time) *User {
	resolvedName := strings.TrimSpace(name)

	if resolvedName == "" {
		resolvedName = strings.TrimSpace(cfg.DefaultName)
	}

	resolvedEmail := strings.TrimSpace(email)

	if cfg.NormalizeEmail {
		resolvedEmail = normalizeEmail(resolvedEmail)
	}

	return &User{
		ID:        strings.TrimSpace(id),
		Name:      resolvedName,
		Email:     resolvedEmail,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
