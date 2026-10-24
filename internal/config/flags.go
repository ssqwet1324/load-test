package config

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"time"
)

type Config struct {
	URL         string
	Requests    int
	Concurrency int
	Timeout     time.Duration
}

// Parse читает флаги командной строки и возвращает валидный Config
func Parse() (Config, error) {
	var cfg Config

	flag.StringVar(&cfg.URL, "url", "", "target URL (required)")
	flag.IntVar(&cfg.Requests, "n", 100, "total number of requests")
	flag.IntVar(&cfg.Concurrency, "c", 10, "number of concurrent workers")
	flag.DurationVar(&cfg.Timeout, "timeout", 5*time.Second, "timeout per request")

	flag.Parse()

	if err := cfg.validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// validate проверяет, что параметры имеют смысл
func (cfg Config) validate() error {
	// URL
	if cfg.URL == "" {
		return errors.New("target URL is required")
	}

	parseURL, err := url.ParseRequestURI(cfg.URL)
	if err != nil {
		return fmt.Errorf("invalid target URL: %w", err)
	}

	if parseURL.Scheme != "http" && parseURL.Scheme != "https" {
		return fmt.Errorf("unsupported scheme %q: only http and https", parseURL.Scheme)
	}

	if parseURL.Hostname() == "" {
		return errors.New("target URL must have a hostname")
	}

	// Requests, Concurrency, Timeout
	if cfg.Requests <= 0 || cfg.Concurrency <= 0 || cfg.Timeout <= 0 {
		return errors.New("-n, -c and -timeout must be positive")
	}

	return nil
}
