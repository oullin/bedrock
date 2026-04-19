module github.com/bedrock/packages/facades/cache

go 1.26.0

require (
	github.com/bedrock/packages/bootstrap v0.0.0
	github.com/bedrock/packages/cache v0.0.0
)

require (
	github.com/bedrock/packages/container v0.0.0 // indirect
	github.com/bedrock/packages/contracts v0.0.0 // indirect
)

replace (
	github.com/bedrock/packages/bootstrap => ../../bootstrap
	github.com/bedrock/packages/cache => ../../cache
	github.com/bedrock/packages/container => ../../container
	github.com/bedrock/packages/contracts => ../../contracts
)
