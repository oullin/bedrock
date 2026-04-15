module github.com/bedrock/packages/bedrock

go 1.26.0

require github.com/bedrock/packages/container v0.0.0

require github.com/bedrock/packages/contracts v0.0.0 // indirect

replace (
	github.com/bedrock/packages/container => ../container
	github.com/bedrock/packages/contracts => ../contracts
)
