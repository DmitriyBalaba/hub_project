package smtp

type Config struct {
	UserName   string `yaml:"user-name" validate:"required"`
	Password   string `yaml:"password" validate:"required"`
	Host       string `yaml:"host" validate:"required"`
	Port       int    `yaml:"port"`
	From       string `yaml:"from" validate:"required"`
	ReplyTo    string `yaml:"reply-to"`
	SenderName string `yaml:"sender-name"`
}
