module github.com/bedrock/packages/facades/auth

go 1.26.0

require (
	github.com/bedrock/packages/auth v0.0.0
	github.com/bedrock/packages/container v0.0.0
	github.com/bedrock/packages/contracts v0.0.0
)

require (
	golang.org/x/crypto v0.50.0 // indirect
)

replace (
	github.com/bedrock/packages/auth => ../../auth
	github.com/bedrock/packages/container => ../../container
	github.com/bedrock/packages/contracts => ../../contracts
	github.com/bedrock/packages/cookie => ../../cookie
	github.com/bedrock/packages/session => ../../session
)
