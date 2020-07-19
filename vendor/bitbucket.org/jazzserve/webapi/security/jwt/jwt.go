package jwt

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

type Service struct {
	config Config
}

func New(c *Config) *Service {
	if c == nil {
		panic("can't create new JWT service with nil config")
	}

	return &Service{
		config: *c,
	}
}

type claims struct {
	jwt.StandardClaims
	Subject json.RawMessage `json:"sub,omitempty"`
}

func (s *Service) newClaims(audience Audience, meta *tokenMeta, subject interface{}) (*claims, error) {
	var now = time.Now().UTC()

	jsonSubject, err := json.Marshal(subject)
	if err != nil {
		return nil, errors.Wrap(err, "can't marshal subject for jwt claims")
	}

	claims := &claims{
		StandardClaims: jwt.StandardClaims{
			Audience:  audience,
			ExpiresAt: now.Add(meta.TTL).Unix(),
			IssuedAt:  now.Unix(),
			Issuer:    s.config.Issuer,
		},
		Subject: jsonSubject,
	}

	return claims, nil
}

func (s *Service) New(audience Audience, subject interface{}) (token string, err error) {
	meta, ok := s.config.Tokens[audience]
	if !ok {
		panic(fmt.Sprintf("jwt token meta for %v not found", audience))
	}

	claims, err := s.newClaims(audience, meta, subject)
	if err != nil {
		return
	}

	token, err = jwt.
		NewWithClaims(jwt.SigningMethodHS512, claims).
		SignedString([]byte(meta.Key))
	if err != nil {
		return
	}

	return
}

func ReadTokenSubject(token *jwt.Token, subject interface{}) error {
	err := json.Unmarshal(token.Claims.(*claims).Subject, subject)
	if err != nil {
		return errors.Wrap(err, "can't unmarshal subject")
	}
	return nil
}

func (s *Service) keyFunc(audience Audience) func(*jwt.Token) (interface{}, error) {
	return func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		meta, ok := s.config.Tokens[audience]
		if !ok {
			panic(fmt.Sprintf("jwt token meta for %v not found", audience))
		}

		claims := token.Claims.(*claims)

		if !claims.VerifyAudience(audience, true) {
			return nil, errors.New("audience mismatch")
		}

		if !claims.VerifyIssuer(s.config.Issuer, true) {
			return nil, errors.New("invalid issuer")
		}

		return []byte(meta.Key), nil
	}
}

func (s *Service) ParseToken(t string, audience Audience) (token *jwt.Token, err error) {
	token, err = jwt.ParseWithClaims(t, &claims{}, s.keyFunc(audience))
	return token, errors.WithStack(err)
}
