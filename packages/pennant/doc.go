// Package pennant provides feature flags. It defines a
// two-level abstraction: Driver (low-level backend) and Decorator (caching
// + event-dispatch wrapper). ArrayDriver provides in-memory storage;
// DatabaseDriver provides SQL-backed persistence. A Manager coordinates
// named driver instances and a ScopedFeatureInteraction provides the fluent
// scope-bound API.
package pennant
