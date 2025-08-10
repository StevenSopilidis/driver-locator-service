package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	UdpAddr               string        `mapstructure:"UDP_ADDRESS"`
	UdpPort               int           `mapstructure:"UDP_PORT"`
	MaxConcurrentRequests int           `mapstructure:"MAX_CONCURRENT_REQUESTS"`
	RedisAddr             string        `mapstructure:"REDIS_ADDRESS"`
	RedisPassword         string        `mapstructure:"REDIS_PASSWORD"`
	RedisDB               int           `mapstructure:"REDIS_DB"`
	TTL                   time.Duration `mapstructure:"TTL"`
}

func NewConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	// Force reading file; if it fails, that's a real error now
	if err := viper.ReadInConfig(); err != nil {
		return config, fmt.Errorf("error reading config file: %w", err)
	}

	viper.AutomaticEnv() // allow env override

	err = viper.Unmarshal(&config)
	return config, err
}
