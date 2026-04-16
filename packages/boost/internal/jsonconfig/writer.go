// Package jsonconfig provides an atomic, concurrency-safe helper for
// upserting entries into JSON config files shared by multiple agents.
package jsonconfig

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

// WriteEntry performs a locked, atomic read-modify-write on a JSON config file.
// It merges serverConfig under root[configKey][serverKey]. If the file does not
// exist, it is created from skeleton (or an empty object when skeleton is nil).
//
// Returns true when the file was modified and false when the entry already
// existed (idempotent no-op). The caller decides how to map that boolean.
func WriteEntry(
	configPath, configKey, serverKey string,
	serverConfig map[string]any,
	skeleton map[string]any,
) (bool, error) {
	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		return false, fmt.Errorf("jsonconfig: mkdir %s: %w", filepath.Dir(configPath), err)
	}

	unlock, err := lockFile(configPath + ".lock")
	if err != nil {
		return false, fmt.Errorf("jsonconfig: lock %s: %w", configPath, err)
	}
	defer unlock()

	root, err := readRoot(configPath, skeleton)
	if err != nil {
		return false, err
	}

	servers := ensureMap(root, configKey)

	if _, exists := servers[serverKey]; exists {
		return false, nil
	}

	servers[serverKey] = serverConfig

	out, err := json.MarshalIndent(root, "", "    ")
	if err != nil {
		return false, fmt.Errorf("jsonconfig: marshal: %w", err)
	}

	if err := atomicWrite(configPath, append(out, '\n'), 0o644); err != nil {
		return false, fmt.Errorf("jsonconfig: write %s: %w", configPath, err)
	}

	return true, nil
}

// readRoot reads and unmarshals the JSON file, falling back to skeleton or an
// empty map when the file does not exist.
func readRoot(path string, skeleton map[string]any) (map[string]any, error) {
	data, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("jsonconfig: read %s: %w", path, err)
	}

	if len(data) > 0 {
		var root map[string]any
		if err := json.Unmarshal(data, &root); err != nil {
			return nil, fmt.Errorf("jsonconfig: unmarshal %s: %w", path, err)
		}
		return root, nil
	}

	if skeleton != nil {
		cp := make(map[string]any, len(skeleton))
		for k, v := range skeleton {
			cp[k] = v
		}
		return cp, nil
	}

	return make(map[string]any), nil
}

// ensureMap returns root[key] as a map, creating it if absent or mistyped.
func ensureMap(root map[string]any, key string) map[string]any {
	if m, ok := root[key].(map[string]any); ok {
		return m
	}
	m := make(map[string]any)
	root[key] = m
	return m
}

// atomicWrite writes data to path via a temp file in the same directory
// followed by an atomic rename.
func atomicWrite(path string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(path)

	tmp, err := os.CreateTemp(dir, ".tmp_*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		os.Remove(tmpName)
		return err
	}

	if err := tmp.Close(); err != nil {
		os.Remove(tmpName)
		return err
	}

	if err := os.Chmod(tmpName, perm); err != nil {
		os.Remove(tmpName)
		return err
	}

	return os.Rename(tmpName, path)
}

// lockFile acquires an exclusive advisory lock on the given path (creating it
// if necessary) and returns an unlock function. This matches the locking
// pattern used in packages/filesystem/lockable_file.go.
func lockFile(path string) (unlock func(), err error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return nil, err
	}

	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		f.Close()
		return nil, err
	}

	return func() {
		syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
		f.Close()
	}, nil
}
