module github.com/bedrock/packages/debugbar

go 1.26.0

require github.com/google/uuid v1.6.0

replace (
	github.com/bedrock/packages/config => ../config
	github.com/bedrock/packages/container => ../container
	github.com/bedrock/packages/contracts => ../contracts
	github.com/bedrock/packages/events => ../events
	github.com/bedrock/packages/httpx => ../httpx
)
