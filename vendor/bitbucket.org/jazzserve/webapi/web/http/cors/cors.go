package cors

import (
	"net/http"

	"github.com/gorilla/handlers"
)

type CORS struct {
	Config  *Config `yaml:"cors" validate:"required"`
	headers []string
	methods []string
}

func New(c *Config, headers, methods []string) *CORS {
	if c == nil || len(headers) == 0 || len(methods) == 0 {
		panic("cannot init CORS with zero input")
	}
	return &CORS{
		Config:  c,
		headers: headers,
		methods: methods,
	}
}

func (cors *CORS) Middleware() func(http.Handler) http.Handler {
	return handlers.CORS(
		handlers.AllowedHeaders(cors.headers),
		handlers.AllowedOrigins(cors.Config.Origins),
		handlers.MaxAge(cors.Config.MaxAge),
		handlers.AllowedMethods(cors.methods),
		handlers.AllowCredentials())
}
