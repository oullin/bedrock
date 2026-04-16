module github.com/bedrock/packages/search/engines/meilisearch

go 1.26.0

require (
	github.com/bedrock/packages/contracts v0.0.0
	github.com/bedrock/packages/search v0.0.0
	github.com/meilisearch/meilisearch-go v0.31.0
)

replace (
	github.com/bedrock/packages/contracts => ../../../../contracts
	github.com/bedrock/packages/search => ../../../
)
