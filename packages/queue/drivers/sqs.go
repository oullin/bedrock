package drivers

import (
	"context"
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

// SQSMessage is a received SQS message.
type SQSMessage struct {
	MessageID     string
	ReceiptHandle string
	Body          string
}

// SQSDriver enqueues jobs via AWS SQS.
type SQSDriver struct {
	client     SQSClient
	queueURLs  map[string]string // queueName → SQS URL
	connection string
}

// NewSQSDriver creates an SQSDriver. queueURLs maps logical queue names to SQS queue URLs.

type sqsJob struct{ BaseJob }

func NewSQSDriver(client SQSClient, queueURLs map[string]string, connection string) *SQSDriver {
	return &SQSDriver{client: client, queueURLs: queueURLs, connection: connection}
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
	attrs, err := d.client.GetQueueAttributes(ctx, d.url(queueName), []string{"ApproximateNumberOfMessages"})

	if err != nil {
		return 0, err
	}

	return parseStatInt(attrs, "ApproximateNumberOfMessages"), nil
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

func (d *SQSDriver) url(queueName string) string {
	if url, ok := d.queueURLs[queueName]; ok {
		return url
	}

	return queueName
}
