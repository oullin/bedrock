module github.com/bedrock/packages/app

go 1.26.0

require (
	github.com/bedrock/packages/anvil/console v0.0.0
	github.com/bedrock/packages/anvil/foundation v0.0.0
	github.com/bedrock/packages/anvil/routing v0.0.0
)

require (
	github.com/fsnotify/fsnotify v1.8.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/bedrock/packages/anvil/config v0.0.0 // indirect
	github.com/bedrock/packages/anvil/http v0.0.0 // indirect
	github.com/bedrock/packages/anvil/support v0.0.0 // indirect
	github.com/bedrock/packages/anvil/view v0.0.0 // indirect
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
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/bedrock/packages/anvil/config => ../anvil/config

replace github.com/bedrock/packages/anvil/console => ../anvil/console

replace github.com/bedrock/packages/anvil/foundation => ../anvil/foundation

replace github.com/bedrock/packages/anvil/http => ../anvil/http

replace github.com/bedrock/packages/anvil/routing => ../anvil/routing

replace github.com/bedrock/packages/anvil/support => ../anvil/support

replace github.com/bedrock/packages/anvil/view => ../anvil/view
