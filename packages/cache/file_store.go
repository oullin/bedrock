package cache

import (
	"context"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var _ Store = (*FileStore)(nil)

// FileStore caches values on the filesystem. Each key is stored as a gob file
// under a two-level directory tree derived from the SHA-256 hash of the key.
type FileStore struct {
	mu     sync.RWMutex
	dir    string
	prefix string
	clock  Clock
	perm   fs.FileMode
}

// fileEntry is the on-disk format for a cached value.
type fileEntry struct {
	ExpiresAt time.Time
	Value     any
}

// NewFileStore creates a FileStore that persists to dir.
func NewFileStore(dir string) *FileStore {
	return &FileStore{dir: dir, perm: 0o644}
}

// NewFileStoreWithOptions creates a FileStore with custom options.
func NewFileStoreWithOptions(dir, prefix string, perm fs.FileMode, clock Clock) *FileStore {
	return &FileStore{dir: dir, prefix: prefix, perm: perm, clock: clock}
}

func (s *FileStore) now() time.Time {
	if s.clock != nil {
		return s.clock.Now()
	}

	return time.Now()
}

func (s *FileStore) GetPrefix() string { return s.prefix }

func (s *FileStore) path(key string) string {
	h := sha256.Sum256([]byte(s.prefix + key))
	hex := hex.EncodeToString(h[:])

	return filepath.Join(s.dir, hex[:2], hex[2:4], hex)
}

func (s *FileStore) read(key string) (*fileEntry, error) {
	f, err := os.Open(s.path(key))
	if err != nil {
		return nil, fmt.Errorf("%w: %q", ErrNotFound, key)
	}
	defer f.Close()

	var entry fileEntry

	if err := gob.NewDecoder(f).Decode(&entry); err != nil {
		return nil, fmt.Errorf("%w: %q", ErrNotFound, key)
	}

	return &entry, nil
}

func (s *FileStore) write(key string, entry fileEntry) error {
	p := s.path(key)

	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}

	f, err := os.CreateTemp(filepath.Dir(p), "cache-")
	if err != nil {
		return err
	}

	if err := gob.NewEncoder(f).Encode(entry); err != nil {
		f.Close()
		os.Remove(f.Name())

		return err
	}

	f.Close()

	return os.Rename(f.Name(), p)
}

func (s *FileStore) Get(_ context.Context, key string) (any, error) {
	s.mu.RLock()
	entry, err := s.read(key)
	s.mu.RUnlock()

	if err != nil {
		return nil, err
	}

	if !entry.ExpiresAt.IsZero() && s.now().After(entry.ExpiresAt) {
		return nil, fmt.Errorf("%w: %q", ErrNotFound, key)
	}

	return entry.Value, nil
}

func (s *FileStore) GetMany(ctx context.Context, keys []string) (map[string]any, error) {
	out := make(map[string]any, len(keys))

	for _, key := range keys {
		if v, err := s.Get(ctx, key); err == nil {
			out[key] = v
		}
	}

	return out, nil
}

func (s *FileStore) Put(_ context.Context, key string, value any, ttl time.Duration) error {
	var exp time.Time
	if ttl > 0 {
		exp = s.now().Add(ttl)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	return s.write(key, fileEntry{ExpiresAt: exp, Value: value})
}

func (s *FileStore) PutMany(ctx context.Context, values map[string]any, ttl time.Duration) error {
	for k, v := range values {
		if err := s.Put(ctx, k, v, ttl); err != nil {
			return err
		}
	}

	return nil
}

func (s *FileStore) Add(ctx context.Context, key string, value any, ttl time.Duration) (bool, error) {
	if _, err := s.Get(ctx, key); err == nil {
		return false, nil
	}

	return true, s.Put(ctx, key, value, ttl)
}

func (s *FileStore) Forever(ctx context.Context, key string, value any) error {
	return s.Put(ctx, key, value, 0)
}

func (s *FileStore) Increment(ctx context.Context, key string, delta int64) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, err := s.read(key)

	var current int64
	var exp time.Time

	if err == nil {
		current, err = toInt64(entry.Value)
		if err != nil {
			return 0, fmt.Errorf("%w: key %q", ErrInvalidValue, key)
		}

		exp = entry.ExpiresAt
	}

	result := current + delta

	return result, s.write(key, fileEntry{ExpiresAt: exp, Value: result})
}

func (s *FileStore) Decrement(ctx context.Context, key string, delta int64) (int64, error) {
	return s.Increment(ctx, key, -delta)
}

func (s *FileStore) Touch(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, err := s.read(key)
	if err != nil {
		return false, nil
	}

	var exp time.Time
	if ttl > 0 {
		exp = s.now().Add(ttl)
	}

	return true, s.write(key, fileEntry{ExpiresAt: exp, Value: entry.Value})
}

func (s *FileStore) Forget(_ context.Context, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := os.Remove(s.path(key))
	if os.IsNotExist(err) {
		return nil
	}

	return err
}

// Flush removes all files under the store's directory.
func (s *FileStore) Flush(_ context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return os.RemoveAll(s.dir)
}
