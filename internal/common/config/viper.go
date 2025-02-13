package config

import (
	"strings"

	"github.com/spf13/viper"
)

func NewViperConfig() error {
	viper.SetConfigName("global")
	viper.SetConfigType("yml")
	viper.AddConfigPath("../common/config")
	viper.AllowEmptyEnv(true)
	viper.EnvKeyReplacer(strings.NewReplacer("-", "_"))
	_ = viper.BindEnv("stripe-key", "STRIPE_KEY")
	viper.AutomaticEnv()
	return viper.ReadInConfig()
}
