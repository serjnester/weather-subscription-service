package configs

import (
	"fmt"
	"net"
	"net/url"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Env        EnvType    `envconfig:"ENV" required:"true" default:"dev"`
	Debug      bool       `envconfig:"DEBUG" default:"false"`
	Port       string     `envconfig:"PORT" default:"8080"`
	DB         DB         `envconfig:"DB"`
	WeatherAPI WeatherAPI `envconfig:"WEATHER_API"`
}

type WeatherAPI struct {
	BaseURL string `split_words:"true" required:"true"`
	Key     string `split_words:"true" required:"true"`
}

type DB struct {
	Host     string `required:"true"`
	Port     string `required:"true"`
	User     string `required:"true"`
	Password string `required:"true"`
	Database string `required:"true"`

	MaxOpenConns    int           `default:"5" split_words:"true"`
	MaxIdleConns    int           `default:"5" split_words:"true"`
	ConnMaxLifetime time.Duration `default:"30s" split_words:"true"`
}

func (c *DB) ConnectionURL() string {
	connURL := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.User, c.Password),
		Host:   net.JoinHostPort(c.Host, c.Port),
		Path:   c.Database,
		RawQuery: url.Values{
			"sslmode": []string{"disable"},
		}.Encode(),
	}

	return connURL.String()
}

func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}
	return &cfg, nil
}
