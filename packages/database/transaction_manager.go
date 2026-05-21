package database

import "sync"

// TransactionManager manages after-commit and after-rollback callbacks
// at each transaction nesting level. It mirrors
// Ref: @bedrock/code-0206
type TransactionManager struct {
	mu        sync.Mutex
	callbacks map[int][]func()
}

// NewTransactionManager creates a new TransactionManager.
func NewTransactionManager() *TransactionManager {
	return &TransactionManager{
		callbacks: make(map[int][]func()),
	}
}

// AfterCommit registers a callback at the given transaction level.
func (m *TransactionManager) AfterCommit(level int, fn func()) {
	m.mu.Lock()

	defer m.mu.Unlock()

	m.callbacks[level] = append(m.callbacks[level], fn)
}

// Commit fires all callbacks registered at the given level and below.
func (m *TransactionManager) Commit(level int) {
	m.mu.Lock()

	var toRun []func()

	for l, fns := range m.callbacks {
		if l <= level {
			toRun = append(toRun, fns...)
			delete(m.callbacks, l)
		}
	}

	m.mu.Unlock()

	for _, fn := range toRun {
		fn()
	}
}

// Rollback discards all callbacks registered at the given level.
func (m *TransactionManager) Rollback(level int) {
	m.mu.Lock()

	defer m.mu.Unlock()

	delete(m.callbacks, level)
}
