package log_test

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bedrock/packages/config"
	"github.com/bedrock/packages/log"
)

func newTestConfig(channels map[string]any) *config.Repository {
	return config.New(map[string]any{
		"logging": map[string]any{
			"default":  "single",
			"channels": channels,
		},
	})
}

func TestManagerDefaultChannel(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	cfg := newTestConfig(map[string]any{
		"single": map[string]any{
			"driver": "single",
			"path":   filepath.Join(dir, "test.log"),
			"level":  "debug",
		},
	})

	m := log.NewManager(cfg)

	if m.GetDefaultDriver() != "single" {
		t.Fatalf("expected default driver 'single', got %q", m.GetDefaultDriver())
	}
}

func TestManagerChannelCaching(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	cfg := newTestConfig(map[string]any{
		"single": map[string]any{
			"driver": "single",
			"path":   filepath.Join(dir, "test.log"),
			"level":  "debug",
		},
	})

	m := log.NewManager(cfg)

	ch1, err := m.Channel("single")
	if err != nil {
		t.Fatalf("Channel: %v", err)
	}

	ch2, err := m.Channel("single")
	if err != nil {
		t.Fatalf("Channel: %v", err)
	}

	if ch1 != ch2 {
		t.Fatal("expected same logger instance for same channel")
	}
}

func TestManagerSingleDriver(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "single.log")
	cfg := newTestConfig(map[string]any{
		"file": map[string]any{
			"driver": "single",
			"path":   path,
			"level":  "debug",
		},
	})

	m := log.NewManager(cfg)

	ch, err := m.Channel("file")
	if err != nil {
		t.Fatalf("Channel: %v", err)
	}

	ch.Info("single driver test")
	m.ForgetChannel("file")

	data, _ := os.ReadFile(path)
	if !strings.Contains(string(data), "single driver test") {
		t.Fatalf("expected file to contain log message, got %q", string(data))
	}
}

func TestManagerDailyDriver(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	cfg := newTestConfig(map[string]any{
		"daily": map[string]any{
			"driver": "daily",
			"path":   filepath.Join(dir, "daily.log"),
			"days":   7,
			"level":  "debug",
		},
	})

	m := log.NewManager(cfg)

	ch, err := m.Channel("daily")
	if err != nil {
		t.Fatalf("Channel: %v", err)
	}

	ch.Info("daily test")
	m.ForgetChannel("daily")

	matches, _ := filepath.Glob(filepath.Join(dir, "daily-*.log"))
	if len(matches) == 0 {
		t.Fatal("expected rotated file to exist")
	}
}

func TestManagerErrorlogDriver(t *testing.T) {
	t.Parallel()

	cfg := newTestConfig(map[string]any{
		"stderr": map[string]any{
			"driver": "errorlog",
			"level":  "debug",
		},
	})

	m := log.NewManager(cfg)

	ch, err := m.Channel("stderr")
	if err != nil {
		t.Fatalf("Channel: %v", err)
	}

	if ch == nil {
		t.Fatal("expected non-nil logger")
	}
}

func TestManagerNullDriver(t *testing.T) {
	t.Parallel()

	cfg := newTestConfig(map[string]any{
		"null": map[string]any{
			"driver": "null",
		},
	})

	m := log.NewManager(cfg)

	ch, err := m.Channel("null")
	if err != nil {
		t.Fatalf("Channel: %v", err)
	}

	ch.Info("null test")
}

func TestManagerStackDriver(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path1 := filepath.Join(dir, "ch1.log")
	path2 := filepath.Join(dir, "ch2.log")

	cfg := newTestConfig(map[string]any{
		"ch1": map[string]any{
			"driver": "single",
			"path":   path1,
			"level":  "debug",
		},
		"ch2": map[string]any{
			"driver": "single",
			"path":   path2,
			"level":  "debug",
		},
		"stack": map[string]any{
			"driver":   "stack",
			"channels": []any{"ch1", "ch2"},
			"level":    "debug",
		},
	})

	m := log.NewManager(cfg)

	ch, err := m.Channel("stack")
	if err != nil {
		t.Fatalf("Channel: %v", err)
	}

	ch.Info("stack test")
	m.ForgetChannel("ch1")
	m.ForgetChannel("ch2")

	data1, _ := os.ReadFile(path1)
	data2, _ := os.ReadFile(path2)

	if !strings.Contains(string(data1), "stack test") {
		t.Fatalf("expected ch1 to contain 'stack test', got %q", string(data1))
	}

	if !strings.Contains(string(data2), "stack test") {
		t.Fatalf("expected ch2 to contain 'stack test', got %q", string(data2))
	}
}

func TestManagerCustomDriver(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	cfg := newTestConfig(map[string]any{
		"custom": map[string]any{
			"driver": "my-driver",
			"level":  "debug",
		},
	})

	m := log.NewManager(cfg)
	m.Extend("my-driver", func(_ log.ChannelConfig) (log.Handler, error) {
		return log.NewStreamHandler(&buf, log.LevelDebug), nil
	})

	ch, err := m.Channel("custom")
	if err != nil {
		t.Fatalf("Channel: %v", err)
	}

	ch.Info("custom driver test")

	if !strings.Contains(buf.String(), "custom driver test") {
		t.Fatalf("expected custom handler to receive log, got %q", buf.String())
	}
}

func TestManagerCustomDriverOverride(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	cfg := newTestConfig(map[string]any{
		"override": map[string]any{
			"driver": "single",
			"level":  "debug",
		},
	})

	m := log.NewManager(cfg)
	m.Extend("single", func(_ log.ChannelConfig) (log.Handler, error) {
		return log.NewStreamHandler(&buf, log.LevelDebug), nil
	})

	ch, err := m.Channel("override")
	if err != nil {
		t.Fatalf("Channel: %v", err)
	}

	ch.Info("override test")

	if !strings.Contains(buf.String(), "override test") {
		t.Fatal("expected custom override to handle record")
	}
}

func TestManagerBuild(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	cfg := newTestConfig(map[string]any{})
	m := log.NewManager(cfg)

	logger, err := m.Build(log.ChannelConfig{
		Driver: log.DriverCustom,
		Via: func(_ log.ChannelConfig) (log.Handler, error) {
			return log.NewStreamHandler(&buf, log.LevelDebug), nil
		},
	})
	if err != nil {
		t.Fatalf("Build: %v", err)
	}

	logger.Info("build test")

	if !strings.Contains(buf.String(), "build test") {
		t.Fatal("expected built logger to write")
	}
}

func TestManagerBuildDoesNotCache(t *testing.T) {
	t.Parallel()

	cfg := newTestConfig(map[string]any{})
	m := log.NewManager(cfg)

	_, err := m.Build(log.ChannelConfig{
		Driver: log.DriverNull,
	})
	if err != nil {
		t.Fatalf("Build: %v", err)
	}

	channels := m.GetChannels()
	if len(channels) != 0 {
		t.Fatalf("expected no cached channels after Build, got %d", len(channels))
	}
}

func TestManagerShareContext(t *testing.T) {
	t.Parallel()

	spy := newSpyHandler(log.LevelDebug)
	cfg := newTestConfig(map[string]any{
		"spy": map[string]any{
			"driver": "my-spy",
			"level":  "debug",
		},
	})

	m := log.NewManager(cfg)
	m.Extend("my-spy", func(_ log.ChannelConfig) (log.Handler, error) {
		return spy, nil
	})

	m.ShareContext(map[string]any{"req_id": "xyz"})

	ch, _ := m.Channel("spy")
	ch.Info("shared context test")

	record := spy.LastRecord()
	if record.Context["req_id"] != "xyz" {
		t.Fatalf("expected req_id = xyz, got %v", record.Context["req_id"])
	}
}

func TestManagerSharedContextMerge(t *testing.T) {
	t.Parallel()

	cfg := newTestConfig(map[string]any{})
	m := log.NewManager(cfg)

	m.ShareContext(map[string]any{"a": 1})
	m.ShareContext(map[string]any{"b": 2})

	ctx := m.SharedContext()
	if ctx["a"] != 1 {
		t.Fatalf("expected a = 1, got %v", ctx["a"])
	}

	if ctx["b"] != 2 {
		t.Fatalf("expected b = 2, got %v", ctx["b"])
	}
}

func TestManagerFlushSharedContext(t *testing.T) {
	t.Parallel()

	cfg := newTestConfig(map[string]any{})
	m := log.NewManager(cfg)

	m.ShareContext(map[string]any{"a": 1})
	m.FlushSharedContext()

	ctx := m.SharedContext()
	if len(ctx) != 0 {
		t.Fatalf("expected empty shared context after flush, got %v", ctx)
	}
}

func TestManagerWithoutContext(t *testing.T) {
	t.Parallel()

	cfg := newTestConfig(map[string]any{})
	m := log.NewManager(cfg)

	m.ShareContext(map[string]any{"a": 1, "b": 2})
	m.WithoutContext("a")

	ctx := m.SharedContext()
	if _, ok := ctx["a"]; ok {
		t.Fatal("expected key 'a' to be removed")
	}

	if ctx["b"] != 2 {
		t.Fatalf("expected b = 2, got %v", ctx["b"])
	}
}

func TestManagerForgetChannel(t *testing.T) {
	t.Parallel()

	cfg := newTestConfig(map[string]any{
		"null": map[string]any{
			"driver": "null",
		},
	})

	m := log.NewManager(cfg)

	_, err := m.Channel("null")
	if err != nil {
		t.Fatalf("Channel: %v", err)
	}

	if len(m.GetChannels()) != 1 {
		t.Fatalf("expected 1 channel, got %d", len(m.GetChannels()))
	}

	m.ForgetChannel("null")

	if len(m.GetChannels()) != 0 {
		t.Fatalf("expected 0 channels after forget, got %d", len(m.GetChannels()))
	}
}

func TestManagerForgetChannelClosesHandler(t *testing.T) {
	t.Parallel()

	closed := false
	spy := &closableSpyHandler{
		spyHandler: newSpyHandler(log.LevelDebug),
		onClose:    func() { closed = true },
	}

	cfg := newTestConfig(map[string]any{
		"closable": map[string]any{
			"driver": "closable",
			"level":  "debug",
		},
	})

	m := log.NewManager(cfg)
	m.Extend("closable", func(_ log.ChannelConfig) (log.Handler, error) {
		return spy, nil
	})

	_, _ = m.Channel("closable")
	m.ForgetChannel("closable")

	if !closed {
		t.Fatal("expected handler Close to be called on ForgetChannel")
	}
}

type closableSpyHandler struct {
	*spyHandler
	onClose func()
}

func (h *closableSpyHandler) Close() error {
	h.onClose()
	return nil
}

func TestManagerTapCallback(t *testing.T) {
	t.Parallel()

	spy := newSpyHandler(log.LevelDebug)
	cfg := newTestConfig(map[string]any{
		"tapped": map[string]any{
			"driver": "spy",
			"level":  "debug",
		},
	})

	m := log.NewManager(cfg)
	m.Extend("spy", func(_ log.ChannelConfig) (log.Handler, error) {
		return spy, nil
	})

	tapped := false
	m.Tap(func(_ *log.Logger) {
		tapped = true
	})

	_, _ = m.Channel("tapped")

	if !tapped {
		t.Fatal("expected tap callback to be called")
	}
}

func TestManagerTapMultipleCallbacks(t *testing.T) {
	t.Parallel()

	spy := newSpyHandler(log.LevelDebug)
	cfg := newTestConfig(map[string]any{
		"tapped": map[string]any{
			"driver": "spy",
			"level":  "debug",
		},
	})

	m := log.NewManager(cfg)
	m.Extend("spy", func(_ log.ChannelConfig) (log.Handler, error) {
		return spy, nil
	})

	var order []int
	m.Tap(func(_ *log.Logger) { order = append(order, 1) })
	m.Tap(func(_ *log.Logger) { order = append(order, 2) })

	_, _ = m.Channel("tapped")

	if len(order) != 2 || order[0] != 1 || order[1] != 2 {
		t.Fatalf("expected taps in order [1, 2], got %v", order)
	}
}

func TestManagerEmergencyFallback(t *testing.T) {
	t.Parallel()

	cfg := newTestConfig(map[string]any{
		"broken": map[string]any{
			"driver": "broken",
			"level":  "debug",
		},
	})

	m := log.NewManager(cfg, log.WithDefaultChannel("broken"))
	m.Extend("broken", func(_ log.ChannelConfig) (log.Handler, error) {
		return nil, errors.New("driver creation failed")
	})

	// Should not panic; uses emergency logger.
	m.Info("fallback test")
}

func TestManagerUnsupportedDriver(t *testing.T) {
	t.Parallel()

	cfg := newTestConfig(map[string]any{
		"bad": map[string]any{
			"driver": "nonexistent",
			"level":  "debug",
		},
	})

	m := log.NewManager(cfg)

	_, err := m.Channel("bad")
	if err == nil {
		t.Fatal("expected error for unsupported driver")
	}

	if !strings.Contains(err.Error(), "unsupported driver") {
		t.Fatalf("expected 'unsupported driver' in error, got %q", err.Error())
	}
}

func TestManagerChannelNotConfigured(t *testing.T) {
	t.Parallel()

	cfg := newTestConfig(map[string]any{})
	m := log.NewManager(cfg)

	_, err := m.Channel("nonexistent")
	if err == nil {
		t.Fatal("expected error for non-configured channel")
	}

	if !errors.Is(err, log.ErrChannelNotFound) {
		t.Fatalf("expected ErrChannelNotFound, got %v", err)
	}
}

func TestManagerImplementsLoggerContract(t *testing.T) {
	t.Parallel()

	cfg := newTestConfig(map[string]any{})
	m := log.NewManager(cfg)

	// Compile-time check is in manager.go but we verify it here too.
	var _ interface {
		Info(string, ...map[string]any)
	} = m
}

func TestManagerDelegatesLevelMethods(t *testing.T) {
	t.Parallel()

	spy := newSpyHandler(log.LevelDebug)
	cfg := newTestConfig(map[string]any{
		"spy": map[string]any{
			"driver": "spy",
			"level":  "debug",
		},
	})

	m := log.NewManager(cfg, log.WithDefaultChannel("spy"))
	m.Extend("spy", func(_ log.ChannelConfig) (log.Handler, error) {
		return spy, nil
	})

	m.Debug("d")
	m.Info("i")
	m.Notice("n")
	m.Warning("w")
	m.Error("e")
	m.Critical("c")
	m.Alert("a")
	m.Emergency("em")
	m.Log(log.LevelInfo, "l")

	records := spy.Records()
	if len(records) != 9 {
		t.Fatalf("expected 9 records, got %d", len(records))
	}
}

func TestManagerSetDefaultDriver(t *testing.T) {
	t.Parallel()

	cfg := newTestConfig(map[string]any{
		"null": map[string]any{
			"driver": "null",
		},
	})

	m := log.NewManager(cfg)
	m.SetDefaultDriver("null")

	if m.GetDefaultDriver() != "null" {
		t.Fatalf("expected default driver 'null', got %q", m.GetDefaultDriver())
	}
}

func TestManagerDriverAlias(t *testing.T) {
	t.Parallel()

	cfg := newTestConfig(map[string]any{
		"null": map[string]any{
			"driver": "null",
		},
	})

	m := log.NewManager(cfg)

	ch, err := m.Driver("null")
	if err != nil {
		t.Fatalf("Driver: %v", err)
	}

	if ch == nil {
		t.Fatal("expected non-nil logger from Driver()")
	}
}

func TestManagerWithEventDispatcher(t *testing.T) {
	t.Parallel()

	spy := newSpyHandler(log.LevelDebug)
	dispatcher := newSpyDispatcher()

	cfg := newTestConfig(map[string]any{
		"spy": map[string]any{
			"driver": "spy",
			"level":  "debug",
		},
	})

	m := log.NewManager(cfg, log.WithEventDispatcher(dispatcher), log.WithDefaultChannel("spy"))
	m.Extend("spy", func(_ log.ChannelConfig) (log.Handler, error) {
		return spy, nil
	})

	m.Info("event test")

	events := dispatcher.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
}

func TestManagerStackMethod(t *testing.T) {
	t.Parallel()

	spy1 := newSpyHandler(log.LevelDebug)
	spy2 := newSpyHandler(log.LevelDebug)

	cfg := newTestConfig(map[string]any{
		"a": map[string]any{
			"driver": "spy-a",
			"level":  "debug",
		},
		"b": map[string]any{
			"driver": "spy-b",
			"level":  "debug",
		},
	})

	m := log.NewManager(cfg)
	m.Extend("spy-a", func(_ log.ChannelConfig) (log.Handler, error) { return spy1, nil })
	m.Extend("spy-b", func(_ log.ChannelConfig) (log.Handler, error) { return spy2, nil })

	logger, err := m.Stack([]string{"a", "b"})
	if err != nil {
		t.Fatalf("Stack: %v", err)
	}

	logger.Info("stack method test")

	if len(spy1.Records()) != 1 {
		t.Fatal("expected spy1 to receive record")
	}

	if len(spy2.Records()) != 1 {
		t.Fatal("expected spy2 to receive record")
	}
}

func TestManagerShareContextWithExistingChannels(t *testing.T) {
	t.Parallel()

	spy := newSpyHandler(log.LevelDebug)
	cfg := newTestConfig(map[string]any{
		"spy": map[string]any{
			"driver": "spy",
			"level":  "debug",
		},
	})

	m := log.NewManager(cfg)
	m.Extend("spy", func(_ log.ChannelConfig) (log.Handler, error) { return spy, nil })

	ch, _ := m.Channel("spy")

	m.ShareContext(map[string]any{"late_ctx": "added"})

	ch.Info("after share")

	record := spy.LastRecord()
	if record.Context["late_ctx"] != "added" {
		t.Fatalf("expected shared context on existing channel, got %v", record.Context)
	}
}

func TestManagerFormatterConfiguration(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	formatter := log.NewLineFormatter()
	formatter.IncludeContext = false

	cfg := newTestConfig(map[string]any{
		"formatted": map[string]any{
			"driver": "custom-fmt",
			"level":  "debug",
		},
	})

	m := log.NewManager(cfg)
	m.Extend("custom-fmt", func(_ log.ChannelConfig) (log.Handler, error) {
		h := log.NewStreamHandler(&buf, log.LevelDebug)
		h.SetFormatter(formatter)
		return h, nil
	})

	ch, _ := m.Channel("formatted")
	ch.Info("fmt test", map[string]any{"secret": "value"})

	if strings.Contains(buf.String(), "secret") {
		t.Fatal("expected formatter to exclude context")
	}
}

func TestManagerProcessorConfiguration(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	cfg := newTestConfig(map[string]any{
		"processed": map[string]any{
			"driver": "custom-proc",
			"level":  "debug",
		},
	})

	m := log.NewManager(cfg)
	m.Extend("custom-proc", func(_ log.ChannelConfig) (log.Handler, error) {
		h := log.NewStreamHandler(&buf, log.LevelDebug)
		h.AddProcessor(log.ProcessorFunc(func(r log.Record) log.Record {
			r.Extra["memory"] = "128MB"
			return r
		}))
		return h, nil
	})

	ch, _ := m.Channel("processed")
	ch.Info("proc test")

	if !strings.Contains(buf.String(), "128MB") {
		t.Fatal("expected processor to add memory info")
	}
}
