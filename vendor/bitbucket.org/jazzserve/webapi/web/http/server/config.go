package server

type Config struct {
	Port            int    `yaml:"port" validate:"required"`
	GracefulTimeout int    `yaml:"graceful-timeout" validate:"required"`
	PathPrefix      string `yaml:"path-prefix"`
}

func (c *Config) GetPathPrefix() string {
	if c == nil {
		return ""
	}
	return c.PathPrefix
}
