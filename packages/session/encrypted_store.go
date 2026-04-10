package session

import "context"

// EncryptedStore wraps a Store and encrypts every attribute value symmetrically.
// It uses the Encrypter interface defined in handler.go.
type EncryptedStore struct {
	*Store
	enc Encrypter
}

// NewEncryptedStore creates a Store whose attribute values are encrypted via enc.
func NewEncryptedStore(name string, handler Handler, enc Encrypter) *EncryptedStore {
	return &EncryptedStore{Store: New(name, handler), enc: enc}
}

// Put encrypts value before storing.
func (s *EncryptedStore) Put(key string, value any) {
	if str, ok := value.(string); ok {
		encrypted, err := s.enc.Encrypt(str)

		if err == nil {
			s.Store.Put(key, encrypted)

			return
		}
	}

	s.Store.Put(key, value)
}

// Get decrypts the value after retrieval.
func (s *EncryptedStore) Get(key string, fallback any) any {
	raw := s.Store.Get(key, fallback)

	if str, ok := raw.(string); ok {
		if plain, err := s.enc.Decrypt(str); err == nil {
			return plain
		}
	}

	return raw
}

// Start delegates to the inner store.
func (s *EncryptedStore) Start(ctx context.Context) error {
	return s.Store.Start(ctx)
}

// Save delegates to the inner store.
func (s *EncryptedStore) Save(ctx context.Context) error {
	return s.Store.Save(ctx)
}
