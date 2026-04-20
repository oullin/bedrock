package horizon

import "errors"

// ErrNoSnapshots is returned when a repository has no captured state.
var ErrNoSnapshots = errors.New("horizon: no snapshots recorded")
