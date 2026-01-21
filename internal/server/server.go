package server

import (
	"errors"
	"github.com/bdzhalalov/resilient-scatter-gather/api"
	"github.com/bdzhalalov/resilient-scatter-gather/config"
	"github.com/bdzhalalov/resilient-scatter-gather/internal/service"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Server struct {
	config *config.Config
	logger *logrus.Logger
	server *http.Server
}

func Init(config *config.Config, logger *logrus.Logger) *Server {
	s := service.NewMockServices()
	router := api.Router(logger, s)

	return &Server{
		config: config,
		logger: logger,
		server: &http.Server{
			Addr:    config.Addr,
			Handler: router,
		},
	}
}

func (s *Server) Run() error {

	s.logger.Info("Running API server on port" + s.config.Addr)

	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.WithError(err).Fatal("Failed to start API server")
	}

	return nil
}
