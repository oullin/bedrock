module github.com/bedrock/packages/helpers

go 1.26.0

require github.com/bedrock/packages/support v0.0.0

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/oklog/ulid/v2 v2.1.1 // indirect
	github.com/yuin/goldmark v1.8.2 // indirect
)

replace github.com/bedrock/packages/support => ../support
