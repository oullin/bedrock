module github.com/bedrock/packages/precognition

go 1.26.0

replace (
	github.com/bedrock/packages/container => ../container
	github.com/bedrock/packages/routing => ../routing
)

require (
	github.com/bedrock/packages/container v0.0.0
	github.com/bedrock/packages/routing v0.0.0
)
