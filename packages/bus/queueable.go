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
func (q *Queueable) OnConnection(connection string) *Queueable {
	q.Connection = connection

	return q
}

// OnQueue sets the queue name.
func (q *Queueable) OnQueue(queue string) *Queueable {
	q.Queue = queue

	return q
}

// WithDelay sets the dispatch delay.
func (q *Queueable) WithDelay(d time.Duration) *Queueable {
	q.Delay = d

	return q
}

// Chain sets a sequence of jobs to run after this one succeeds.
func (q *Queueable) Chain(jobs ...any) *Queueable {
	q.ChainJobs = jobs

	return q
}

// AppendToChain appends jobs to the existing chain.
func (q *Queueable) AppendToChain(jobs ...any) *Queueable {
	q.ChainJobs = append(q.ChainJobs, jobs...)

	return q
}

// Batchable embeds batch membership information into a job struct.
type Batchable struct {
	BatchID string
}

// Batching reports whether the job is part of a batch.
func (b *Batchable) Batching() bool { return b.BatchID != "" }

// WithBatchID sets the batch ID.
func (b *Batchable) WithBatchID(id string) { b.BatchID = id }
