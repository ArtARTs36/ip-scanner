package app

import (
	"fmt"

	env "github.com/caarlos0/env/v11"
	"github.com/pkg/errors"

	"github.com/artarts36/ip-scanner/internal/config"
)

// Config struct.
type Config struct {
	config.Config

	GRPC GRPCConfig `envPrefix:"GRPC_"`
}

type GRPCConfig struct {
	Port          int  `env:"PORT,required"`
	UseReflection bool `env:"USE_REFLECTION"`
}

// ReadConfig func.
func ReadConfig(prefix string) (*Config, error) {
	c := &Config{}
	opts := env.Options{
		Prefix: prefix,
	}

	if err := env.ParseWithOptions(c, opts); err != nil {
		return nil, errors.Wrap(err, "init config failed")
	}

	if err := c.Validate(prefix); err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	return c, nil
}

func (c *Config) Validate(envPrefix string) error {
	return c.Config.Validate(envPrefix)
}
