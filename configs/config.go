package configs

import (
	"github.com/spf13/viper"
)

var cfg = &conf{}

type conf struct {
	BrasilAPIURL string `mapstructure:"BRASIL_API_URL"`
	ViaCEPURL    string `mapstructure:"VIA_CEP_URL"`
}

func NewConfig() *conf {
	return cfg
}

func init() {
	viper.SetConfigName("app_config")
	viper.SetConfigType("yaml")
	viper.SetConfigFile("./resources/env.yaml")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
}
