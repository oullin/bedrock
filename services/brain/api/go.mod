module github.com/bedrock/services/brain/api

go 1.26.0

require (
	github.com/bedrock/packages/console v0.0.0
	github.com/bedrock/packages/filesystem v0.0.0
	github.com/bedrock/packages/httpx v0.0.0
	github.com/bedrock/packages/pipeline v0.0.0
	github.com/fsnotify/fsnotify v1.10.1
	golang.org/x/tools v0.45.0
)

require (
	github.com/bedrock/packages/container v0.0.0 // indirect
	github.com/bedrock/packages/contracts v0.0.0 // indirect
	golang.org/x/mod v0.36.0 // indirect
	golang.org/x/sync v0.20.0 // indirect
	golang.org/x/sys v0.44.0 // indirect
)

replace (
	github.com/bedrock/packages/collection => ../../../packages/collection
	github.com/bedrock/packages/conditionable => ../../../packages/conditionable
	github.com/bedrock/packages/console => ../../../packages/console
	github.com/bedrock/packages/container => ../../../packages/container
	github.com/bedrock/packages/contracts => ../../../packages/contracts
	github.com/bedrock/packages/cookie => ../../../packages/cookie
	github.com/bedrock/packages/encryption => ../../../packages/encryption
	github.com/bedrock/packages/filesystem => ../../../packages/filesystem
	github.com/bedrock/packages/httpx => ../../../packages/httpx
	github.com/bedrock/packages/jsonx => ../../../packages/jsonx
	github.com/bedrock/packages/pipeline => ../../../packages/pipeline
	github.com/bedrock/packages/process => ../../../packages/process
	github.com/bedrock/packages/redis => ../../../packages/redis
	github.com/bedrock/packages/session => ../../../packages/session
	github.com/bedrock/packages/str => ../../../packages/str
	github.com/bedrock/packages/support => ../../../packages/support
)
