package config

import (
	"fmt"
	"path"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		HTTP     `yaml:"http"`
		Database `yaml:"database"`
		Log      `yaml:"log"`
		Redis    `yaml:"redis"`
		Email    `yaml:"email"`
		S3Data   `yaml:"s3"`
	}

	HTTP struct {
		Port        string        `env-required:"true" yaml:"port" env:"SERVER_PORT"`
		Address     string        `env-required:"true" yaml:"address" env:"SERVER_ADDRESS"`
		Timeout     string        `env-required:"true" yaml:"timeout"`
		IdleTimeout time.Duration `env-required:"true" yaml:"idle_timeout"`
	}

	Database struct {
		Conn        string `env-required:"true" env:"POSTGRES_CONN"`
		MaxPoolSize int    `env-required:"true" yaml:"max_pool_size" env:"MAX_POOL_SIZE"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"level" env:"LOG_LEVEL"`
	}

	Redis struct {
		Addr     string `env-required:"true" yaml:"addr" env:"REDIS_ADDR"`
		Password string `env-required:"true" yaml:"password" env:"REDIS_PASSWORD"`
		DB       int    `env-required:"true" yaml:"db" env:"REDIS_DB"`
	}

	Email struct {
		FromEmail string `env-required:"true" yaml:"from_email" env:"FROM_EMAIL"`
		Password  string `env-required:"true" yaml:"password" env:"FROM_EMAIL_PASSWORD"`
		SMTP      string `env-required:"true" yaml:"smtp" env:"FROM_EMAIL_SMTP"`
		Addr      string `env-required:"true" yaml:"smtp_addr" env:"SMTP_ADDR"`
	}

	S3Data struct {
		BucketName       string `env-required:"true" yaml:"port" env:"BUCKET_NAME"`
		Region           string `env-required:"true" yaml:"port" env:"REGION"`
		EndpointResolver string `env-required:"true" yaml:"port" env:"ENDPOINT_RESOLVER"`
	}
)

func NewConfig(configPath string) (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig(path.Join("./", configPath), cfg)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	err = cleanenv.UpdateEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("error updating env: %w", err)
	}

	return cfg, nil
}
