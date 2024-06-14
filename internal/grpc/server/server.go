package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/kviatkovsky/User-Management-gRPC/internal/config"
	"github.com/kviatkovsky/User-Management-gRPC/internal/domain/models/users/repositories"
	"github.com/kviatkovsky/User-Management-gRPC/internal/grpc/services/users"
	grpcUser "github.com/kviatkovsky/User-Management-gRPC/internal/grpc/user"
	"github.com/kviatkovsky/User-Management-gRPC/pkg/client/postgresql"

	"google.golang.org/grpc"
)

func StartServer(log *slog.Logger, cfg *config.Config) {
	postgresqlClient, err := postgresql.NewClient(context.TODO(), cfg.PostgresQL)

	if err != nil {
		log.Error("Error connecting to DB", "Error", err)
	}
	repository := repositories.NewUserRepository(postgresqlClient, log)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPC.Port))
	if err != nil {
		log.Error("failed to listen: %v", err)
	}

	s := users.Server{Logger: log, Repo: repository}

	grpcServer := grpc.NewServer()

	grpcUser.RegisterUserServiceServer(grpcServer, &s)

	log.Info("server started", "port", cfg.GRPC.Port)

	if err := grpcServer.Serve(lis); err != nil {
		log.Error("failed to serve", "error", err)
	}

}
