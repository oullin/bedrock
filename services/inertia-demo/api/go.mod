module github.com/bedrock/services/inertia-demo/api

go 1.26.0

require (
	github.com/bedrock/packages/config v0.0.0
	github.com/bedrock/packages/encryption v0.0.0
	github.com/bedrock/packages/inertia v0.0.0
	github.com/bedrock/packages/seo v0.0.0
	github.com/bedrock/packages/validation v0.0.0
	github.com/bedrock/packages/wayfinder v0.0.0
	github.com/spf13/viper v1.20.1
	golang.org/x/crypto v0.50.0
	modernc.org/sqlite v1.48.0
)

require (
	github.com/bedrock/packages/container v0.0.0 // indirect
	github.com/bedrock/packages/contracts v0.0.0 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/fsnotify/fsnotify v1.8.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/ncruces/go-strftime v1.0.0 // indirect
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	github.com/sagikazarmark/locafero v0.7.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.12.0 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.9.0 // indirect
	golang.org/x/sys v0.43.0 // indirect
	golang.org/x/text v0.36.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	modernc.org/libc v1.70.0 // indirect
	modernc.org/mathutil v1.7.1 // indirect
	modernc.org/memory v1.11.0 // indirect
)

replace (
	github.com/bedrock/packages/config => ../../../packages/config
	github.com/bedrock/packages/container => ../../../packages/container
	github.com/bedrock/packages/contracts => ../../../packages/contracts
	github.com/bedrock/packages/encryption => ../../../packages/encryption
	github.com/bedrock/packages/inertia => ../../../packages/inertia
	github.com/bedrock/packages/seo => ../../../packages/seo
	github.com/bedrock/packages/validation => ../../../packages/validation
	github.com/bedrock/packages/wayfinder => ../../../packages/wayfinder
)
