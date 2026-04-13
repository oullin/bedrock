module github.com/bedrock/packages/hashing

go 1.26.0

replace github.com/bedrock/packages/contracts => ../contracts

require (
	github.com/bedrock/packages/contracts v0.0.0-00010101000000-000000000000
	golang.org/x/crypto v0.50.0
)

require golang.org/x/sys v0.43.0 // indirect
