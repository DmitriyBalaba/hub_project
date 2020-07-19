package jwt

import (
	"github.com/pkg/errors"
	"time"
)

type Config struct {
	Issuer string                  `validate:"required" yaml:"issuer"`
	Tokens map[Audience]*tokenMeta `validate:"dive" yaml:"tokens"`
}

type Audience = string

type tokenMeta struct {
	TTL time.Duration `validate:"required" yaml:"ttl"`
	Key string        `validate:"required" yaml:"key"`
}

var requiredTokens = make(map[string]struct{})

func AddRequiredTokens(audience ...string) {
	for i := range audience {
		requiredTokens[audience[i]] = struct{}{}
	}
}

func checkRequired(tokens map[Audience]*tokenMeta) error {
	for k := range requiredTokens {
		if _, ok := tokens[k]; !ok {
			return errors.Errorf("required %s jwt token not found", k)
		}
	}
	return nil
}
