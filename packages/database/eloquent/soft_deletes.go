package eloquent

import (
	"time"

	"github.com/bedrock/packages/database/query"
)

// SoftDeletes adds soft delete behaviour to a model. Models with this
// embedded struct will use a deleted_at column instead of hard deleting.
type SoftDeletes struct {
	deletedAtColumn string
	forceDeleting   bool
}

// InitSoftDeletes sets the default deleted_at column.
func (sd *SoftDeletes) InitSoftDeletes() {
	if sd.deletedAtColumn == "" {
		sd.deletedAtColumn = "deleted_at"
	}
}

// GetDeletedAtColumn returns the deleted_at column name.
func (sd *SoftDeletes) GetDeletedAtColumn() string {
	if sd.deletedAtColumn == "" {
		return "deleted_at"
	}
	return sd.deletedAtColumn
}

// SetDeletedAtColumn sets the deleted_at column name.
func (sd *SoftDeletes) SetDeletedAtColumn(column string) {
	sd.deletedAtColumn = column
}

// IsForceDeleting returns whether the model is being force-deleted.
func (sd *SoftDeletes) IsForceDeleting() bool {
	return sd.forceDeleting
}

// Trashed checks if a model instance has been soft-deleted.
func (sd *SoftDeletes) Trashed(attrs *HasAttributes) bool {
	deletedAt := attrs.GetAttribute(sd.GetDeletedAtColumn())
	return deletedAt != nil
}

// RunSoftDelete sets the deleted_at attribute on the model.
func (sd *SoftDeletes) RunSoftDelete(attrs *HasAttributes) {
	attrs.SetAttribute(sd.GetDeletedAtColumn(), time.Now().Format("2006-01-02 15:04:05"))
}

// Restore clears the deleted_at attribute on the model.
func (sd *SoftDeletes) Restore(attrs *HasAttributes) {
	attrs.SetAttribute(sd.GetDeletedAtColumn(), nil)
}

// SoftDeleteScope is a global scope that excludes soft-deleted models.
type SoftDeleteScope struct {
	column string
}

// NewSoftDeleteScope creates a new soft delete global scope.
func NewSoftDeleteScope(column string) *SoftDeleteScope {
	return &SoftDeleteScope{column: column}
}

// Apply adds the "deleted_at is null" constraint.
func (s *SoftDeleteScope) Apply(builder *query.Builder) {
	builder.WhereNull(s.column)
}

// WithTrashed returns a scope function that includes soft-deleted models.
func WithTrashed() ScopeFunc {
	return func(builder *query.Builder) {
		// Remove the soft delete where clause by rebuilding without it.
		// In practice, this means the EloquentBuilder should track and
		// skip the SoftDeleteScope when this scope is applied.
	}
}

// OnlyTrashed returns a scope function that only returns soft-deleted models.
func OnlyTrashed(column string) ScopeFunc {
	return func(builder *query.Builder) {
		builder.WhereNotNull(column)
	}
}
