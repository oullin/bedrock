package mail

import "io"

// Attachment describes a file attached to an email message.
type Attachment struct {
	// Path is the filesystem path for path-based attachments.
	Path string
	// Data provides the attachment content for data-based attachments.
	Data func() (io.Reader, error)
	// Disk is the storage disk name for storage-based attachments.
	Disk string
	// Name is the display filename shown to the recipient.
	Name string
	// Mime is the MIME content type (e.g. "application/pdf").
	Mime string
	// Inline marks the attachment as an inline embed (Content-ID).
	Inline bool
}

// As sets the display filename and returns the attachment.
func (a *Attachment) As(name string) *Attachment {
	a.Name = name

	return a
}

// WithMime sets the MIME content type and returns the attachment.
func (a *Attachment) WithMime(mime string) *Attachment {
	a.Mime = mime

	return a
}

// IsEquivalent reports whether a matches other. Fields on other that are
// empty are treated as wildcards (i.e. they match any value on a).
func (a *Attachment) IsEquivalent(other *Attachment) bool {
	if other.Path != "" && a.Path != other.Path {
		return false
	}

	if other.Disk != "" && a.Disk != other.Disk {
		return false
	}

	if other.Name != "" && a.Name != other.Name {
		return false
	}

	if other.Mime != "" && a.Mime != other.Mime {
		return false
	}

	// At least one field must match positively.
	if other.Path == "" && other.Name == "" && other.Disk == "" {
		return false
	}

	return true
}

// FromPath creates a new path-based Attachment.
func FromPath(path string) *Attachment {
	return &Attachment{Path: path}
}

// FromData creates a new data-based Attachment with the given name.
func FromData(data func() (io.Reader, error), name string) *Attachment {
	return &Attachment{Data: data, Name: name}
}

// FromStorage creates a new storage-based Attachment on the default disk.
func FromStorage(path string) *Attachment {
	return &Attachment{Path: path}
}

// FromStorageDisk creates a new storage-based Attachment on the given disk.
func FromStorageDisk(disk, path string) *Attachment {
	return &Attachment{Disk: disk, Path: path}
}
