// Ref: @bedrock/code-0254
// upstream framework 13.x. It defines the FailedJobProvider contract plus
// the optional Countable and Prunable extensions, and ships five
// implementations that mirror the upstream providers:
//
//   - DatabaseFailedJobProvider     — integer-keyed SQL table
//   - DatabaseUuidFailedJobProvider — UUID-keyed SQL table
//   - FileFailedJobProvider         — single JSON file on disk
//   - DynamoDbFailedJobProvider     — DynamoDB-backed, mockable client
//   - NullFailedJobProvider         — no-op for testing / disabled state
//
// Each provider has its own test file under packages/queue/failed/.
package failed
