module github.com/bedrock/packages/auth

go 1.26.0

require golang.org/x/crypto v0.37.0

replace (
	github.com/bedrock/packages/cookie => ../cookie
	github.com/bedrock/packages/session => ../session
)
