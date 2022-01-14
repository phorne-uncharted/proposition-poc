package env

import (
	"sync"

	"github.com/caarlos0/env"
)

var (
	cfg  *Config
	once sync.Once
)

// Config represents the application configuration state loaded from env vars.
type Config struct {
	AllowedSitesFile string `env:"ALLOWED_SITES_FILE" envDefault:"allowed-sites.txt"`
	AppPort          string `env:"PORT" envDefault:"8090"`
}

// LoadConfig loads the config from the environment if necessary and returns a copy.
func LoadConfig() (Config, error) {
	var err error
	once.Do(func() {
		cfg = &Config{}
		err = env.Parse(cfg)
		if err != nil {
			cfg = &Config{}
		}
	})
	return *cfg, err
}
