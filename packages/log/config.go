package log

import (
	"strings"

	"github.com/bedrock/packages/config"
)

// ChannelConfig holds the parsed configuration for a logging channel.
type ChannelConfig struct {
	Driver     Driver
	Path       string
	Level      Level
	Days       int
	Channels   []string
	Formatter  Formatter
	Processors []Processor
	With       map[string]any
	Via        func(config ChannelConfig) (Handler, error)
	Bubble     bool
	Permission int
	Tap        []func(*Logger)
}

// ParseChannelConfig extracts channel configuration from the given config
// repository using the key "logging.channels.<channel>".
func ParseChannelConfig(cfg *config.Repository, channel string) (ChannelConfig, error) {
	prefix := "logging.channels." + channel

	if !cfg.Has(prefix) {
		return ChannelConfig{}, ErrChannelNotFound
	}

	cc := ChannelConfig{
		Driver: Driver(stringOr(cfg, prefix+".driver", "single")),
		Path:   stringOr(cfg, prefix+".path", ""),
		Days:   intOr(cfg, prefix+".days", 7),
		Bubble: true,
		With:   make(map[string]any),
	}

	levelStr := stringOr(cfg, prefix+".level", "debug")

	level, err := ParseLevel(levelStr)
	if err != nil {
		cc.Level = LevelDebug
	} else {
		cc.Level = level
	}

	if channels := cfg.Get(prefix + ".channels"); channels != nil {
		switch v := channels.(type) {
		case []any:
			for _, c := range v {
				if s, ok := c.(string); ok {
					cc.Channels = append(cc.Channels, s)
				}
			}
		case string:
			for _, s := range strings.Split(v, ",") {
				s = strings.TrimSpace(s)
				if s != "" {
					cc.Channels = append(cc.Channels, s)
				}
			}
		}
	}

	return cc, nil
}

func stringOr(cfg *config.Repository, key string, fallback string) string {
	v, err := cfg.String(key, fallback)
	if err != nil {
		return fallback
	}

	return v
}

func intOr(cfg *config.Repository, key string, fallback int) int {
	v, err := cfg.Integer(key, fallback)
	if err != nil {
		return fallback
	}

	return v
}
