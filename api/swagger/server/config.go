package server

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"time"
)

type Config struct {
	JwtSecret   string `split_words:"true"`
	JwtDuration string `split_words:"true"`
}

func LoadConfig() (*Config, error) {
	// Load config
	var c Config
	var err error
	err = envconfig.Process("detectordag", &c)
	if err != nil {
		return nil, err
	}
	// Ensure duration is valid
	dur, err := c.ParseDuration()
	if err != nil {
		return nil, err
	}
	if dur.Seconds() < 1 {
		return nil, fmt.Errorf("JWT expiry duration insufficient: %f", dur.Seconds())
	}
	return &c, nil
}

func (c *Config) ParseDuration() (time.Duration, error) {
	return time.ParseDuration(c.JwtDuration)
}
