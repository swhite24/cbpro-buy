package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type (
	// Config defines configuration for purchasing from coinbase pro
	Config struct {
		// Client configuration
		BaseURL    string `json:"base_url"`
		Key        string `json:"key" mapstructure:"key"`
		Passphrase string `json:"passphrase" mapstructure:"passphrase"`
		Secret     string `json:"secret" mapstructure:"secret"`
		UseSandbox bool   `json:"sandbox" mapstructure:"sandbox"`

		// Purchase configuration
		Currency    string  `json:"currency"`
		Product     string  `json:"product"`
		UseCoinbase bool    `json:"coinbase" mapstructure:"coinbase"`
		AutoDeposit bool    `json:"autodeposit" mapstructure:"autodeposit"`
		Amount      float64 `json:"amount"`
	}
)

// InitializeConfig delivers the initialized config
func InitializeConfig(flags *pflag.FlagSet) *Config {
	viper.BindPFlags(flags)

	viper.SetEnvPrefix("CBPRO_BUY")
	viper.AutomaticEnv()

	c := Config{BaseURL: "https://api.pro.coinbase.com"}
	viper.Unmarshal(&c)
	if c.UseSandbox {
		c.BaseURL = "https://api-public.sandbox.pro.coinbase.com"
	}
	return &c
}
