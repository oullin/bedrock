// Package failed contains the Go port of Illuminate\Queue\Failed\* from
// laravel/framework 13.x. It defines the FailedJobProvider contract plus
// the optional Countable and Prunable extensions, and ships five
// implementations that mirror the Laravel providers:
//
//   - DatabaseFailedJobProvider     — integer-keyed SQL table
//   - DatabaseUuidFailedJobProvider — UUID-keyed SQL table
//   - FileFailedJobProvider         — single JSON file on disk
//   - DynamoDbFailedJobProvider     — DynamoDB-backed, mockable client
//   - NullFailedJobProvider         — no-op for testing / disabled state
//
// Tests under this package are a 1:1 port of the four PHPUnit suites
// (DatabaseFailedJobProviderTest, DatabaseUuidFailedJobProviderTest,
// FileFailedJobProviderTest, DynamoDbFailedJobProviderTest). Every ported
// test carries a `// Port of …` header so scripts/queue-parity.sh can
// match it back to its upstream counterpart.
package failed
