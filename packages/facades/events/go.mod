module github.com/bedrock/packages/facades/events

go 1.26.0

require (
	github.com/bedrock/packages/bedrock v0.0.0
	github.com/bedrock/packages/events v0.0.0
)

require (
	github.com/bedrock/packages/container v0.0.0 // indirect
	github.com/bedrock/packages/contracts v0.0.0 // indirect
)

replace (
	github.com/bedrock/packages/bedrock => ../../bedrock
	github.com/bedrock/packages/container => ../../container
	github.com/bedrock/packages/contracts => ../../contracts
	github.com/bedrock/packages/events => ../../events
)
