package failed

import "time"

// Ref: @bedrock/code-0260
// no-op; Log returns an empty id, counters return zero, and Forget
// reports success. Useful for disabling failed-job tracking in tests.
type NullFailedJobProvider struct{}

// NewNullFailedJobProvider constructs a NullFailedJobProvider.
func NewNullFailedJobProvider() *NullFailedJobProvider { return &NullFailedJobProvider{} }

func (NullFailedJobProvider) Log(_, _, _ string, _ error) (string, error) { return "", nil }
func (NullFailedJobProvider) IDs(_ string) ([]string, error)              { return []string{}, nil }
func (NullFailedJobProvider) All() ([]FailedJob, error)                   { return nil, nil }
func (NullFailedJobProvider) Find(_ string) (*FailedJob, error)           { return nil, nil }
func (NullFailedJobProvider) Forget(_ string) (bool, error)               { return true, nil }
func (NullFailedJobProvider) Flush(_ int) error                           { return nil }
func (NullFailedJobProvider) Count(_, _ string) (int64, error)            { return 0, nil }
func (NullFailedJobProvider) Prune(_ time.Time) (int64, error)            { return 0, nil }
