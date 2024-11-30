package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log/slog"
	"os"
	"time"
)

type Config struct {
	Env     string        `mapstructure:"ENV"`
	Host    string        `mapstructure:"HOST"`
	Port    string        `mapstructure:"PORT"`
	Timeout time.Duration `mapstructure:"TIMEOUT"`
	Address string        `mapstructure:"API_ADRESS"`
}

// read config from ./config/
func MustReadConfig() Config {
	var cfg Config

	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AddConfigPath("./configs/")

	if err := viper.ReadInConfig(); err != nil {
		slog.Error(err.Error())
		panic(fmt.Errorf("error reading config file: %w", err))
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		panic(fmt.Errorf("unable to unmarshal into struct: %v", err))
	}

	return cfg
}

func SetEnvSecret(secret string) error {

	err := os.Setenv("SECRET_KEY", secret)

	return err
}
