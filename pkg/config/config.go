package config

type Configuration struct {
	App      App      `mapstructure:"app" json:"app" yaml:"app"`
	Telegram Telegram `mapstructure:"telegram" json:"telegram" yaml:"telegram"`
	DB       DB       `mapstructure:"db" json:"db" yaml:"db"`
}
