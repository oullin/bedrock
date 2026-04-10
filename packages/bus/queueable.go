package bus

import "time"

// Queueable embeds queue routing options into a command/job struct.
type Queueable struct {
	Connection string
	Queue      string
	Delay      time.Duration
	ChainJobs  []any
	Middleware []Pipe
}

// OnConnection sets the queue connection.

// OnQueue sets the queue name.

// WithDelay sets the dispatch delay.

// WithoutDelay removes any dispatch delay.

// Chain sets a sequence of jobs to run after this one succeeds.

// AppendToChain appends jobs to the existing chain.

// PrependToChain inserts jobs at the beginning of the chain.

// Through sets the middleware pipeline for this job.

// GetQueue returns the queue name.

// GetConnection returns the connection name.

// AllOnConnection sets the connection on this job and all chain jobs
// that support the OnConnection method.

// AllOnQueue sets the queue on this job and all chain jobs
// that support the OnQueue method.

// Batchable embeds batch membership information into a job struct.
type Batchable struct {
	BatchID   string
	batchInst *Batch
}

func (q *Queueable) OnConnection(connection string) *Queueable {
	q.Connection = connection

	return q
}

func (q *Queueable) OnQueue(queue string) *Queueable {
	q.Queue = queue

	return q
}

func (q *Queueable) WithDelay(d time.Duration) *Queueable {
	q.Delay = d

	return q
}

func (q *Queueable) WithoutDelay() *Queueable {
	q.Delay = 0

	return q
}

func (q *Queueable) Chain(jobs ...any) *Queueable {
	q.ChainJobs = jobs

	return q
}

func (q *Queueable) AppendToChain(jobs ...any) *Queueable {
	q.ChainJobs = append(q.ChainJobs, jobs...)

	return q
}

func (q *Queueable) PrependToChain(jobs ...any) *Queueable {
	q.ChainJobs = append(jobs, q.ChainJobs...)

	return q
}

func (q *Queueable) Through(pipes ...Pipe) *Queueable {
	q.Middleware = pipes

	return q
}

func (q *Queueable) GetQueue() string { return q.Queue }

func (q *Queueable) GetConnection() string { return q.Connection }

func (q *Queueable) AllOnConnection(connection string) *Queueable {
	q.Connection = connection

	for _, job := range q.ChainJobs {
		if c, ok := job.(interface{ OnConnection(string) *Queueable }); ok {
			c.OnConnection(connection)
		}
	}

	return q
}

func (q *Queueable) AllOnQueue(queue string) *Queueable {
	q.Queue = queue

	for _, job := range q.ChainJobs {
		if c, ok := job.(interface{ OnQueue(string) *Queueable }); ok {
			c.OnQueue(queue)
		}
	}

	return q
}

// Batching reports whether the job is part of a batch.
func (b *Batchable) Batching() bool { return b.BatchID != "" }

// WithBatchID sets the batch ID.
func (b *Batchable) WithBatchID(id string) { b.BatchID = id }

// Batch returns the parent Batch instance, if set.
func (b *Batchable) Batch() *Batch { return b.batchInst }

// SetBatch sets the parent Batch instance.
func (b *Batchable) SetBatch(batch *Batch) { b.batchInst = batch }
