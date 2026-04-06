package tests

import (
	"testing"

	"github.com/bedrock/demo/bootstrap"
	"github.com/bedrock/packages/foundation"
)

func newTestApp(t *testing.T) *foundation.Application {
	t.Helper()

	app, err := bootstrap.New()
	if err != nil {
		t.Fatalf("bootstrap.New: %v", err)
	}

	return app
}
