package config

type Telegram struct {
	Token              string  `mapstructure:"token" json:"token" yaml:"token"`
	TronScanApiKey     string  `mapstructure:"tron_scan_api_key" json:"tron_scan_api_key" yaml:"tron_scan_api_key"`
	GridApiKey         string  `mapstructure:"grid_api_key" json:"grid_api_key" yaml:"grid_api_key"`
	AliasKey           string  `mapstructure:"alias_key" yaml:"alias_key"`
	PrivateKey         string  `mapstructure:"private_key" yaml:"private_key"`
	ReceiveAddress     string  `mapstructure:"receive_address" yaml:"receive_address"`
	ReceiveAddressIcon string  `mapstructure:"receive_address_icon" yaml:"receive_address_icon"`
	SendAddress        string  `mapstructure:"send_address" yaml:"send_address"`
	Ratio              float64 `mapstructure:"ratio" yaml:"ratio"`
	AdvanceAmount      float64 `mapstructure:"advance_amount" yaml:"advance_amount"`
	ThresholdValue     float64 `mapstructure:"threshold_value" yaml:"threshold_value"`
}
