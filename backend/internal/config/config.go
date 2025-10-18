package config

import "github.com/kelseyhightower/envconfig"

// Config структура конфигурации приложения
type Config struct {
	Port        string `envconfig:"PORT"`
	DatabaseURL string `envconfig:"POSTGRES_URL"`
	RedisHost   string `envconfig:"REDIS_HOST"`
}

// Load загружает конфигурацию из переменных окружения
func Load() (*Config, error) {
	cfg := &Config{}
	err := envconfig.Process("", cfg)
	return cfg, err
}
