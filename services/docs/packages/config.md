# config

<!-- laravel-docs: configuration.md#configuration -->
<!-- laravel-docs: configuration.md#accessing-configuration-values -->

Configuration repository with dot-notation access.

## Overview

The `config` package provides a Laravel-inspired configuration store backed by
[Viper](https://github.com/spf13/viper). It supports dot-notation key access,
nested map values, YAML file loading, and environment variable binding.

**Module:** `github.com/bedrock/packages/config`

```bash
go get github.com/bedrock/packages/config@latest
```

## Creating a Repository

```go
// From a map of values
repo := config.New(map[string]any{
    "app": map[string]any{
        "name": "Bedrock",
        "env":  "production",
    },
    "database": map[string]any{
        "host": "localhost",
        "port": 5432,
    },
})
```

## Reading Values

```go
repo.Get("app.name")                    // "Bedrock"
repo.Get("database.port")              // 5432
repo.Get("missing.key", "default")     // "default"

repo.Has("app.name")  // true
repo.Has("missing")   // false
```

## Writing Values

```go
repo.Set("app.debug", true)
```

## Loading from a YAML File

Use `NewFromViper` to configure Viper before wrapping it:

```go
v := viper.New()
v.SetConfigFile("config/app.yaml")
v.AutomaticEnv()
_ = v.ReadInConfig()

repo := config.NewFromViper(v)
```

## All Values

```go
all := repo.All() // map[string]any with every key in all sources
```
