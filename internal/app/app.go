package app

import (
	"log/slog"
	"user_service/internal/app/userService"
	"user_service/internal/service"
	"user_service/internal/storage/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	UserService *userApp.App
}

func NewApp(
	log *slog.Logger,
	grpcPort int,
	db *pgxpool.Pool,
) *App {
	storage := postgres.New(db)

	userService := service.New(storage, log)

	grpcApp := userApp.New(log, grpcPort, userService)

	return &App{
		UserService: grpcApp,
	}
}
