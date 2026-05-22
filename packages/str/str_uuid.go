package str

import (
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
)

var (
	uuidMu       sync.Mutex
	uuidFactory  func() string // nil = use google/uuid
	uuidSequence []string      // consumed in order when set
	uuidFallback func() string // used when sequence is exhausted

	ulidMu       sync.Mutex
	ulidFactory  func() string
	ulidSequence []string
	ulidFallback func() string
)

// StrUuid generates a UUID version 4.
// If a factory or sequence is installed (via FreezeUuids/CreateUuidsUsing),
// it is used instead.
func StrUuid() string {
	uuidMu.Lock()

	defer uuidMu.Unlock()

	if uuidFactory != nil {
		return uuidFactory()
	}

	if len(uuidSequence) > 0 {
		val := uuidSequence[0]
		uuidSequence = uuidSequence[1:]

		return val
	}

	if uuidFallback != nil {
		return uuidFallback()
	}

	return uuid.New().String()
}

// StrOrderedUuid generates a time-sortable UUID version 7.
func StrOrderedUuid() string {
	uuidMu.Lock()

	defer uuidMu.Unlock()

	if uuidFactory != nil {
		return uuidFactory()
	}

	if len(uuidSequence) > 0 {
		val := uuidSequence[0]
		uuidSequence = uuidSequence[1:]

		return val
	}

	v7, err := uuid.NewV7()

	if err != nil {
		return uuid.New().String()
	}

	return v7.String()
}

// StrUuid7 generates a UUID version 7 with optional time.
func StrUuid7(t ...time.Time) string {
	return StrOrderedUuid()
}

// FreezeUuids installs a factory that returns a deterministic UUID string.
// Returns a cleanup function — call it (typically via defer) to restore normal generation.
// NOTE: do NOT combine with t.Parallel() at the top test level since this mutates global state.
func FreezeUuids(factory func() string) func() {
	uuidMu.Lock()
	prev := uuidFactory
	prevSeq := uuidSequence
	prevFallback := uuidFallback
	uuidFactory = factory
	uuidMu.Unlock()

	return func() {
		uuidMu.Lock()
		uuidFactory = prev
		uuidSequence = prevSeq
		uuidFallback = prevFallback
		uuidMu.Unlock()
	}
}

// CreateUuidsUsing installs a custom UUID factory.
func CreateUuidsUsing(factory func() string) {
	uuidMu.Lock()
	uuidFactory = factory
	uuidMu.Unlock()
}

// CreateUuidsUsingSequence sets a sequence of UUIDs to use for generation.
// Once the sequence is exhausted, whenMissing is called if provided.
func CreateUuidsUsingSequence(sequence []string, whenMissing ...func() string) func() {
	uuidMu.Lock()
	prevFactory := uuidFactory
	prevSeq := uuidSequence
	prevFallback := uuidFallback
	uuidSequence = make([]string, len(sequence))
	copy(uuidSequence, sequence)
	uuidFactory = nil

	if len(whenMissing) > 0 {
		uuidFallback = whenMissing[0]
	}

	uuidMu.Unlock()

	return func() {
		uuidMu.Lock()
		uuidFactory = prevFactory
		uuidSequence = prevSeq
		uuidFallback = prevFallback
		uuidMu.Unlock()
	}
}

// CreateUuidsNormally resets UUID generation to the default.
func CreateUuidsNormally() {
	uuidMu.Lock()
	uuidFactory = nil
	uuidSequence = nil
	uuidFallback = nil
	uuidMu.Unlock()
}

// StrUlid generates a ULID.
func StrUlid() string {
	ulidMu.Lock()

	defer ulidMu.Unlock()

	if ulidFactory != nil {
		return ulidFactory()
	}

	if len(ulidSequence) > 0 {
		val := ulidSequence[0]
		ulidSequence = ulidSequence[1:]

		return val
	}

	if ulidFallback != nil {
		return ulidFallback()
	}

	entropy := rand.New(rand.NewSource(time.Now().UnixNano())) //nolint:gosec
	ms := ulid.Timestamp(time.Now())
	id, err := ulid.New(ms, entropy)

	if err != nil {
		panic("support: failed to generate ULID: " + err.Error())
	}

	return id.String()
}

// FreezeUlids installs a factory that returns deterministic ULIDs.
// Returns a cleanup function.
func FreezeUlids(factory func() string) func() {
	ulidMu.Lock()
	prev := ulidFactory
	prevSeq := ulidSequence
	prevFallback := ulidFallback
	ulidFactory = factory
	ulidMu.Unlock()

	return func() {
		ulidMu.Lock()
		ulidFactory = prev
		ulidSequence = prevSeq
		ulidFallback = prevFallback
		ulidMu.Unlock()
	}
}

// CreateUlidsUsing installs a custom ULID factory.
func CreateUlidsUsing(factory func() string) {
	ulidMu.Lock()
	ulidFactory = factory
	ulidMu.Unlock()
}

// CreateUlidsUsingSequence sets a sequence of ULIDs to use for generation.
func CreateUlidsUsingSequence(sequence []string, whenMissing ...func() string) func() {
	ulidMu.Lock()
	prevFactory := ulidFactory
	prevSeq := ulidSequence
	prevFallback := ulidFallback
	ulidSequence = make([]string, len(sequence))
	copy(ulidSequence, sequence)
	ulidFactory = nil

	if len(whenMissing) > 0 {
		ulidFallback = whenMissing[0]
	}

	ulidMu.Unlock()

	return func() {
		ulidMu.Lock()
		ulidFactory = prevFactory
		ulidSequence = prevSeq
		ulidFallback = prevFallback
		ulidMu.Unlock()
	}
}

// CreateUlidsNormally resets ULID generation to the default.
func CreateUlidsNormally() {
	ulidMu.Lock()
	ulidFactory = nil
	ulidSequence = nil
	ulidFallback = nil
	ulidMu.Unlock()
}

// ResetFactoryState resets all string generation factories to defaults.
func ResetFactoryState() {
	CreateUuidsNormally()
	CreateUlidsNormally()
	CreateRandomStringsNormally()
}
