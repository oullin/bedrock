package debugbar

import (
	"crypto/md5"
	"fmt"
	"time"
)

// EntryResult represents a persisted DebugBar entry as returned by repository
// queries. It mirrors Upstream's EntryResult class and is suitable for JSON
// serialisation to the dashboard API.
type EntryResult struct {
	ID         string         `json:"id"`
	Sequence   int64          `json:"sequence"`
	BatchID    string         `json:"batch_id"`
	Type       string         `json:"type"`
	FamilyHash string         `json:"family_hash,omitempty"`
	Content    map[string]any `json:"content"`
	CreatedAt  time.Time      `json:"created_at"`
	Tags       []string       `json:"tags"`
	Avatar     string         `json:"avatar,omitempty"`
}

// GenerateAvatar sets the avatar URL using a Gravatar MD5 hash derived from
// the user email found in Content["user"]["email"], mirroring the PHP
// Avatar::url() helper. A custom avatar resolver can override this.
func (r *EntryResult) GenerateAvatar(resolver func(user map[string]any) string) {
	userRaw, ok := r.Content["user"]
	if !ok {
		return
	}

	user, ok := userRaw.(map[string]any)
	if !ok {
		return
	}

	if resolver != nil {
		r.Avatar = resolver(user)
		return
	}

	email, _ := user["email"].(string)
	if email == "" {
		return
	}

	hash := fmt.Sprintf("%x", md5.Sum([]byte(email)))
	r.Avatar = fmt.Sprintf("https://www.gravatar.com/avatar/%s?s=200&d=mm", hash)
}
