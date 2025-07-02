package config

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string       `yaml:"env" env:"ENV" env-default:"local" env-required:"true"`
	StorageLink *StorageLink `yaml:"storage_link"`
	HTTPServer  *HTTPServer  `yaml:"http_server"`
}

type StorageLink struct {
	SQLDriver   string `yaml:"sql_driver" env-default:"postgres"`
	SQLUser     string `yaml:"sql_user" env-default:"postgres"`
	SQLPassword string `yaml:"sql_password"`
	SQLHost     string `yaml:"sql_host" env-default:"localhost"`
	SQLPort     string `yaml:"sql_port" env-default:"5433"`
	SQLDBName   string `yaml:"sql_dbname" env-default:"url-shortener"`
	SQLSSLMode  string `yaml:"sql_sslmode" env-default:"disable"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8082"`
	Timeout     time.Duration `yaml:"timeout" env-default:"16s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist at path: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Failed to read config file: %s", err)
	}

	GetStorageLink(&cfg)

	return &cfg
}

func GetStorageLink(cfg *Config) string {
	storageLink := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.StorageLink.SQLDriver,
		cfg.StorageLink.SQLUser,
		cfg.StorageLink.SQLPassword,
		cfg.StorageLink.SQLHost,
		cfg.StorageLink.SQLPort,
		cfg.StorageLink.SQLDBName,
		cfg.StorageLink.SQLSSLMode,
	)

	if storageLink == "" {
		log.Fatal("Storage link is empty in the config file")
	}

	slog.Debug("Storage link set in config",
		slog.String("storage_link", storageLink))

	return storageLink
}
