package app

import (
	"blob/internal/delivery/telegram"
	"github.com/spf13/viper"
)

type Config struct {
	Telegram telegram.Config `json:"telegram" yaml:"telegram"`
}

func LoadConfig(path string) (config *Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
