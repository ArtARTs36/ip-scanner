package storage

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
)

type Storage interface {
	GetActualTime(ctx context.Context) (time.Time, error)
	Load(ctx context.Context) (*DB, error)
}

type Config struct {
	Local *LocalConfig `envPrefix:"LOCAL_"`
	S3    *S3Config    `envPrefix:"S3_"`

	Order      []Type            `env:"ORDER"`
	AutoUpdate *AutoUpdateConfig `envPrefix:"AUTO_UPDATE_" env:",init"`
}

type AutoUpdateConfig struct {
	Interval time.Duration `env:"INTERVAL"`
	Order    Type          `env:"ORDER"`
}

type Type string

const (
	TypeUnknown = ""
	TypeLocal   = "local"
	TypeS3      = "s3"
)

func (o Type) Validate() error {
	if !o.Valid() {
		return fmt.Errorf("order have unexpected value %q. possible values: [%s]", o, strings.Join([]string{
			TypeLocal, TypeS3,
		}, ", "))
	}
	return nil
}

func (o Type) Valid() bool {
	switch o {
	case TypeLocal:
		return true
	case TypeS3:
		return true
	}
	return false
}

func (c *Config) Validate(envPrefix string) error {
	if slices.Contains(c.Order, TypeS3) && c.S3 == nil {
		s3Cfg, err := env.ParseAsWithOptions[S3Config](env.Options{
			Prefix: fmt.Sprintf("%sS3_", envPrefix),
		})
		if err != nil {
			return fmt.Errorf("parse s3 config: %w", err)
		}

		c.S3 = &s3Cfg
	}

	if slices.Contains(c.Order, TypeLocal) {
		localCfg, err := env.ParseAsWithOptions[LocalConfig](env.Options{
			Prefix: fmt.Sprintf("%sLOCAL_", envPrefix),
		})
		if err != nil {
			return fmt.Errorf("parse local config: %w", err)
		}

		c.Local = &localCfg
	}

	if c.S3 == nil && c.Local == nil {
		return errors.New("storage driver must be filled")
	}

	for _, o := range c.Order {
		if err := o.Validate(); err != nil {
			return fmt.Errorf("order: %w", err)
		}
	}

	if err := c.AutoUpdate.Validate(); err != nil {
		return fmt.Errorf("validate auto-update: %w", err)
	}

	return nil
}

func (c *AutoUpdateConfig) Validate() error {
	if !c.Enabled() {
		return nil
	}

	if err := c.Order.Validate(); err != nil {
		return fmt.Errorf("order: %w", err)
	}

	return nil
}

func (c *AutoUpdateConfig) Enabled() bool {
	return c.Interval > 0
}
