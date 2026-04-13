module github.com/bedrock/packages/notifications

go 1.26.0

require (
	github.com/bedrock/packages/bus v0.0.0
	github.com/bedrock/packages/contracts v0.0.0
)

replace (
	github.com/bedrock/packages/bus => ../bus
	github.com/bedrock/packages/cache => ../cache
	github.com/bedrock/packages/contracts => ../contracts
	github.com/bedrock/packages/queue => ../queue
)
