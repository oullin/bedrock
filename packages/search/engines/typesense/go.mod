module github.com/bedrock/packages/search/engines/typesense

go 1.26.0

require (
	github.com/bedrock/packages/contracts v0.0.0
	github.com/bedrock/packages/search v0.0.0
	github.com/typesense/typesense-go/v3 v3.0.0
)

replace (
	github.com/bedrock/packages/contracts => ../../../contracts
	github.com/bedrock/packages/search => ../..
)
