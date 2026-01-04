package config

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Server   ServerConfig   `toml:"server"`
	Database DatabaseConfig `toml:"database"`
	JWT      JWTConfig      `toml:"jwt"`
}

type ServerConfig struct {
	Port  int  `toml:"port"`
	Debug bool `toml:"debug"`
}

type DatabaseConfig struct {
	Driver   string `toml:"driver"`
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Name     string `toml:"name"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Path     string `toml:"path"` // For SQLite
}

type JWTConfig struct {
	Secret           string `toml:"secret"`
	TokenExpiryHours int    `toml:"token_expiry_hours"`
}

// LoadConfig loads the configuration from the TOML file
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	err = toml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// GetDSN returns the database connection string based on the driver
func (c *DatabaseConfig) GetDSN() string {
	switch c.Driver {
	case "sqlite":
		return c.Path
	case "postgres":
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			c.Host, c.Port, c.User, c.Password, c.Name)
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			c.User, c.Password, c.Host, c.Port, c.Name)
	default:
		return c.Path
	}
}

// GetServerAddress returns the server address with port
func (c *ServerConfig) GetServerAddress() string {
	return fmt.Sprintf(":%d", c.Port)
}
