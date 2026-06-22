package support

import (
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	BaseURL string `yaml:"base_url"`
	Timeout string `yaml:"timeout"`
	Retry   int    `yaml:"retry"`
}

var cfg *Config

func GetConfig() *Config {
	if cfg == nil {
		cfg = loadConfig()
	}
	return cfg
}

func GetTimeout() time.Duration {
	d, err := time.ParseDuration(GetConfig().Timeout)
	if err != nil {
		return 30 * time.Second``
	}
	return d
}``

func loadConfig() *Config {
	c := &Config{
		BaseURL: "http://localhost:8080",
		Timeout: "30s",
		Retry:   0,
	}

	data, err := os.ReadFile("config.yaml")
	if err != nil {
		root := findProjectRoot()
		data, err = os.ReadFile(filepath.Join(root, "config.yaml"))
		if err != nil {
			return c
		}
	}

	yaml.Unmarshal(data, c)

	// Env vars override config
	if url := os.Getenv("BASE_URL"); url != "" {
		c.BaseURL = url
	}

	return c
}
