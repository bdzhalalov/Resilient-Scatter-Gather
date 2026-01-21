package main

import (
	"github.com/bdzhalalov/resilient-scatter-gather/config"
	"github.com/bdzhalalov/resilient-scatter-gather/internal/server"
	"github.com/bdzhalalov/resilient-scatter-gather/pkg/logger"
)

func main() {
	cfg := config.InitConfig()

	log := logger.Logger(&cfg)

	apiServer := server.Init(&cfg, log)

	if err := apiServer.Run(); err != nil {
		log.Fatalf("Can't start server: %v", err)
	}
}
