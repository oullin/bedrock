// Package cache provides caching primitives. It defines a
// two-level abstraction: Store (low-level backend operations) and Repository
// (high-level helpers including remember, tags, and distributed locks).
// Multiple concrete store implementations are provided under stores/.
package cache
