module github.com/bedrock/packages/facades/queue

go 1.26.0

require (
	github.com/bedrock/app v0.0.0
	github.com/bedrock/packages/queue v0.0.0
)

require (
	github.com/bedrock/packages/container v0.0.0 // indirect
	github.com/bedrock/packages/contracts v0.0.0 // indirect
)

replace (
	github.com/bedrock/app => ../../app
	github.com/bedrock/packages/cache => ../../cache
	github.com/bedrock/packages/container => ../../container
	github.com/bedrock/packages/contracts => ../../contracts
	github.com/bedrock/packages/queue => ../../queue
)
