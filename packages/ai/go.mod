module github.com/bedrock/packages/ai

go 1.26.0

require (
	github.com/bedrock/packages/container v0.0.0
	github.com/bedrock/packages/contracts v0.0.0
	github.com/bedrock/packages/events v0.0.0
	github.com/bedrock/packages/pipeline v0.0.0
	github.com/bedrock/packages/queue v0.0.0
	github.com/google/uuid v1.6.0
)

replace (
	github.com/bedrock/packages/container => ../container
	github.com/bedrock/packages/contracts => ../contracts
	github.com/bedrock/packages/events    => ../events
	github.com/bedrock/packages/pipeline  => ../pipeline
	github.com/bedrock/packages/queue     => ../queue
)
