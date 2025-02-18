package config

import "github.com/spf13/viper"

type Config struct {
	DatabaseURL             string `mapstructure:"DATABASE_URL"`
	Port                    string `mapstructure:"PORT"`
	GoogleClientID          string `mapstructure:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret      string `mapstructure:"GOOGLE_CLIENT_SECRET"`
	GoogleClientCallbackUrl string `mapstructure:"GOOGLE_CLIENT_CALLBACK_URL"`
}

func LoadConfig() (config Config, err error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return config, err
	}

	err = viper.Unmarshal(&config)
	return
}
