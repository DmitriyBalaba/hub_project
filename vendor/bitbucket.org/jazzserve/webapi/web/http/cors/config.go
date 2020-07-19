package cors

type Config struct {
	Origins []string `yaml:"origins" validate:"required"`
	MaxAge  int      `yaml:"max-age" validate:"required"`
}
