package store

import "time"

const (
	network            = "tcp"
	defaultCookiePath  = "/"
	defaultMaxIdle     = 1
	defaultIdleTimeOut = 5
)

type Config struct {
	Address string         `yaml:"address" validate:"required"`
	Key     string         `yaml:"key" validate:"required"`
	MaxAge  *time.Duration `yaml:"max-age" validate:"required"`
	DB      int            `yaml:"db" validate:"required"`
}
