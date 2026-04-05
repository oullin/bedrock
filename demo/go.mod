module github.com/gollin/demo

go 1.26.0

require (
	github.com/gollin/packages/framework/console v0.0.0
	github.com/gollin/packages/framework/foundation v0.0.0
	github.com/gollin/packages/framework/routing v0.0.0
)

require (
	github.com/fsnotify/fsnotify v1.8.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/gollin/packages/framework/config v0.0.0 // indirect
	github.com/gollin/packages/framework/http v0.0.0 // indirect
	github.com/gollin/packages/framework/support v0.0.0 // indirect
	github.com/gollin/packages/framework/view v0.0.0 // indirect
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

replace github.com/gollin/packages/framework/config => ../packages/framework/config

replace github.com/gollin/packages/framework/console => ../packages/framework/console

replace github.com/gollin/packages/framework/foundation => ../packages/framework/foundation

replace github.com/gollin/packages/framework/http => ../packages/framework/http

replace github.com/gollin/packages/framework/routing => ../packages/framework/routing

replace github.com/gollin/packages/framework/support => ../packages/framework/support

replace github.com/gollin/packages/framework/view => ../packages/framework/view
