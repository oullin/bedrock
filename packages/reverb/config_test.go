package reverb_test

import (
	"testing"

	"github.com/bedrock/packages/reverb"
)

func TestDefaultConfig_Host(t *testing.T) {
	t.Parallel()

	cfg := reverb.DefaultConfig()
	if cfg.Host != "0.0.0.0" {
		t.Errorf("expected Host %q, got %q", "0.0.0.0", cfg.Host)
	}
}

func TestDefaultConfig_Port(t *testing.T) {
	t.Parallel()

	cfg := reverb.DefaultConfig()
	if cfg.Port != 8080 {
		t.Errorf("expected Port %d, got %d", 8080, cfg.Port)
	}
}

func TestDefaultConfig_MaxRequestSize(t *testing.T) {
	t.Parallel()

	cfg := reverb.DefaultConfig()
	if cfg.MaxRequestSize != 10_000 {
		t.Errorf("expected MaxRequestSize %d, got %d", int64(10_000), cfg.MaxRequestSize)
	}
}

func TestAppConfig_Defaults(t *testing.T) {
	t.Parallel()

	app := reverb.NewApp(reverb.AppConfig{ID: "a"})

	if app.PingInterval() != 60 {
		t.Errorf("expected PingInterval %d, got %d", 60, app.PingInterval())
	}
	if app.ActivityTimeout() != 30 {
		t.Errorf("expected ActivityTimeout %d, got %d", 30, app.ActivityTimeout())
	}
	if app.ClientEventsMode() != "none" {
		t.Errorf("expected ClientEventsMode %q, got %q", "none", app.ClientEventsMode())
	}
}
