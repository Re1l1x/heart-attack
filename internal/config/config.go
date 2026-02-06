package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
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
		panic("missed CONFIG_PATH environment variable")
	}

	var err error
	if _, err = os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var config Config

	if err = cleanenv.ReadConfig(configPath, &config); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &config
}

func Empty() *Config {
	return &Config{}
}
