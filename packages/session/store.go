package session

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
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

// NewWithID creates a session store with a specific ID, validating it first.
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

// SetID sets the session ID. Returns ErrInvalidID if the ID is not valid.
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

// Push appends a value to a slice attribute.
func (s *Store) Push(key string, value any) {
	s.mu.Lock()

	defer s.mu.Unlock()

	existing, ok := s.attributes[key]

	if !ok {
		s.attributes[key] = []any{value}

		return
	}

	if sl, ok := existing.([]any); ok {
		s.attributes[key] = append(sl, value)
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

// Only returns a map containing only the specified keys.
func (s *Store) Only(keys ...string) map[string]any {
	s.mu.RLock()

	defer s.mu.RUnlock()

	result := make(map[string]any, len(keys))

	for _, key := range keys {
		if v, ok := s.attributes[key]; ok {
			result[key] = v
		}
	}

	return result
}

// Except returns all attributes except the specified keys.
func (s *Store) Except(keys ...string) map[string]any {
	s.mu.RLock()

	defer s.mu.RUnlock()

	exclude := make(map[string]bool, len(keys))

	for _, k := range keys {
		exclude[k] = true
	}

	result := make(map[string]any, len(s.attributes))

	for k, v := range s.attributes {
		if !exclude[k] {
			result[k] = v
		}
	}

	return result
}

// Replace merges the given key-value pairs into the session.
func (s *Store) Replace(values map[string]any) {
	s.mu.Lock()

	defer s.mu.Unlock()

	for k, v := range values {
		s.attributes[k] = v
	}
}

// Remove retrieves and removes a value (nil fallback).
func (s *Store) Remove(key string) any {
	return s.Pull(key, nil)
}

// Increment increments a numeric session value by the given amount.
func (s *Store) Increment(key string, amount int64) int64 {
	s.mu.Lock()

	defer s.mu.Unlock()

	var current int64

	if v, ok := s.attributes[key]; ok {
		current = toSessionInt64(v)
	}

	result := current + amount
	s.attributes[key] = result

	return result
}

// Decrement decrements a numeric session value.
func (s *Store) Decrement(key string, amount int64) int64 {
	return s.Increment(key, -amount)
}

// Remember retrieves a value or stores the result of the callback.
func (s *Store) Remember(key string, callback func() any) any {
	s.mu.Lock()

	defer s.mu.Unlock()

	if v, ok := s.attributes[key]; ok {
		return v
	}

	value := callback()
	s.attributes[key] = value

	return value
}

// HasAny reports whether any of the given keys exist with non-nil values.
func (s *Store) HasAny(keys ...string) bool {
	s.mu.RLock()

	defer s.mu.RUnlock()

	for _, key := range keys {
		if v, ok := s.attributes[key]; ok && v != nil {
			return true
		}
	}

	return false
}

// Flash stores a value available only for the next request.
func (s *Store) Flash(key string, value any) {
	s.mu.Lock()

	defer s.mu.Unlock()

	s.attributes[key] = value
	s.pushFlashKey(key)
	s.removeFromOldFlash(key)
}

// Now stores a value for the current request only (expires this request).
func (s *Store) Now(key string, value any) {
	s.mu.Lock()

	defer s.mu.Unlock()

	s.attributes[key] = value
	old := s.getFlashOld()
	s.setFlashOld(append(old, key))
}

// FlashInput stores input data as "old input" for the next request.
func (s *Store) FlashInput(values map[string]any) {
	s.Flash("_old_input", values)
}

// GetOldInput retrieves a previously flashed input value.
func (s *Store) GetOldInput(key string, fallback any) any {
	s.mu.RLock()

	defer s.mu.RUnlock()

	raw, ok := s.attributes["_old_input"]

	if !ok {
		return fallback
	}

	m, ok := raw.(map[string]any)

	if !ok {
		return fallback
	}

	if key == "" {
		return m
	}

	if v, ok := m[key]; ok {
		return v
	}

	return fallback
}

// HasOldInput reports whether old input data exists for the given key.
func (s *Store) HasOldInput(key string) bool {
	s.mu.RLock()

	defer s.mu.RUnlock()

	raw, ok := s.attributes["_old_input"]

	if !ok {
		return false
	}

	m, ok := raw.(map[string]any)

	if !ok {
		return false
	}

	if key == "" {
		return len(m) > 0
	}

	_, exists := m[key]

	return exists
}

// Reflash keeps all current flash data for an additional request.
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

// Token returns the CSRF token, generating one if absent.
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

// Regenerate generates a new session ID, optionally destroying the old data.
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

// Migrate generates a new session ID (upstream alias for Regenerate).
func (s *Store) Migrate(ctx context.Context, destroy bool) error {
	return s.Regenerate(ctx, destroy)
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
	if url, ok := s.Get("_previous_url", "").(string); ok {
		return url
	}

	return ""
}

// SetPreviousURL stores the previously visited URL.
func (s *Store) SetPreviousURL(url string) {
	s.Put("_previous_url", url)
}

// PasswordConfirmed marks that the user has recently confirmed their password.
func (s *Store) PasswordConfirmed() {
	s.Put("auth.password_confirmed_at", nowUnix())
}

// PasswordConfirmedAt returns the Unix timestamp when the password was last confirmed.
func (s *Store) PasswordConfirmedAt() int64 {
	v := s.Get("auth.password_confirmed_at", int64(0))

	return toSessionInt64(v)
}

// GetHandler returns the underlying session handler.
func (s *Store) GetHandler() Handler {
	s.mu.RLock()

	defer s.mu.RUnlock()

	return s.handler
}

// SetHandler replaces the session handler and returns the previous one.
func (s *Store) SetHandler(handler Handler) Handler {
	s.mu.Lock()

	defer s.mu.Unlock()

	old := s.handler
	s.handler = handler

	return old
}

// HandlerNeedsRequest reports whether the handler implements RequestAware.
func (s *Store) HandlerNeedsRequest() bool {
	s.mu.RLock()

	defer s.mu.RUnlock()

	_, ok := s.handler.(RequestAware)

	return ok
}

// SetRequestOnHandler forwards the HTTP request to the handler if it
// implements RequestAware.
func (s *Store) SetRequestOnHandler(r *http.Request) {
	s.mu.RLock()

	defer s.mu.RUnlock()

	if ra, ok := s.handler.(RequestAware); ok {
		ra.SetRequest(r)
	}
}

// SetExists forwards the existence flag to the handler if it implements
// ExistenceAware.
func (s *Store) SetExists(exists bool) {
	s.mu.RLock()

	defer s.mu.RUnlock()

	if ea, ok := s.handler.(ExistenceAware); ok {
		ea.SetExists(exists)
	}
}

// IsValidID reports whether id is a valid session ID (40-char hex string).
func IsValidID(id string) bool {
	return isValidID(id)
}

// ID returns the session ID (alias for GetID).
func (s *Store) ID() string {
	return s.GetID()
}

// HasPreviousURL reports whether a previous URL has been stored.
func (s *Store) HasPreviousURL() bool {
	return s.PreviousURL() != ""
}

// PreviousRoute returns the previously matched route name.
func (s *Store) PreviousRoute() string {
	if route, ok := s.Get("_previous_route", "").(string); ok {
		return route
	}

	return ""
}

// SetPreviousRoute stores the previously matched route name.
func (s *Store) SetPreviousRoute(route string) {
	s.Put("_previous_route", route)
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
		m = map[string]any{"old": []any{}, "new": []any{}}
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

// ageFlashData removes old flash keys and promotes new→old. Caller holds write lock.
func (s *Store) ageFlashData() {
	old := s.getFlashOld()

	for _, key := range old {
		delete(s.attributes, key)
	}

	s.setFlashOld(s.getFlashNew())
	s.setFlashNew(nil)
}

// --- helpers ---

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

func toSessionInt64(v any) int64 {
	switch n := v.(type) {
	case int:
		return int64(n)
	case int64:
		return n
	case float64:
		return int64(n)
	case int32:
		return int64(n)
	default:
		return 0
	}
}

func toAnySlice(ss []string) []any {
	result := make([]any, len(ss))

	for i, s := range ss {
		result[i] = s
	}

	return result
}

func nowUnix() int64 {
	return time.Now().Unix()
}
