package memory

type Config struct {
	Key    string `yaml:"key" validate:"required"`
	MaxAge int    `yaml:"max-age" validate:"required"`
}
