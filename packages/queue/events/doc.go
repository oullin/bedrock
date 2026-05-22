// Ref: @bedrock/code-0232
// class from upstream framework 13.x. Each upstream event is a plain Go
// struct; the worker, manager, and drivers emit them via the queue
// package's EventEmitter interface.
//
// The package deliberately uses `any` for job-shaped fields and `error`
// for exception-shaped fields. Using the queue.Job interface directly
// would create an import cycle, and consumers type-assert to the concrete
// job they care about anyway (the same pattern upstream listeners use).
package events
