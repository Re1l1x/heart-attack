package config

import (
	"log/slog"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"github.com/jus1d/kypidbot/internal/lib/logger/sl"
)

const (
	EnvLocal       = "local"
	EnvDevelopment = "dev"
	EnvProduction  = "prod"
)

type Config struct {
	Env      string   `yaml:"env" env-required:"true"`
	Telegram Telegram `yaml:"telegram" env-required:"true"`
	Ollama   Ollama   `yaml:"ollama" env-required:"true"`
	Postgres Postgres `yaml:"postgres" env-required:"true"`
}

type Ollama struct {
	Host  string `yaml:"host" env-required:"true"`
	Port  string `yaml:"port" env-required:"true"`
	Model string `yaml:"model" env-required:"true"`
}

type Telegram struct {
	Token string `yaml:"token" env-required:"true"`
}

type Postgres struct {
	Host     string `yaml:"host" env-required:"true"`
	Port     string `yaml:"port" env-required:"true"`
	User     string `yaml:"user" env-required:"true"`
	Name     string `yaml:"name" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	ModeSSL  string `yaml:"sslmode" env-required:"true"`
}

// MustLoad loads config to a new Config instance and return it
func MustLoad() *Config {
	_ = godotenv.Load()

	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		slog.Error("missed CONFIG_PATH parameter")
		os.Exit(1)
	}

	var err error
	if _, err = os.Stat(configPath); os.IsNotExist(err) {
		slog.Error("config file does not exist", slog.String("path", configPath))
		os.Exit(1)
	}

	var config Config

	if err = cleanenv.ReadConfig(configPath, &config); err != nil {
		slog.Error("cannot read config", sl.Err(err))
		os.Exit(1)
	}

	return &config
}

func Empty() *Config {
	return &Config{}
}
