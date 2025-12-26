package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"user_service/internal/app"
	"user_service/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

const (
	local = "local"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)

	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}

	conn, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		panic("Error connecting to database")
	}

	if err := conn.Ping(context.Background()); err != nil {
		panic("Error pinging database")
	}

	log.Info("Connected to database")

	aplication := app.NewApp(log, cfg.GRPC.Port, conn)

	go aplication.UserService.Run()
	log.Info("App started")

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Info("App stopping")
	aplication.UserService.Stop()
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case local:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	}
	return log
}
