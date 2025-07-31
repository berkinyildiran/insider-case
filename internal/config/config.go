package config

import (
	"fmt"
	"github.com/berkinyildiran/insider-case/internal/cache"
	"github.com/berkinyildiran/insider-case/internal/database"
	"github.com/berkinyildiran/insider-case/internal/scheduler"
	"github.com/berkinyildiran/insider-case/internal/sender"
	"github.com/berkinyildiran/insider-case/internal/server"
	"github.com/berkinyildiran/insider-case/internal/validator"
	"github.com/spf13/viper"
)

type Config struct {
	validator *validator.Validator
	Cache     *cache.Config     `mapstructure:"cache" validate:"required"`
	Database  *database.Config  `mapstructure:"database" validate:"required"`
	Scheduler *scheduler.Config `mapstructure:"scheduler" validate:"required"`
	Sender    *sender.Config    `mapstructure:"sender" validate:"required"`
	Server    *server.Config    `mapstructure:"server" validate:"required"`

	viper *viper.Viper
}

func NewConfig(validator *validator.Validator) *Config {
	return &Config{
		validator: validator,

		viper: viper.New(),
	}
}

func (c *Config) Load(name string, path string) error {
	c.viper.SetConfigName(name)
	c.viper.SetConfigType("yaml")
	c.viper.AddConfigPath(path)

	if err := c.viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := c.viper.Unmarshal(c); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}

func (c *Config) Validate() error {
	if err := c.validator.Validate(c); err != nil {
		return fmt.Errorf("failed to validate config: %w", err)
	}

	return nil
}
