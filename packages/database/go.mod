module github.com/bedrock/packages/database

go 1.26.0

require (
	github.com/bedrock/packages/container v0.0.0
	github.com/bedrock/packages/contracts v0.0.0
	github.com/bedrock/packages/pagination v0.0.0
	github.com/bedrock/packages/support v0.0.0
)

replace (
	github.com/bedrock/packages/container => ../container
	github.com/bedrock/packages/contracts => ../contracts
	github.com/bedrock/packages/pagination => ../pagination
	github.com/bedrock/packages/support => ../support
)
