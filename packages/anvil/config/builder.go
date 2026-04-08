package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

// Builder loads YAML configuration files and env overrides into a repository.
type Builder struct {
	dir    string
	appEnv string
	env    map[string]string
}

// NewBuilder creates a configuration builder rooted at dir.
func NewBuilder(dir string) *Builder {
	return &Builder{
		dir: filepath.Clean(dir),
	}
}

// WithAppEnv sets the application environment overlay.
func (b *Builder) WithAppEnv(appEnv string) *Builder {
	b.appEnv = strings.TrimSpace(appEnv)

	return b
}

// WithEnv overrides the process environment used for env var resolution.
func (b *Builder) WithEnv(env map[string]string) *Builder {
	if env == nil {
		b.env = nil

		return b
	}

	b.env = make(map[string]string, len(env))

	for key, value := range env {
		b.env[key] = value
	}

	return b
}

// Build loads the config repository.
func (b *Builder) Build(ctx context.Context) (*Repository, error) {
	_ = ctx

	items := map[string]any{}

	if err := b.loadDir(items, b.dir); err != nil {
		return nil, err
	}

	appEnv := b.appEnv

	if appEnv == "" {
		appEnv = strings.TrimSpace(b.lookupEnv("APP_ENV"))
	}

	if appEnv != "" {
		if err := b.loadDir(items, filepath.Join(b.dir, appEnv)); err != nil {
			return nil, err
		}
	}

	applyEnvOverrides(items, b.lookupEnv)

	return NewRepository(items), nil
}

func (b *Builder) loadDir(target map[string]any, dir string) error {
	info, err := os.Stat(dir)

	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return fmt.Errorf("stat config dir %s: %w", dir, err)
	}

	if !info.IsDir() {
		return fmt.Errorf("config path %s is not a directory", dir)
	}

	entries, err := os.ReadDir(dir)

	if err != nil {
		return fmt.Errorf("read config dir %s: %w", dir, err)
	}

	files := make([]string, 0, len(entries))

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()

		if strings.HasSuffix(name, ".yml") || strings.HasSuffix(name, ".yaml") {
			files = append(files, filepath.Join(dir, name))
		}
	}

	slices.Sort(files)

	for _, file := range files {
		if err := loadFile(target, file); err != nil {
			return err
		}
	}

	return nil
}

func (b *Builder) lookupEnv(key string) string {
	if b.env != nil {
		return b.env[key]
	}

	return os.Getenv(key)
}

func loadFile(target map[string]any, file string) error {
	v := viper.New()
	v.SetConfigFile(file)

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("read config file %s: %w", file, err)
	}

	namespace := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))
	target[namespace] = mergeValue(target[namespace], normalizeMap(v.AllSettings()))

	return nil
}

func normalizeMap(input map[string]any) map[string]any {
	output := make(map[string]any, len(input))

	for key, value := range input {
		output[key] = normalizeValue(value)
	}

	return output
}

func normalizeValue(value any) any {
	switch typed := value.(type) {
	case map[string]any:
		return normalizeMap(typed)
	case map[any]any:
		out := make(map[string]any, len(typed))

		for key, item := range typed {
			out[fmt.Sprint(key)] = normalizeValue(item)
		}

		return out
	case []any:
		out := make([]any, len(typed))

		for index, item := range typed {
			out[index] = normalizeValue(item)
		}

		return out
	default:
		return typed
	}
}

func mergeValue(current any, incoming any) any {
	currentMap, currentOK := current.(map[string]any)
	incomingMap, incomingOK := incoming.(map[string]any)

	if !currentOK || !incomingOK {
		return incoming
	}

	merged := make(map[string]any, len(currentMap))

	for key, value := range currentMap {
		merged[key] = value
	}

	for key, value := range incomingMap {
		merged[key] = mergeValue(merged[key], value)
	}

	return merged
}

func applyEnvOverrides(items map[string]any, lookup func(string) string) {
	for _, key := range flattenedKeys(items, "") {
		envKey := strings.ToUpper(strings.ReplaceAll(key, ".", "_"))
		envValue := lookup(envKey)

		if envValue == "" {
			continue
		}

		setByPath(items, key, coerceEnvValue(getByPath(items, key), envValue))
	}
}

func flattenedKeys(items map[string]any, prefix string) []string {
	keys := []string{}

	for key, value := range items {
		fullKey := key

		if prefix != "" {
			fullKey = prefix + "." + key
		}

		child, ok := value.(map[string]any)

		if !ok {
			keys = append(keys, fullKey)

			continue
		}

		keys = append(keys, flattenedKeys(child, fullKey)...)
	}

	slices.Sort(keys)

	return keys
}

func getByPath(items map[string]any, path string) any {
	segments := strings.Split(path, ".")
	current := any(items)

	for _, segment := range segments {
		mapped, ok := current.(map[string]any)

		if !ok {
			return nil
		}

		current = mapped[segment]
	}

	return current
}

func setByPath(items map[string]any, path string, value any) {
	segments := strings.Split(path, ".")
	current := items

	for _, segment := range segments[:len(segments)-1] {
		next := current[segment]
		child, ok := next.(map[string]any)

		if !ok {
			child = map[string]any{}
			current[segment] = child
		}

		current = child
	}

	current[segments[len(segments)-1]] = value
}

func coerceEnvValue(current any, raw string) any {
	switch current.(type) {
	case bool:
		parsed, err := strconv.ParseBool(strings.TrimSpace(raw))

		if err == nil {
			return parsed
		}
	case int:
		parsed, err := strconv.Atoi(strings.TrimSpace(raw))

		if err == nil {
			return parsed
		}
	case int64:
		parsed, err := strconv.ParseInt(strings.TrimSpace(raw), 10, 64)

		if err == nil {
			return parsed
		}
	case float64:
		parsed, err := strconv.ParseFloat(strings.TrimSpace(raw), 64)

		if err == nil {
			return parsed
		}
	case []any, []string:
		parts := strings.Split(raw, ",")
		values := make([]string, 0, len(parts))

		for _, part := range parts {
			values = append(values, strings.TrimSpace(part))
		}

		return values
	}

	return raw
}
