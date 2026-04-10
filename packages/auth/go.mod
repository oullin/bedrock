module github.com/bedrock/packages/auth

go 1.26.0

require golang.org/x/crypto v0.37.0

require github.com/bedrock/packages/contracts/auth v0.0.0-00010101000000-000000000000 // indirect

replace (
	github.com/bedrock/packages/contracts/auth => ../contracts/auth
	github.com/bedrock/packages/cookie => ../cookie
	github.com/bedrock/packages/session => ../session
)
