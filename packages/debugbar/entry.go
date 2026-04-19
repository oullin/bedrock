package debugbar

import (
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
)

// EntryUser holds the authenticated user snapshot attached to an entry,
// mirroring the user() method on Upstream's IncomingEntry.
type EntryUser struct {
	ID    any    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// IncomingEntry represents a single telemetry entry being recorded before it
// is persisted. It mirrors Upstream's IncomingEntry class and provides a fluent
// builder interface.
type IncomingEntry struct {
	UUID       string
	BatchID    string
	Type       string
	FamilyHash string
	User       *EntryUser
	Content    map[string]any
	Tags       []string
	RecordedAt time.Time
	Hostname   string
}

// NewEntry creates a new IncomingEntry of the given type with a fresh ordered
// UUID, the current hostname, and the current timestamp.
func NewEntry(entryType string, content map[string]any) *IncomingEntry {
	hostname, _ := os.Hostname()

	e := &IncomingEntry{
		UUID:       uuid.New().String(),
		Type:       entryType,
		Content:    content,
		Tags:       []string{},
		RecordedAt: time.Now(),
		Hostname:   hostname,
	}

	return e
}

// WithBatchID sets the batch identifier and returns the entry for chaining.
func (e *IncomingEntry) WithBatchID(id string) *IncomingEntry {
	e.BatchID = id

	return e
}

// WithType sets the entry type and returns the entry for chaining.
func (e *IncomingEntry) WithType(t string) *IncomingEntry {
	e.Type = t

	return e
}

// WithFamilyHash sets the family hash used to group similar entries (e.g.,
// repeated identical SQL queries) and returns the entry for chaining.
func (e *IncomingEntry) WithFamilyHash(hash string) *IncomingEntry {
	e.FamilyHash = hash

	return e
}

// WithUser attaches authenticated user information to the entry, also adding a
// "user:{id}" tag. Returns the entry for chaining.
func (e *IncomingEntry) WithUser(user *EntryUser) *IncomingEntry {
	e.User = user

	if user != nil {
		e.Content["user"] = map[string]any{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		}
		e.AddTags(fmt.Sprintf("user:%v", user.ID))
	}

	return e
}

// WithTags replaces the entry's tag slice and returns the entry for chaining.
func (e *IncomingEntry) WithTags(tags []string) *IncomingEntry {
	e.Tags = tags

	return e
}

// AddTags appends tags (deduplicating) and returns the entry for chaining.
func (e *IncomingEntry) AddTags(tags ...string) *IncomingEntry {
	existing := make(map[string]struct{}, len(e.Tags))

	for _, t := range e.Tags {
		existing[t] = struct{}{}
	}

	for _, t := range tags {
		if _, ok := existing[t]; !ok {
			e.Tags = append(e.Tags, t)
			existing[t] = struct{}{}
		}
	}

	return e
}

// HasMonitoredTag reports whether any of the entry's tags appear in the
// monitored set, mirroring IncomingEntry::hasMonitoredTag().
func (e *IncomingEntry) HasMonitoredTag(monitored []string) bool {
	set := make(map[string]struct{}, len(monitored))

	for _, m := range monitored {
		set[m] = struct{}{}
	}

	for _, t := range e.Tags {
		if _, ok := set[t]; ok {
			return true
		}
	}

	return false
}

// IsRequest reports whether this is an HTTP request entry.
func (e *IncomingEntry) IsRequest() bool { return e.Type == EntryTypeRequest }

// IsFailedRequest reports whether this is a failed HTTP request (status ≥ 500).
func (e *IncomingEntry) IsFailedRequest() bool {
	if e.Type != EntryTypeRequest {
		return false
	}

	status, _ := e.Content["response_status"].(int)

	return status >= 500
}

// IsQuery reports whether this is a database query entry.
func (e *IncomingEntry) IsQuery() bool { return e.Type == EntryTypeQuery }

// IsSlowQuery reports whether this is a slow database query entry.
func (e *IncomingEntry) IsSlowQuery() bool {
	if e.Type != EntryTypeQuery {
		return false
	}

	slow, _ := e.Content["slow"].(bool)

	return slow
}

// IsEvent reports whether this is an application event entry.
func (e *IncomingEntry) IsEvent() bool { return e.Type == EntryTypeEvent }

// IsCache reports whether this is a cache operation entry.
func (e *IncomingEntry) IsCache() bool { return e.Type == EntryTypeCache }

// IsGate reports whether this is an authorization gate entry.
func (e *IncomingEntry) IsGate() bool { return e.Type == EntryTypeGate }

// IsFailedJob reports whether this is a failed queued job entry.
func (e *IncomingEntry) IsFailedJob() bool {
	if e.Type != EntryTypeJob {
		return false
	}

	status, _ := e.Content["status"].(string)

	return status == "failed"
}

// IsException reports whether this is an exception entry.
func (e *IncomingEntry) IsException() bool { return e.Type == EntryTypeException }

// IsLog reports whether this is a log message entry.
func (e *IncomingEntry) IsLog() bool { return e.Type == EntryTypeLog }

// IsScheduledTask reports whether this is a scheduled task entry.
func (e *IncomingEntry) IsScheduledTask() bool { return e.Type == EntryTypeScheduledTask }

// IsClientRequest reports whether this is an outbound HTTP client request entry.
func (e *IncomingEntry) IsClientRequest() bool { return e.Type == EntryTypeClientRequest }

// IsMail reports whether this is a mail message entry.
func (e *IncomingEntry) IsMail() bool { return e.Type == EntryTypeMail }

// IsNotification reports whether this is a notification entry.
func (e *IncomingEntry) IsNotification() bool { return e.Type == EntryTypeNotification }

// ToMap serialises the entry to a plain map for storage, mirroring
// IncomingEntry::toArray().
func (e *IncomingEntry) ToMap() map[string]any {
	m := map[string]any{
		"uuid":        e.UUID,
		"batch_id":    e.BatchID,
		"type":        e.Type,
		"family_hash": e.FamilyHash,
		"content":     e.Content,
		"tags":        e.Tags,
		"recorded_at": e.RecordedAt.UTC().Format(time.RFC3339Nano),
		"hostname":    e.Hostname,
	}

	if e.User != nil {
		m["user_id"] = e.User.ID
	}

	return m
}
