package userApp

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	jwtt "user_service/internal/lib/jwt"
	"user_service/internal/server"

	"google.golang.org/grpc"
)

type App struct {
	log         *slog.Logger
	gRPCServer  *grpc.Server
	port        int
	userService server.UserService
}

func New(log *slog.Logger, port int, userService server.UserService) *App {
	gRPCServer := grpc.NewServer(
		grpc.UnaryInterceptor(jwtt.AuthInterceptor(os.Getenv("APP_SECRET"))),
	)

	server.Register(gRPCServer, userService)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) Run() error {
	const op = "UserService.Run"

	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.port),
	)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Info("starting gRPC server")

	if err = a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (a *App) Stop() {
	const op = "UserService.Stop"

	a.log.With(slog.String("op", op)).Info("stopping gRPC server")

	a.gRPCServer.GracefulStop()
}
