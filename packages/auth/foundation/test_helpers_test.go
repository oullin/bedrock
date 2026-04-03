package foundation_test

import (
	"testing"

	"github.com/gollin/packages/auth/memory"
)

func newInMemoryUserRepository(t *testing.T) *memory.InMemoryUserRepository {
	t.Helper()

	repo, err := memory.NewInMemoryUserRepository()

	if err != nil {
		t.Fatalf("NewInMemoryUserRepository: %v", err)
	}

	return repo
}
