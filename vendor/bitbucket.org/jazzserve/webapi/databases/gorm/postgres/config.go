package postgres

import (
	"fmt"
	"strings"
)

const (
	Driver        = "postgres"
	DefaultSchema = "public"
)

type Config struct {
	User            string `yaml:"user" validate:"required"`
	Password        string `yaml:"password" validate:"required"`
	Host            string `yaml:"host" validate:"required"`
	Port            int    `yaml:"port" validate:"required"`
	Database        string `yaml:"database" validate:"required"`
	MaxIdleConns    int    `yaml:"max-idle-conns"`
	MaxOpenConns    int    `yaml:"max-open-conns" validate:"required"`
	ApplicationName string `yaml:"application-name"`
	SSLMode         string `yaml:"ssl-mode" validate:"required"`
	Schema          string `yaml:"schema"`
}

func (c *Config) composeConnectionString() string {
	var schema = DefaultSchema
	if s := strings.TrimSpace(c.Schema); s != "" {
		schema = s
	}

	return fmt.Sprintf("%s://%s:%s@%s:%v/%s?application_name=%s&sslmode=%s&search_path=%s",
		Driver, c.User, c.Password, c.Host, c.Port, c.Database,
		c.ApplicationName, c.SSLMode, schema)
}
