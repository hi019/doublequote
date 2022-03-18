package domain

import (
	"fmt"
)

type Config struct {
	HTTP struct {
		Port        string `mapstructure:"port"`
		Domain      string `mapstructure:"domain"`
		FrontendURL string `mapstructure:"frontend_url"`
	} `mapstructure:"http"`

	App struct {
		Secret                   string `mapstructure:"secret"`
		RequireEmailVerification bool   `mapstructure:"require_email_verification"`
		DataFolder               string `mapstructure:"data_folder"`
	} `mapstructure:"app"`

	Database struct {
		URL string `toml:"url"`
	}

	Redis struct {
		URL string `toml:"url"`
	}

	SMTP struct {
		URL      string `toml:"url"`
		Host     string `toml:"host"`
		Identity string `toml:"identity"`
		Username string `toml:"url"`
		Password string `toml:"url"`
	}
}

func (c *Config) BaseURL() string {
	var domain string
	var protocol string

	if c.HTTP.Domain != "" {
		domain = c.HTTP.Domain
		protocol = "https"
	} else {
		domain = "localhost"
		protocol = "http"
	}

	return fmt.Sprintf("%s://%s:%s", protocol, domain, c.HTTP.Port)
}

func (c *Config) ApiURL() string {
	return fmt.Sprintf("%s/api", c.BaseURL())
}

func (c *Config) FrontendURL() string {
	return c.BaseURL()
}
