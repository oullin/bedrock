package jobqueue

import "errors"

// ErrNoSnapshots is returned when a repository has no captured state.
var ErrNoSnapshots = errors.New("jobqueue: no snapshots recorded")
