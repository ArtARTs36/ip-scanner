package config

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/artarts36/ip-scanner/internal/infrastructure/storage"

	goMetrics "github.com/artarts36/go-metrics"
)

type Config struct {
	Log LogConfig `envPrefix:"LOG_"`

	Metrics Metrics `envPrefix:"METRICS_"`

	Storage storage.Config `envPrefix:"STORAGE_"`
}

type LogConfig struct {
	Level slog.Level `env:"LEVEL"`
}

type Metrics struct {
	Server struct {
		Addr    string        `env:"ADDR"`
		Timeout time.Duration `env:"TIMEOUT" envDefault:"30s"`
	} `envPrefix:"SERVER_"`
}

func (m *Metrics) ToGoMetrics() goMetrics.Config {
	addr := m.Server.Addr
	if addr == "" {
		addr = ":8081"
	}

	return goMetrics.Config{
		Server: goMetrics.ServerConfig{
			Addr:    addr,
			Timeout: m.Server.Timeout,
		},
		Namespace: "ip_scanner",
	}
}

func (c *Config) Validate(envPrefix string) error {
	err := c.Storage.Validate(fmt.Sprintf("%s%s", envPrefix, "STORAGE_"))
	if err != nil {
		return fmt.Errorf("validate storage: %w", err)
	}
	return nil
}
