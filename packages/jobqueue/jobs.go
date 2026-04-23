package jobqueue

import (
	"sort"
	"sync"
	"time"
)

// JobStatus describes where a queue job sits in the JobQueue lifecycle.
type JobStatus string

// JobPending indicates a job is waiting to be processed.

// JobReserved indicates a worker has reserved the job.

// JobCompleted indicates a job finished successfully.

// JobFailed indicates a job failed.

// JobRecord stores queue job metadata tracked by JobQueue.
type JobRecord struct {
	ID                  string
	Name                string
	Queue               string
	Type                string
	Tags                []string
	Status              JobStatus
	PushedAt            time.Time
	ReservedAt          time.Time
	CompletedAt         time.Time
	FailedAt            time.Time
	DelayUntil          time.Time
	FailedTagsExpireAt  time.Time
	CompletionStored    bool
	ReservedWasMigrated bool
}

// JobRepository stores queue job lifecycle records in memory.
type JobRepository struct {
	mu      sync.RWMutex
	now     func() time.Time
	jobs    map[string]JobRecord
	recent  []JobRecord
	failed  map[string]JobRecord
	tagJobs map[string]map[string]struct{}
}

const (
	JobPending JobStatus = "pending"

	JobReserved JobStatus = "reserved"

	JobCompleted JobStatus = "completed"

	JobFailed JobStatus = "failed"
)

// NewJobRepository creates an empty in-memory job repository.
func NewJobRepository(now func() time.Time) *JobRepository {
	if now == nil {
		now = time.Now
	}

	return &JobRepository{
		now:     now,
		jobs:    make(map[string]JobRecord),
		failed:  make(map[string]JobRecord),
		tagJobs: make(map[string]map[string]struct{}),
	}
}

// StorePending records a pending job.
func (r *JobRepository) StorePending(job JobRecord) {
	r.mu.Lock()

	defer r.mu.Unlock()

	if job.Status == "" {
		job.Status = JobPending
	}

	if job.PushedAt.IsZero() {
		job.PushedAt = r.now()
	}

	job = cloneJob(job)
	r.jobs[job.ID] = job
	r.indexTags(job)
}

// Pending returns pending jobs for a queue with cursor-style pagination.
func (r *JobRepository) Pending(queue string, offset, limit int) []JobRecord {
	r.mu.RLock()

	defer r.mu.RUnlock()

	jobs := make([]JobRecord, 0, len(r.jobs))
	now := r.now()

	for _, job := range r.jobs {
		if job.Queue != queue || job.Status != JobPending {
			continue
		}

		if !job.DelayUntil.IsZero() && job.DelayUntil.After(now) {
			continue
		}

		jobs = append(jobs, cloneJob(job))
	}

	sortJobs(jobs)

	return pageJobs(jobs, offset, limit)
}

// Recent returns recent completed jobs.
func (r *JobRepository) Recent(limit int) []JobRecord {
	r.mu.RLock()

	defer r.mu.RUnlock()

	jobs := make([]JobRecord, len(r.recent))

	for i, job := range r.recent {
		jobs[i] = cloneJob(job)
	}

	sortJobs(jobs)

	if limit > 0 && len(jobs) > limit {
		jobs = jobs[len(jobs)-limit:]
	}

	return jobs
}

// TrimRecent keeps only the most recent completed jobs.
func (r *JobRepository) TrimRecent(max int) {
	r.mu.Lock()

	defer r.mu.Unlock()

	if max < 0 {
		max = 0
	}

	sortJobs(r.recent)

	if len(r.recent) > max {
		r.recent = append([]JobRecord(nil), r.recent[len(r.recent)-max:]...)
	}
}

// MarkReserved marks a pending job as reserved.
func (r *JobRepository) MarkReserved(id string, at time.Time) bool {
	r.mu.Lock()

	defer r.mu.Unlock()

	job, ok := r.jobs[id]

	if !ok {
		return false
	}

	job.Status = JobReserved
	job.ReservedAt = at
	r.jobs[id] = job

	return true
}

// MigrateStaleReserved returns stale reserved jobs to pending.
func (r *JobRepository) MigrateStaleReserved(before time.Time) int {
	r.mu.Lock()

	defer r.mu.Unlock()

	migrated := 0

	for id, job := range r.jobs {
		if job.Status != JobReserved || job.ReservedAt.After(before) {
			continue
		}

		job.Status = JobPending
		job.ReservedWasMigrated = true
		r.jobs[id] = job
		migrated++
	}

	return migrated
}

// MarkComplete moves a job out of pending storage and optionally stores it as recent.
func (r *JobRepository) MarkComplete(id string, storeCompletion bool, at time.Time) bool {
	r.mu.Lock()

	defer r.mu.Unlock()

	job, ok := r.jobs[id]

	if !ok {
		return false
	}

	delete(r.jobs, id)
	r.unindexTags(job)

	if storeCompletion {
		job.Status = JobCompleted
		job.CompletedAt = at
		job.CompletionStored = true
		r.recent = append(r.recent, cloneJob(job))
	}

	return true
}

// MarkFailed moves a job into failed storage and retains its tags.
func (r *JobRepository) MarkFailed(id string, at time.Time, tagTTL time.Duration) bool {
	r.mu.Lock()

	defer r.mu.Unlock()

	job, ok := r.jobs[id]

	if !ok {
		return false
	}

	delete(r.jobs, id)
	r.unindexTags(job)

	job.Status = JobFailed
	job.FailedAt = at

	if tagTTL > 0 {
		job.FailedTagsExpireAt = at.Add(tagTTL)
	}

	r.failed[id] = cloneJob(job)

	return true
}

// FindFailed returns a failed job by id.
func (r *JobRepository) FindFailed(id string) (JobRecord, bool) {
	r.mu.RLock()

	defer r.mu.RUnlock()

	job, ok := r.failed[id]

	if !ok {
		return JobRecord{}, false
	}

	return cloneJob(job), true
}

// DeleteFailed deletes a failed job.
func (r *JobRepository) DeleteFailed(id string) bool {
	r.mu.Lock()

	defer r.mu.Unlock()

	if _, ok := r.failed[id]; !ok {
		return false
	}

	delete(r.failed, id)

	return true
}

// Release delays a job until the provided time.
func (r *JobRepository) Release(id string, until time.Time) bool {
	r.mu.Lock()

	defer r.mu.Unlock()

	job, ok := r.jobs[id]

	if !ok {
		return false
	}

	job.Status = JobPending
	job.DelayUntil = until
	r.jobs[id] = job

	return true
}

// MigrateReleased clears delays for jobs due at or before now.
func (r *JobRepository) MigrateReleased(now time.Time) int {
	r.mu.Lock()

	defer r.mu.Unlock()

	migrated := 0

	for id, job := range r.jobs {
		if job.DelayUntil.IsZero() || job.DelayUntil.After(now) {
			continue
		}

		job.DelayUntil = time.Time{}
		r.jobs[id] = job
		migrated++
	}

	return migrated
}

// PurgeQueue removes recent jobs for a queue.
func (r *JobRepository) PurgeQueue(queue string) int {
	r.mu.Lock()

	defer r.mu.Unlock()

	kept := r.recent[:0]
	removed := 0

	for _, job := range r.recent {
		if job.Queue == queue {
			removed++

			continue
		}

		kept = append(kept, job)
	}

	r.recent = kept

	return removed
}

// JobIDsForTag returns pending job ids indexed for a tag.
func (r *JobRepository) JobIDsForTag(tag string, offset, limit int) []string {
	r.mu.RLock()

	defer r.mu.RUnlock()

	tagged := r.tagJobs[tag]
	ids := make([]string, 0, len(tagged))

	for id := range tagged {
		ids = append(ids, id)
	}

	sort.Strings(ids)

	if offset >= len(ids) {
		return nil
	}

	end := len(ids)

	if limit > 0 && offset+limit < end {
		end = offset + limit
	}

	return append([]string(nil), ids[offset:end]...)
}

func (r *JobRepository) indexTags(job JobRecord) {
	for _, tag := range job.Tags {
		if r.tagJobs[tag] == nil {
			r.tagJobs[tag] = make(map[string]struct{})
		}

		r.tagJobs[tag][job.ID] = struct{}{}
	}
}

func (r *JobRepository) unindexTags(job JobRecord) {
	for _, tag := range job.Tags {
		delete(r.tagJobs[tag], job.ID)

		if len(r.tagJobs[tag]) == 0 {
			delete(r.tagJobs, tag)
		}
	}
}

func sortJobs(jobs []JobRecord) {
	sort.SliceStable(jobs, func(i, j int) bool {
		if jobs[i].PushedAt.Equal(jobs[j].PushedAt) {
			return jobs[i].ID < jobs[j].ID
		}

		return jobs[i].PushedAt.Before(jobs[j].PushedAt)
	})
}

func pageJobs(jobs []JobRecord, offset, limit int) []JobRecord {
	if offset >= len(jobs) {
		return nil
	}

	end := len(jobs)

	if limit > 0 && offset+limit < end {
		end = offset + limit
	}

	return append([]JobRecord(nil), jobs[offset:end]...)
}

func cloneJob(job JobRecord) JobRecord {
	tags := make([]string, len(job.Tags))
	copy(tags, job.Tags)
	job.Tags = tags

	return job
}
