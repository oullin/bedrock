package fake

import "sync/atomic"

// fakeIDCounter is used to generate unique fake IDs across tests.
var fakeIDCounter atomic.Int64
