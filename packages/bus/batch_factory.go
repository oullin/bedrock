package bus

import "time"

// BatchFactory constructs Batch instances from repository data.
type BatchFactory struct {
	dispatcher QueueingDispatcher
}

// NewBatchFactory creates a BatchFactory.
func NewBatchFactory(dispatcher QueueingDispatcher) *BatchFactory {
	return &BatchFactory{dispatcher: dispatcher}
}

// Make creates a Batch with all fields populated, including repo and dispatcher references.
func (f *BatchFactory) Make(
	repo BatchRepository,
	id, name string,
	totalJobs, pendingJobs, failedJobs int,
	failedJobIDs []string,
	options map[string]any,
	createdAt time.Time,
	cancelledAt, finishedAt *time.Time,
) *Batch {
	if failedJobIDs == nil {
		failedJobIDs = []string{}
	}

	if options == nil {
		options = make(map[string]any)
	}

	return &Batch{
		ID:           id,
		Name:         name,
		TotalJobs:    totalJobs,
		PendingJobs:  pendingJobs,
		FailedJobs:   failedJobs,
		FailedJobIDs: failedJobIDs,
		Options:      options,
		CreatedAt:    createdAt,
		CancelledAt:  cancelledAt,
		FinishedAt:   finishedAt,
		repo:         repo,
		dispatcher:   f.dispatcher,
	}
}
