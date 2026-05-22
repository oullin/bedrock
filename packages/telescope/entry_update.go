package telescope

// TagsChange groups tag mutations applied to a stored entry.
type TagsChange struct {
	Add    []string
	Remove []string
}

// EntryUpdate describes a mutation to apply to a persisted Telescope entry.
type EntryUpdate struct {
	UUID    string
	Type    string
	Changes map[string]any
	Tags    TagsChange
}

// NewEntryUpdate creates an EntryUpdate for the given entry with the provided
func NewEntryUpdate(entryUUID, entryType string, changes map[string]any) *EntryUpdate {
	if changes == nil {
		changes = map[string]any{}
	}

	return &EntryUpdate{
		UUID:    entryUUID,
		Type:    entryType,
		Changes: changes,
	}
}

// Change merges additional field changes into the update and returns the
// update for chaining.
func (u *EntryUpdate) Change(changes map[string]any) *EntryUpdate {
	for k, v := range changes {
		u.Changes[k] = v
	}

	return u
}

// AddTags appends tags to be added on the stored entry (deduplicating) and
// returns the update for chaining.
func (u *EntryUpdate) AddTags(tags ...string) *EntryUpdate {
	existing := make(map[string]struct{}, len(u.Tags.Add))

	for _, t := range u.Tags.Add {
		existing[t] = struct{}{}
	}

	for _, t := range tags {
		if _, ok := existing[t]; !ok {
			u.Tags.Add = append(u.Tags.Add, t)
			existing[t] = struct{}{}
		}
	}

	return u
}

// RemoveTags appends tags to be removed from the stored entry (deduplicating)
// and returns the update for chaining.
func (u *EntryUpdate) RemoveTags(tags ...string) *EntryUpdate {
	existing := make(map[string]struct{}, len(u.Tags.Remove))

	for _, t := range u.Tags.Remove {
		existing[t] = struct{}{}
	}

	for _, t := range tags {
		if _, ok := existing[t]; !ok {
			u.Tags.Remove = append(u.Tags.Remove, t)
			existing[t] = struct{}{}
		}
	}

	return u
}
