package config

type App struct {
	Env     string   `mapstructure:"env" json:"env" yaml:"env"`
	AppName string   `mapstructure:"app_name" json:"app_name" yaml:"app_name"`
	Support string   `mapstructure:"support" json:"support" yaml:"support"`
	Groups  []string `mapstructure:"groups" json:"groups" yaml:"groups"`
	License string   `mapstructure:"license" json:"license" yaml:"license"`
}
