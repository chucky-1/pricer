// Package config has a configuration structure
package config

// Config contains configuration data
type Config struct {
	Host string `env:"HOST" envDefault:"localhost"`
	Port string `env:"PORT" envDefault:"6379"`
}
