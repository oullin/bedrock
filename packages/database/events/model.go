package events

// Model lifecycle events. These are dispatched at various points during
// model persistence operations.

// Creating is dispatched before a new model is inserted.
type Creating struct{ Model any }

// Created is dispatched after a new model is inserted.
type Created struct{ Model any }

// Updating is dispatched before an existing model is updated.
type Updating struct{ Model any }

// Updated is dispatched after an existing model is updated.
type Updated struct{ Model any }

// Saving is dispatched before a model is saved (insert or update).
type Saving struct{ Model any }

// Saved is dispatched after a model is saved.
type Saved struct{ Model any }

// Deleting is dispatched before a model is deleted.
type Deleting struct{ Model any }

// Deleted is dispatched after a model is deleted.
type Deleted struct{ Model any }

// Restoring is dispatched before a soft-deleted model is restored.
type Restoring struct{ Model any }

// Restored is dispatched after a soft-deleted model is restored.
type Restored struct{ Model any }

// Trashed is dispatched after a model is soft-deleted.
type Trashed struct{ Model any }

// ForceDeleting is dispatched before a model is force-deleted.
type ForceDeleting struct{ Model any }

// ForceDeleted is dispatched after a model is force-deleted.
type ForceDeleted struct{ Model any }

// Replicating is dispatched when a model is being replicated.
type Replicating struct{ Model any }

// Retrieved is dispatched after a model is retrieved from the database.
type Retrieved struct{ Model any }
