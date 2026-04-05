package tests

import (
	"testing"

	"github.com/gollin/demo/bootstrap"
	"github.com/gollin/packages/illuminate/foundation"
)

func newTestApp(t *testing.T) *foundation.Application {
	t.Helper()

	app, err := bootstrap.New()
	if err != nil {
		t.Fatalf("bootstrap.New: %v", err)
	}

	return app
}
