package database

import dbevents "github.com/bedrock/packages/database/events"

// Event type aliases re-exported from the events subpackage for convenience.
type (
	// Connection events.
	QueryExecuted         = dbevents.QueryExecuted
	TransactionBeginning  = dbevents.TransactionBeginning
	TransactionCommitting = dbevents.TransactionCommitting
	TransactionCommitted  = dbevents.TransactionCommitted
	TransactionRolledBack = dbevents.TransactionRolledBack
	ConnectionEstablished = dbevents.ConnectionEstablished
	StatementPrepared     = dbevents.StatementPrepared

	// Model lifecycle events.
	ModelCreating      = dbevents.Creating
	ModelCreated       = dbevents.Created
	ModelUpdating      = dbevents.Updating
	ModelUpdated       = dbevents.Updated
	ModelSaving        = dbevents.Saving
	ModelSaved         = dbevents.Saved
	ModelDeleting      = dbevents.Deleting
	ModelDeleted       = dbevents.Deleted
	ModelRestoring     = dbevents.Restoring
	ModelRestored      = dbevents.Restored
	ModelTrashed       = dbevents.Trashed
	ModelForceDeleting = dbevents.ForceDeleting
	ModelForceDeleted  = dbevents.ForceDeleted
	ModelReplicating   = dbevents.Replicating
	ModelRetrieved     = dbevents.Retrieved

	// Migration events.
	MigrationStarted    = dbevents.MigrationStarted
	MigrationEnded      = dbevents.MigrationEnded
	MigrationSkipped    = dbevents.MigrationSkipped
	MigrationsStarted   = dbevents.MigrationsStarted
	MigrationsEnded     = dbevents.MigrationsEnded
	NoPendingMigrations = dbevents.NoPendingMigrations
	MigrationsPruned    = dbevents.MigrationsPruned

	// Schema events.
	SchemaDumped = dbevents.SchemaDumped
	SchemaLoaded = dbevents.SchemaLoaded
)
