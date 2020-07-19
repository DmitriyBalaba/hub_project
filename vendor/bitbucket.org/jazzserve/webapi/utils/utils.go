package utils

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"os"
)

const (
	ErrorExitCode   = -1
	SuccessExitCode = 0
)

// CheckCloseError checks for closeErr on Close() and sets to err if it's empty
func CheckCloseError(c io.Closer, err *error) {
	if err == nil {
		panic("invalid call to CheckCloseError with nil err ptr")
	}
	closeErr := c.Close()
	if closeErr != nil && *err == nil {
		*err = errors.Wrap(closeErr, "can't close")
	}
}

func GenerateRandomString(length int) (result string, err error) {
	data := make([]byte, length)
	if _, err = rand.Read(data); err != nil {
		return
	}
	return hex.EncodeToString(data[:]), nil
}

func InitDefaultZeroLogger() {
	// make log produce human-friendly, colorized output
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	// adding caller function to log messages globally
	log.Logger = log.With().Caller().Logger()
}
