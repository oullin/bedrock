package tests

import (
	"io"
	"net/http"
	"testing"
)

type testApp struct {
	handler http.Handler
	config  *testConfig
}

type testConfig struct {
	values map[string]string
}

func (c *testConfig) String(key string) (string, error) {
	return c.values[key], nil
}

func (a *testApp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.handler.ServeHTTP(w, r)
}

func (a *testApp) HandleCommand(args []string, stdout io.Writer, _ io.Writer) int {
	_ = args
	_, _ = io.WriteString(stdout, "stub")

	return 0
}

func (a *testApp) Config() *testConfig {
	return a.config
}

func newTestApp(_ *testing.T) *testApp {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, "Bedrock")
	})

	mux.HandleFunc("GET /up", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, "ok")
	})

	return &testApp{
		handler: mux,
		config: &testConfig{
			values: map[string]string{
				"app.name": "Bedrock Demo",
			},
		},
	}
}
