package session

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"
)

// Store manages HTTP session state including attributes, flash data, and
// CSRF tokens. It is safe for concurrent use.
type Store struct {
	mu         sync.RWMutex
	id         string
	name       string
	attributes map[string]any
	handler    Handler
	started    bool
}

// New creates a session store with a generated ID.
func New(name string, handler Handler) *Store {
	return &Store{
		name:       name,
		id:         generateID(),
		attributes: make(map[string]any),
		handler:    handler,
	}
}

// NewWithID creates a session store with a specific ID.
func NewWithID(name string, handler Handler, id string) *Store {
	s := &Store{
		name:       name,
		attributes: make(map[string]any),
		handler:    handler,
	}

	if isValidID(id) {
		s.id = id
	} else {
		s.id = generateID()
	}

	return s
}

// Start loads session data from the handler and ages flash data. Returns
// ErrAlreadyStarted if called twice without an intervening Invalidate.
func (s *Store) Start(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.started {
		return ErrAlreadyStarted
	}

	data, err := s.handler.Read(ctx, s.id)
	if err != nil {
		return fmt.Errorf("session: read: %w", err)
	}

	if data != "" {
		attrs, err := deserialize(data)
		if err != nil {
			return fmt.Errorf("session: deserialize: %w", err)
		}

		s.attributes = attrs
	}

	s.ageFlashData()
	s.started = true

	return nil
}

// Save serializes the session attributes and writes them to the handler.
func (s *Store) Save(ctx context.Context) error {
	s.mu.RLock()
	data, err := serialize(s.attributes)
	id := s.id
	s.mu.RUnlock()

	if err != nil {
		return fmt.Errorf("session: serialize: %w", err)
	}

	return s.handler.Write(ctx, id, data)
}

// IsStarted reports whether the session has been started.
func (s *Store) IsStarted() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.started
}

// GetID returns the session ID.
func (s *Store) GetID() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.id
}

// SetID sets the session ID. Returns ErrInvalidID if the ID is not 40 hex
// characters.
func (s *Store) SetID(id string) error {
	if !isValidID(id) {
		return fmt.Errorf("%w: %q", ErrInvalidID, id)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.id = id

	return nil
}

// GetName returns the session name.
func (s *Store) GetName() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.name
}

// SetName sets the session name.
func (s *Store) SetName(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.name = name
}

// Get retrieves an attribute value or returns the fallback.
func (s *Store) Get(key string, fallback any) any {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if v, ok := s.attributes[key]; ok {
		return v
	}

	return fallback
}

// Put stores an attribute value.
func (s *Store) Put(key string, value any) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.attributes[key] = value
}

// Has reports whether a non-nil value exists for the key.
func (s *Store) Has(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	v, ok := s.attributes[key]

	return ok && v != nil
}

// Exists reports whether the key exists, even if nil.
func (s *Store) Exists(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.attributes[key]

	return ok
}

// Missing reports whether the key does not exist.
func (s *Store) Missing(key string) bool {
	return !s.Exists(key)
}

// Pull retrieves and removes a value.
func (s *Store) Pull(key string, fallback any) any {
	s.mu.Lock()
	defer s.mu.Unlock()

	v, ok := s.attributes[key]
	if !ok {
		return fallback
	}

	delete(s.attributes, key)

	return v
}

// Push appends a value to a slice attribute. If the key does not exist, a
// new slice is created.
func (s *Store) Push(key string, value any) {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, ok := s.attributes[key]
	if !ok {
		s.attributes[key] = []any{value}
		return
	}

	if slice, ok := existing.([]any); ok {
		s.attributes[key] = append(slice, value)
	} else {
		s.attributes[key] = []any{existing, value}
	}
}

// All returns a shallow copy of all attributes.
func (s *Store) All() map[string]any {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make(map[string]any, len(s.attributes))
	for k, v := range s.attributes {
		result[k] = v
	}

	return result
}

// Forget removes one or more keys.
func (s *Store) Forget(keys ...string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, key := range keys {
		delete(s.attributes, key)
	}
}

// Flush removes all attributes.
func (s *Store) Flush() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.attributes = make(map[string]any)
}

// Flash stores a value available only for the next request.
func (s *Store) Flash(key string, value any) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.attributes[key] = value
	s.pushFlashKey(key)
	s.removeFromOldFlash(key)
}

// Reflash keeps all flash data for an additional request.
func (s *Store) Reflash() {
	s.mu.Lock()
	defer s.mu.Unlock()

	old := s.getFlashOld()

	s.setFlashNew(append(s.getFlashNew(), old...))
	s.setFlashOld(nil)
}

// Keep keeps specific flash keys for an additional request.
func (s *Store) Keep(keys ...string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, key := range keys {
		s.pushFlashKey(key)
		s.removeFromOldFlash(key)
	}
}

// Token returns the CSRF token, generating one if it does not exist.
func (s *Store) Token() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	if token, ok := s.attributes["_token"].(string); ok && token != "" {
		return token
	}

	token := generateToken()
	s.attributes["_token"] = token

	return token
}

// RegenerateToken generates a new CSRF token.
func (s *Store) RegenerateToken() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.attributes["_token"] = generateToken()
}

// Regenerate generates a new session ID, optionally destroying old data.
func (s *Store) Regenerate(ctx context.Context, destroy bool) error {
	if destroy {
		s.mu.RLock()
		oldID := s.id
		s.mu.RUnlock()

		if err := s.handler.Destroy(ctx, oldID); err != nil {
			return fmt.Errorf("session: destroy old: %w", err)
		}
	}

	s.mu.Lock()
	s.id = generateID()
	s.mu.Unlock()

	return nil
}

// Invalidate flushes all data and regenerates the session ID.
func (s *Store) Invalidate(ctx context.Context) error {
	s.Flush()

	s.mu.Lock()
	s.started = false
	s.mu.Unlock()

	return s.Regenerate(ctx, true)
}

// PreviousURL returns the previously visited URL.
func (s *Store) PreviousURL() string {
	v := s.Get("_previous_url", "")

	if url, ok := v.(string); ok {
		return url
	}

	return ""
}

// SetPreviousURL stores the previously visited URL.
func (s *Store) SetPreviousURL(url string) {
	s.Put("_previous_url", url)
}

// --- flash helpers (caller must hold lock) ---

func (s *Store) getFlashMeta() map[string]any {
	raw, ok := s.attributes["_flash"]
	if !ok {
		return nil
	}

	if m, ok := raw.(map[string]any); ok {
		return m
	}

	return nil
}

func (s *Store) ensureFlashMeta() map[string]any {
	m := s.getFlashMeta()
	if m == nil {
		m = map[string]any{
			"old": []any{},
			"new": []any{},
		}
		s.attributes["_flash"] = m
	}

	return m
}

func (s *Store) getFlashOld() []string {
	m := s.getFlashMeta()
	if m == nil {
		return nil
	}

	return toStringSlice(m["old"])
}

func (s *Store) getFlashNew() []string {
	m := s.getFlashMeta()
	if m == nil {
		return nil
	}

	return toStringSlice(m["new"])
}

func (s *Store) setFlashOld(keys []string) {
	m := s.ensureFlashMeta()
	m["old"] = toAnySlice(keys)
}

func (s *Store) setFlashNew(keys []string) {
	m := s.ensureFlashMeta()
	m["new"] = toAnySlice(keys)
}

func (s *Store) pushFlashKey(key string) {
	newKeys := s.getFlashNew()

	for _, k := range newKeys {
		if k == key {
			return
		}
	}

	s.setFlashNew(append(newKeys, key))
}

func (s *Store) removeFromOldFlash(key string) {
	old := s.getFlashOld()

	filtered := make([]string, 0, len(old))
	for _, k := range old {
		if k != key {
			filtered = append(filtered, k)
		}
	}

	s.setFlashOld(filtered)
}

// ageFlashData removes old flash keys and promotes new to old. Caller
// must hold the write lock.
func (s *Store) ageFlashData() {
	old := s.getFlashOld()

	for _, key := range old {
		delete(s.attributes, key)
	}

	s.setFlashOld(s.getFlashNew())
	s.setFlashNew(nil)
}

// --- ID and token generation ---

func generateID() string {
	b := make([]byte, 20)
	_, _ = rand.Read(b)

	return hex.EncodeToString(b)
}

func generateToken() string {
	b := make([]byte, 20)
	_, _ = rand.Read(b)

	return hex.EncodeToString(b)
}

func isValidID(id string) bool {
	if len(id) != 40 {
		return false
	}

	for _, c := range id {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}

	return true
}

// --- serialization ---

func serialize(attrs map[string]any) (string, error) {
	b, err := json.Marshal(attrs)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func deserialize(data string) (map[string]any, error) {
	var attrs map[string]any
	if err := json.Unmarshal([]byte(data), &attrs); err != nil {
		return nil, err
	}

	return attrs, nil
}

// --- slice conversion helpers ---

func toStringSlice(v any) []string {
	if v == nil {
		return nil
	}

	switch sl := v.(type) {
	case []string:
		return sl
	case []any:
		result := make([]string, 0, len(sl))
		for _, item := range sl {
			if s, ok := item.(string); ok {
				result = append(result, s)
			}
		}
		return result
	}

	return nil
}

func toAnySlice(ss []string) []any {
	result := make([]any, len(ss))
	for i, s := range ss {
		result[i] = s
	}

	return result
}
