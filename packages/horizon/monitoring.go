package horizon

import (
	"sort"
	"sync"
)

// CompletedJob stores a completed job retained for monitored tags.
type CompletedJob struct {
	ID   string
	Tags []string
}

// MonitoringRepository tracks monitored tags and retained completed jobs.
type MonitoringRepository struct {
	mu        sync.RWMutex
	tags      map[string]struct{}
	completed map[string]CompletedJob
}

// NewMonitoringRepository creates an empty monitoring repository.
func NewMonitoringRepository() *MonitoringRepository {
	return &MonitoringRepository{
		tags:      make(map[string]struct{}),
		completed: make(map[string]CompletedJob),
	}
}

// Monitor starts monitoring the given tags.
func (r *MonitoringRepository) Monitor(tags []string) {
	r.mu.Lock()

	defer r.mu.Unlock()

	for _, tag := range tags {
		r.tags[tag] = struct{}{}
	}
}

// Monitoring returns the monitored tags in stable order.
func (r *MonitoringRepository) Monitoring() []string {
	r.mu.RLock()

	defer r.mu.RUnlock()

	tags := make([]string, 0, len(r.tags))

	for tag := range r.tags {
		tags = append(tags, tag)
	}

	sort.Strings(tags)

	return tags
}

// IsMonitoring reports whether any of the supplied tags are monitored.
func (r *MonitoringRepository) IsMonitoring(tags []string) bool {
	r.mu.RLock()

	defer r.mu.RUnlock()

	for _, tag := range tags {
		if _, ok := r.tags[tag]; ok {
			return true
		}
	}

	return false
}

// StopMonitoring removes monitored tags and prunes jobs no longer monitored.
func (r *MonitoringRepository) StopMonitoring(tags []string) {
	r.mu.Lock()

	defer r.mu.Unlock()

	for _, tag := range tags {
		delete(r.tags, tag)
	}

	for id, job := range r.completed {
		if !hasAnyTag(r.tags, job.Tags) {
			delete(r.completed, id)
		}
	}
}

// RecordCompletedJob stores a completed job only when one of its tags is monitored.
func (r *MonitoringRepository) RecordCompletedJob(job CompletedJob) bool {
	r.mu.Lock()

	defer r.mu.Unlock()

	if !hasAnyTag(r.tags, job.Tags) {
		return false
	}

	r.completed[job.ID] = cloneCompletedJob(job)

	return true
}

// CompletedJobs returns completed jobs, optionally filtered by tag.
func (r *MonitoringRepository) CompletedJobs(tag string) []CompletedJob {
	r.mu.RLock()

	defer r.mu.RUnlock()

	jobs := make([]CompletedJob, 0, len(r.completed))

	for _, job := range r.completed {
		if tag == "" || jobHasTag(job, tag) {
			jobs = append(jobs, cloneCompletedJob(job))
		}
	}

	sort.SliceStable(jobs, func(i, j int) bool {
		return jobs[i].ID < jobs[j].ID
	})

	return jobs
}

// CompletedJobsPage returns monitored completed jobs with offset/limit paging.
func (r *MonitoringRepository) CompletedJobsPage(tag string, offset, limit int) []CompletedJob {
	jobs := r.CompletedJobs(tag)

	if offset >= len(jobs) {
		return nil
	}

	end := len(jobs)

	if limit > 0 && offset+limit < end {
		end = offset + limit
	}

	return append([]CompletedJob(nil), jobs[offset:end]...)
}

func hasAnyTag(monitored map[string]struct{}, tags []string) bool {
	for _, tag := range tags {
		if _, ok := monitored[tag]; ok {
			return true
		}
	}

	return false
}

func jobHasTag(job CompletedJob, tag string) bool {
	for _, candidate := range job.Tags {
		if candidate == tag {
			return true
		}
	}

	return false
}

func cloneCompletedJob(job CompletedJob) CompletedJob {
	tags := make([]string, len(job.Tags))
	copy(tags, job.Tags)
	job.Tags = tags

	return job
}
