// Package telescope provides a debugging and introspection tool for Go
// applications, mirroring Telescope. It captures HTTP requests,
// database queries, exceptions, log messages, events, queued jobs, cache
// operations, mail, notifications, model changes, views, commands, scheduled
// tasks, Redis commands, authorization gates, outbound HTTP client calls, and
// debug dumps. Every captured item is a typed, UUID-keyed entry grouped into
// batches and persisted via a pluggable repository contract.
package telescope
