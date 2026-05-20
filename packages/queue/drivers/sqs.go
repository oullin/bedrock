package drivers

import (
	"context"
	"strings"
	"time"

	"github.com/bedrock/packages/queue"
)

// SQSClient is the interface for an AWS SQS client.
type SQSClient interface {
	// SendMessage sends a message to the given queue URL. Returns message ID.
	SendMessage(ctx context.Context, queueURL string, body string, delay time.Duration) (string, error)
	// SendMessageBatch sends multiple messages. Returns a slice of message IDs.
	SendMessageBatch(ctx context.Context, queueURL string, bodies []string) ([]string, error)
	// ReceiveMessage polls for messages. Returns up to maxMessages.
	ReceiveMessage(ctx context.Context, queueURL string, maxMessages int, waitSeconds int) ([]SQSMessage, error)
	// DeleteMessage deletes a message by receipt handle.
	DeleteMessage(ctx context.Context, queueURL string, receiptHandle string) error
	// ChangeMessageVisibility changes the visibility timeout (release semantics).
	ChangeMessageVisibility(ctx context.Context, queueURL string, receiptHandle string, visibility time.Duration) error
	// GetQueueAttributes returns queue attributes. Returns map with ApproximateNumberOfMessages etc.
	GetQueueAttributes(ctx context.Context, queueURL string, attributes []string) (map[string]string, error)
}

// SQSFIFOSender is an optional interface implemented by SQSClient
// implementations that support FIFO queue extras (MessageGroupId and
// MessageDeduplicationId). The SQSDriver probes for this interface via
// a type assertion before calling it from PushFIFO, so regular
// SQSClient implementations continue to work unchanged.
type SQSFIFOSender interface {
	SendMessageFIFO(ctx context.Context, queueURL, body, messageGroupID, messageDeduplicationID string, delay time.Duration) (string, error)
}

// SQSQueueLister is an optional interface implemented by SQSClient
// implementations that can list the SQS queues visible to the
// configured credentials. The driver uses it to populate QueueNames
// so the manager-level cross-queue inspection helpers can fan out
// without the caller pre-declaring every queue name.
type SQSQueueLister interface {
	ListQueueURLs(ctx context.Context, prefix string) ([]string, error)
}

// SQSMessage is a received SQS message.
type SQSMessage struct {
	MessageID     string
	ReceiptHandle string
	Body          string
}

// SQSDriver enqueues jobs via AWS SQS.
type SQSDriver struct {
	client       SQSClient
	queueURLs    map[string]string // queueName → SQS URL (legacy explicit mapping)
	connection   string
	prefix       string // URL prefix, e.g. "https://sqs.us-east-1.amazonaws.com/1234/"
	defaultQueue string // default logical queue name when "" is supplied
	suffix       string // optional name suffix (e.g. "-staging") applied before ".fifo"
}

// NewSQSDriver creates an SQSDriver. queueURLs maps logical queue names to SQS queue URLs.

type sqsJob struct{ BaseJob }

func NewSQSDriver(client SQSClient, queueURLs map[string]string, connection string) *SQSDriver {
	return &SQSDriver{client: client, queueURLs: queueURLs, connection: connection}
}

// SetPrefix sets the SQS queue URL prefix. Returns the
// driver for chaining. When set, GetQueue composes the full URL from
// prefix + queue name (with optional suffix and FIFO-awareness),
// mirroring Illuminate\Queue\SqsQueue::getQueue().
func (d *SQSDriver) SetPrefix(prefix string) *SQSDriver {
	d.prefix = prefix

	return d
}

// SetDefault sets the default logical queue name used when the caller
// passes an empty queue string. Mirrors the upstream $default constructor
// argument.
func (d *SQSDriver) SetDefault(name string) *SQSDriver {
	d.defaultQueue = name

	return d
}

// SetSuffix sets the name suffix applied by GetQueue before any
// trailing ".fifo", mirroring the upstream $suffix.
func (d *SQSDriver) SetSuffix(suffix string) *SQSDriver {
	d.suffix = suffix

	return d
}

// GetQueue resolves a logical queue name to its SQS URL, replicating
// Illuminate\Queue\SqsQueue::getQueue() semantics:
//
//   - Empty input falls back to the configured default queue.
//   - Already-qualified URLs (starts with http:// or https://) are
//     returned unchanged.
//   - A plain name is prefixed with the prefix and has the configured
//     suffix appended, preserving any trailing ".fifo" marker.
//   - The suffix is only applied once even if it's already present in
//     the name.
func (d *SQSDriver) GetQueue(name string) string {
	if name == "" {
		name = d.defaultQueue
	}

	if strings.HasPrefix(name, "http://") || strings.HasPrefix(name, "https://") {
		return name
	}

	if url, ok := d.queueURLs[name]; ok {
		return url
	}

	return d.suffixQueue(name)
}

func (d *SQSDriver) suffixQueue(name string) string {
	prefix := strings.TrimRight(d.prefix, "/")

	if strings.HasSuffix(name, ".fifo") {
		base := strings.TrimSuffix(name, ".fifo")

		return prefix + "/" + ensureSuffix(base, d.suffix) + ".fifo"
	}

	return prefix + "/" + ensureSuffix(name, d.suffix)
}

func ensureSuffix(name, suffix string) string {
	if suffix == "" {
		return name
	}

	if strings.HasSuffix(name, suffix) {
		return name
	}

	return name + suffix
}

// IsFIFO reports whether the resolved SQS URL for a logical queue name
// targets a FIFO queue (i.e. ends in ".fifo").
func (d *SQSDriver) IsFIFO(name string) bool {
	return strings.HasSuffix(d.GetQueue(name), ".fifo")
}

// PushFIFO sends a payload with FIFO-queue-specific options. If the
// underlying client implements SQSFIFOSender, its SendMessageFIFO is
// used; otherwise the call degrades to plain SendMessage (in which case
// group/dedup are silently dropped — callers that need FIFO support
// must pair the driver with a FIFO-aware client).
func (d *SQSDriver) PushFIFO(ctx context.Context, queueName string, payload []byte, messageGroupID, messageDeduplicationID string) (string, error) {
	url := d.GetQueue(queueName)

	if sender, ok := d.client.(SQSFIFOSender); ok {
		return sender.SendMessageFIFO(ctx, url, string(payload), messageGroupID, messageDeduplicationID, 0)
	}

	return d.client.SendMessage(ctx, url, string(payload), 0)
}

func (d *SQSDriver) Push(ctx context.Context, queueName string, payload []byte) (string, error) {
	return d.client.SendMessage(ctx, d.url(queueName), string(payload), 0)
}

func (d *SQSDriver) PushDelayed(ctx context.Context, queueName string, payload []byte, delay time.Duration) (string, error) {
	return d.client.SendMessage(ctx, d.url(queueName), string(payload), delay)
}

func (d *SQSDriver) PushMultiple(ctx context.Context, queueName string, payloads [][]byte) ([]string, error) {
	bodies := make([]string, len(payloads))

	for i, p := range payloads {
		bodies[i] = string(p)
	}

	return d.client.SendMessageBatch(ctx, d.url(queueName), bodies)
}

func (d *SQSDriver) Pop(ctx context.Context, queueName string) (queue.Job, error) {
	msgs, err := d.client.ReceiveMessage(ctx, d.url(queueName), 1, 20)

	if err != nil || len(msgs) == 0 {
		return nil, queue.ErrNoJob
	}

	msg := msgs[0]
	queueURL := d.url(queueName)

	job := &sqsJob{
		BaseJob: BaseJob{
			id:         msg.MessageID,
			payload:    []byte(msg.Body),
			queue:      queueName,
			connection: d.connection,
		},
	}
	job.deleteFunc = func() error {
		return d.client.DeleteMessage(ctx, queueURL, msg.ReceiptHandle)
	}

	job.releaseFunc = func(delay time.Duration) error {
		return d.client.ChangeMessageVisibility(ctx, queueURL, msg.ReceiptHandle, delay)
	}

	job.failFunc = func(_ error) error {
		return job.deleteFunc()
	}

	return job, nil
}

func (d *SQSDriver) Size(ctx context.Context, queueName string) (int64, error) {
	attrs, err := d.client.GetQueueAttributes(ctx, d.url(queueName), []string{
		"ApproximateNumberOfMessages",
		"ApproximateNumberOfMessagesDelayed",
		"ApproximateNumberOfMessagesNotVisible",
	})

	if err != nil {
		return 0, err
	}

	return parseStatInt(attrs, "ApproximateNumberOfMessages") +
		parseStatInt(attrs, "ApproximateNumberOfMessagesDelayed") +
		parseStatInt(attrs, "ApproximateNumberOfMessagesNotVisible"), nil
}

func (d *SQSDriver) PendingSize(ctx context.Context, queueName string) (int64, error) {
	return d.Size(ctx, queueName)
}

func (d *SQSDriver) DelayedSize(_ context.Context, _ string) (int64, error) { return 0, nil }

func (d *SQSDriver) ReservedSize(ctx context.Context, queueName string) (int64, error) {
	attrs, err := d.client.GetQueueAttributes(ctx, d.url(queueName), []string{"ApproximateNumberOfMessagesNotVisible"})

	if err != nil {
		return 0, err
	}

	return parseStatInt(attrs, "ApproximateNumberOfMessagesNotVisible"), nil
}

func (d *SQSDriver) ConnectionName() string { return d.connection }

// QueueNames lists the SQS queues visible to the driver. When the
// configured client implements SQSQueueLister, the result is the
// (prefix-filtered) set of queue URLs reported by the AWS API
// transformed into their tail-segment names. When it does not, the
// driver falls back to the statically-configured queueURLs map. Both
// paths return an ErrNotSupported only if there is no inventory to
// report at all.
func (d *SQSDriver) QueueNames(ctx context.Context) ([]string, error) {
	if lister, ok := d.client.(SQSQueueLister); ok {
		urls, err := lister.ListQueueURLs(ctx, d.prefix)

		if err != nil {
			return nil, err
		}

		out := make([]string, 0, len(urls))

		for _, u := range urls {
			if name := tailSegment(u); name != "" {
				out = append(out, name)
			}
		}

		return out, nil
	}

	if len(d.queueURLs) == 0 {
		return nil, queue.ErrNotSupported
	}

	out := make([]string, 0, len(d.queueURLs))

	for name := range d.queueURLs {
		out = append(out, name)
	}

	return out, nil
}

// PendingJobs returns ErrNotSupported. SQS does not expose a way to
// peek at messages currently sitting in a queue without consuming
// them — even short-poll ReceiveMessage hides messages from other
// consumers under the visibility timeout, which is destructive. Use
// the queue's Size method for a count instead.
func (d *SQSDriver) PendingJobs(_ context.Context, _ string) ([]queue.InspectedJob, error) {
	return nil, queue.ErrNotSupported
}

// DelayedJobs returns ErrNotSupported. SQS only reports the count of
// delayed messages via the ApproximateNumberOfMessagesDelayed
// attribute; the messages themselves are inaccessible until they
// become visible.
func (d *SQSDriver) DelayedJobs(_ context.Context, _ string) ([]queue.InspectedJob, error) {
	return nil, queue.ErrNotSupported
}

// ReservedJobs returns ErrNotSupported. SQS reports
// ApproximateNumberOfMessagesNotVisible as the count of reserved
// jobs but the messages themselves belong to the consumer that holds
// the lease and cannot be enumerated from outside.
func (d *SQSDriver) ReservedJobs(_ context.Context, _ string) ([]queue.InspectedJob, error) {
	return nil, queue.ErrNotSupported
}

// tailSegment returns the segment after the last '/' in s. Used to
// distil a fully-qualified SQS URL down to its queue name.
func tailSegment(s string) string {
	if i := strings.LastIndexByte(s, '/'); i >= 0 {
		return s[i+1:]
	}

	return s
}

func (d *SQSDriver) url(queueName string) string {
	if d.prefix != "" || d.defaultQueue != "" || d.suffix != "" {
		return d.GetQueue(queueName)
	}

	if url, ok := d.queueURLs[queueName]; ok {
		return url
	}

	return queueName
}
