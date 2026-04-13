package events

import "context"

// InvokeQueuedClosure handles the execution of a QueuedClosure when it is
// processed by the queue worker.
type InvokeQueuedClosure struct{}

// Handle invokes the queued closure with the given event.
func (i *InvokeQueuedClosure) Handle(ctx context.Context, closure *QueuedClosure, event any) error {
	_, err := closure.closure(ctx, event)

	return err
}

// Failed is called when the queued closure job fails. It invokes the catch
// callback if one was registered.
func (i *InvokeQueuedClosure) Failed(ctx context.Context, closure *QueuedClosure, err error) {
	if closure.catchFn != nil {
		closure.catchFn(ctx, err)
	}
}
