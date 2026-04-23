package jobqueue

import (
	"sort"
	"sync"
	"time"
)

// JobMeasurement records one completed job sample.
type JobMeasurement struct {
	Job         string
	Queue       string
	Runtime     time.Duration
	CompletedAt time.Time
}

// MetricSummary exposes throughput and average runtime for a metric bucket.
type MetricSummary struct {
	Throughput     int
	AverageRuntime time.Duration
}

// MetricsSnapshot stores a point-in-time performance snapshot.
type MetricsSnapshot struct {
	RecordedAt      time.Time
	JobsProcessed   int
	JobThroughput   map[string]int
	QueueThroughput map[string]int
}

// MetricsRepository stores in-memory throughput and runtime measurements.
type MetricsRepository struct {
	mu        sync.RWMutex
	total     metricBucket
	jobs      map[string]metricBucket
	queues    map[string]metricBucket
	snapshots []MetricsSnapshot
}

type metricBucket struct {
	count   int
	runtime time.Duration
}

// NewMetricsRepository creates an empty metrics repository.
func NewMetricsRepository() *MetricsRepository {
	return &MetricsRepository{
		jobs:   make(map[string]metricBucket),
		queues: make(map[string]metricBucket),
	}
}

// RecordJob records one completed job for aggregate metrics.
func (r *MetricsRepository) RecordJob(measurement JobMeasurement) {
	r.mu.Lock()

	defer r.mu.Unlock()

	r.total = r.total.add(measurement.Runtime)

	if measurement.Job != "" {
		r.jobs[measurement.Job] = r.jobs[measurement.Job].add(measurement.Runtime)
	}

	if measurement.Queue != "" {
		r.queues[measurement.Queue] = r.queues[measurement.Queue].add(measurement.Runtime)
	}
}

// Total returns aggregate throughput and runtime.
func (r *MetricsRepository) Total() MetricSummary {
	r.mu.RLock()

	defer r.mu.RUnlock()

	return r.total.summary()
}

// Job returns metrics for a job class.
func (r *MetricsRepository) Job(name string) MetricSummary {
	r.mu.RLock()

	defer r.mu.RUnlock()

	return r.jobs[name].summary()
}

// Queue returns metrics for a queue.
func (r *MetricsRepository) Queue(name string) MetricSummary {
	r.mu.RLock()

	defer r.mu.RUnlock()

	return r.queues[name].summary()
}

// Jobs returns all job classes that have metric samples.
func (r *MetricsRepository) Jobs() []string {
	r.mu.RLock()

	defer r.mu.RUnlock()

	jobs := make([]string, 0, len(r.jobs))

	for job := range r.jobs {
		jobs = append(jobs, job)
	}

	sort.Strings(jobs)

	return jobs
}

// SnapshotPerformance stores and returns the current aggregate counters.
func (r *MetricsRepository) SnapshotPerformance(recordedAt time.Time) MetricsSnapshot {
	r.mu.Lock()

	defer r.mu.Unlock()

	snapshot := MetricsSnapshot{
		RecordedAt:      recordedAt,
		JobsProcessed:   r.total.count,
		JobThroughput:   bucketCounts(r.jobs),
		QueueThroughput: bucketCounts(r.queues),
	}

	r.snapshots = append(r.snapshots, snapshot)

	if len(r.snapshots) > 24 {
		r.snapshots = append([]MetricsSnapshot(nil), r.snapshots[len(r.snapshots)-24:]...)
	}

	return cloneMetricsSnapshot(snapshot)
}

// Snapshots returns retained performance snapshots.
func (r *MetricsRepository) Snapshots() []MetricsSnapshot {
	r.mu.RLock()

	defer r.mu.RUnlock()

	snapshots := make([]MetricsSnapshot, len(r.snapshots))

	for i, snapshot := range r.snapshots {
		snapshots[i] = cloneMetricsSnapshot(snapshot)
	}

	return snapshots
}

// JobsProcessedPerMinuteSince returns the throughput since a prior snapshot.
func (r *MetricsRepository) JobsProcessedPerMinuteSince(snapshot MetricsSnapshot, now time.Time) float64 {
	r.mu.RLock()

	defer r.mu.RUnlock()

	elapsed := now.Sub(snapshot.RecordedAt).Minutes()

	if elapsed <= 0 {
		return 0
	}

	processed := r.total.count - snapshot.JobsProcessed

	if processed <= 0 {
		return 0
	}

	return float64(processed) / elapsed
}

func (b metricBucket) add(runtime time.Duration) metricBucket {
	b.count++
	b.runtime += runtime

	return b
}

func (b metricBucket) summary() MetricSummary {
	if b.count == 0 {
		return MetricSummary{}
	}

	return MetricSummary{
		Throughput:     b.count,
		AverageRuntime: b.runtime / time.Duration(b.count),
	}
}

func bucketCounts(buckets map[string]metricBucket) map[string]int {
	counts := make(map[string]int, len(buckets))

	for name, bucket := range buckets {
		counts[name] = bucket.count
	}

	return counts
}

func cloneMetricsSnapshot(snapshot MetricsSnapshot) MetricsSnapshot {
	snapshot.JobThroughput = cloneIntMap(snapshot.JobThroughput)
	snapshot.QueueThroughput = cloneIntMap(snapshot.QueueThroughput)

	return snapshot
}

func cloneIntMap(values map[string]int) map[string]int {
	cloned := make(map[string]int, len(values))

	for key, value := range values {
		cloned[key] = value
	}

	return cloned
}
