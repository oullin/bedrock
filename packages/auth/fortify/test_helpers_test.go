package fortify_test

import (
	"testing"

	auth "github.com/gollin/packages/auth"
	"github.com/gollin/packages/auth/memory"
)

func newDefaultPasswordHasher(t *testing.T) auth.DefaultPasswordHasher {
	t.Helper()

	hasher, err := auth.NewDefaultPasswordHasher()

	if err != nil {
		t.Fatalf("NewDefaultPasswordHasher: %v", err)
	}

	return hasher
}

func newInMemoryUserRepository(t *testing.T, hasher ...auth.PasswordHasher) *memory.InMemoryUserRepository {
	t.Helper()

	repo, err := memory.NewInMemoryUserRepository(hasher...)

	if err != nil {
		t.Fatalf("NewInMemoryUserRepository: %v", err)
	}

	return repo
}
