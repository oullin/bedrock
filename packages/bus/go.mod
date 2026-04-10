module github.com/bedrock/packages/bus

go 1.26.0

require github.com/bedrock/packages/queue v0.0.0

replace (
	github.com/bedrock/packages/cache => ../cache
	github.com/bedrock/packages/queue => ../queue
)
