package server

import (
	"bitbucket.org/jazzserve/webapi/utils"
	"context"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

type Server struct {
	*http.Server
	config *Config
}

func New(config *Config) *Server {
	if config == nil {
		panic("cannot create new server with nil config")
	}
	return &Server{
		Server: &http.Server{
			Addr: ":" + strconv.Itoa(config.Port),
			// TODO: set read and write timeouts as well
		},
		config: config,
	}
}

func (s *Server) waitForInterrupt() {
	waitTime := time.Second * time.Duration(s.config.GracefulTimeout)

	// creating a channel and waiting for Interrupt signal sent
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)
	<-signalChannel

	// Interrupt signal has been sent so the server shuts down with delay
	log.Info().Msgf("Shutting down (%d seconds timeout)", s.config.GracefulTimeout)
	ctx, cancel := context.WithTimeout(context.Background(), waitTime)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal().Msg("Cannot shut down: " + err.Error())
		os.Exit(utils.ErrorExitCode)
	}
	log.Info().Msg("Bye")
}

func (s *Server) Run() {
	if s.Handler == nil {
		panic("cannot run server with nil handler: use SetHandler()")
	}
	go func() {
		log.Info().Msgf("Listen and serve on port: %d", s.config.Port)
		if err := s.Server.ListenAndServe(); err != nil {
			log.Fatal().Msg("Cannot run server: " + err.Error())
			os.Exit(utils.ErrorExitCode)
		}
	}()
	s.waitForInterrupt()
}

func (s *Server) SetHandler(h http.Handler) {
	s.Handler = h
}
