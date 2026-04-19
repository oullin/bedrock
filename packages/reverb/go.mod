module github.com/bedrock/packages/reverb

go 1.26.0

require (
	github.com/bedrock/packages/container v0.0.0
	github.com/bedrock/packages/contracts v0.0.0
	github.com/bedrock/packages/redis v0.0.0
	nhooyr.io/websocket v1.8.17
)

require (
	github.com/klauspost/compress v1.10.3 // indirect
	github.com/redis/go-redis/v9 v9.7.3 // indirect
	github.com/bsm/redislock v0.9.4 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
)

replace (
	github.com/bedrock/packages/container => ../container
	github.com/bedrock/packages/contracts => ../contracts
	github.com/bedrock/packages/redis => ../redis
)
