module github.com/bedrock/packages/facades/log

go 1.26.0

require (
	github.com/bedrock/packages/bootstrap v0.0.0
	github.com/bedrock/packages/log v0.0.0
)

require (
	github.com/bedrock/packages/config v0.0.0-00010101000000-000000000000 // indirect
	github.com/bedrock/packages/container v0.0.0 // indirect
	github.com/bedrock/packages/contracts v0.0.0 // indirect
	github.com/fsnotify/fsnotify v1.8.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/sagikazarmark/locafero v0.7.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.12.0 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/spf13/viper v1.20.1 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.9.0 // indirect
	golang.org/x/sys v0.43.0 // indirect
	golang.org/x/text v0.36.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/bedrock/packages/bootstrap => ../../bootstrap
	github.com/bedrock/packages/config => ../../config
	github.com/bedrock/packages/container => ../../container
	github.com/bedrock/packages/contracts => ../../contracts
	github.com/bedrock/packages/events => ../../events
	github.com/bedrock/packages/log => ../../log
)
