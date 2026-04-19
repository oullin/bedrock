module github.com/bedrock/packages/search/engines/algolia

go 1.26.0

require (
	github.com/bedrock/packages/contracts v0.0.0
	github.com/bedrock/packages/search v0.0.0
	github.com/algolia/algoliasearch-client-go/v4 v4.0.0
)

replace (
	github.com/bedrock/packages/contracts => ../../../contracts
	github.com/bedrock/packages/search => ../..
)
