// Package notifications provides a Laravel-inspired notification system for
// sending messages across multiple channels (mail, database, broadcast, and
// custom drivers). Notifications are dispatched through a channel manager that
// lazily resolves drivers, supports queued delivery via the bus package, and
// fires lifecycle events (sending, sent, failed) through the event dispatcher.
package notifications
