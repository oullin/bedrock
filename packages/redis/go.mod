module github.com/bedrock/packages/redis

go 1.26.1

require (
	github.com/bedrock/packages/container v0.0.0
	github.com/redis/go-redis/v9 v9.7.0
)

replace github.com/bedrock/packages/container => ../container

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
)
