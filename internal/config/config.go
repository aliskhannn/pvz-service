package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"os"
	"time"
)

type Config struct {
	Server   `yaml:"server"`
	JWT      `yaml:"jwt"`
	Database `yaml:"database"`
}

type Server struct {
	HTTPPort    string `yaml:"http_port"`
	GRPCPort    string `yaml:"grpc_port"`
	MetricsPort string `yaml:"metrics_port"`
	Mode        string `yaml:"mode"`
}

type JWT struct {
	Secret string
	TTL    time.Duration `yaml:"ttl"`
}

type Database struct {
	User     string
	Password string
	Name     string
	Host     string
	Port     string
	SSLMode  string `yaml:"ssl_mode"`
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	cfg.JWT.Secret = os.Getenv("JWT_SECRET")
	cfg.Database.User = os.Getenv("DB_USER")
	cfg.Database.Password = os.Getenv("DB_PASSWORD")
	cfg.Database.Name = os.Getenv("DB_NAME")
	cfg.Database.Host = os.Getenv("DB_HOST")
	cfg.Database.Port = os.Getenv("DB_PORT")

	ttl, err := time.ParseDuration(viper.GetString("jwt.ttl"))
	if err != nil {
		return nil, err
	}
	cfg.JWT.TTL = ttl

	return &cfg, nil
}
