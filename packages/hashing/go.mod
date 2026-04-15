module github.com/bedrock/packages/hashing

go 1.26.0

replace (
	github.com/bedrock/packages/container => ../container
	github.com/bedrock/packages/contracts => ../contracts
)

require (
	github.com/bedrock/packages/container v0.0.0
	github.com/bedrock/packages/contracts v0.0.0
	golang.org/x/crypto v0.50.0
)

require golang.org/x/sys v0.43.0 // indirect
