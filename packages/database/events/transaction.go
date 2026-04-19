package events

// TransactionBeginning is dispatched when a database transaction starts.
type TransactionBeginning struct {
	ConnectionName string
}

// TransactionCommitted is dispatched after a transaction has been committed.
type TransactionCommitted struct {
	ConnectionName string
}

// TransactionCommitting is dispatched before a transaction is committed.
type TransactionCommitting struct {
	ConnectionName string
}

// TransactionRolledBack is dispatched after a transaction has been rolled back.
type TransactionRolledBack struct {
	ConnectionName string
}
