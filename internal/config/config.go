package config

import "github.com/spf13/viper"

type Config struct {
	ServerPort string `mapstructure:"SERVER_PORT"`
	DBSource   string `mapstructure:"DB_SOURCE"`
	RedisAddr  string `mapstructure:"REDIS_ADDR"`
	SentryDSN  string `mapstructure:"SENTRY_DSN"`
}

func Load() (*Config, error) {
	// Tell viper to look for environment variables
	viper.AutomaticEnv()

	// You can also tell it to read from a file (optional)
	// viper.SetConfigName("config")
	// viper.AddConfigPath(".")
	// viper.ReadInConfig()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
