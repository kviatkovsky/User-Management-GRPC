package main

import (
	"github.com/kviatkovsky/User-Management-gRPC/internal/config"
	"github.com/kviatkovsky/User-Management-gRPC/internal/grpc/server"
	"github.com/kviatkovsky/User-Management-gRPC/internal/logger"
)

func main() {
	cfg := config.MustLoad()
	log := logger.SetupLogger(cfg.Env)
	log.Info("starting application")
	server.StartServer(log, cfg)
}
