package inception

import (
	"context"
	"io"
)

// UpdatesProfilePhotos uploads a new profile photo.
type UpdatesProfilePhotos interface {
	Update(ctx context.Context, user HasTeams, photo io.Reader) (string, error)
}

// DeletesProfilePhotos removes the user's profile photo.
type DeletesProfilePhotos interface {
	Delete(ctx context.Context, user HasTeams) error
}

// DeletesUsers permanently deletes a user account.
type DeletesUsers interface {
	Delete(ctx context.Context, user HasTeams) error
}
