package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
	"user_service/internal/storage"
)

type User struct {
	Username    string
	PhoneNumber string
	Birthday    time.Time
	FirstName   string
	LastName    string
	Id          int64
}

type Storage interface {
	Get(ctx context.Context, username string) (User, error)
	New(ctx context.Context, username string, id int64) error
}

type UserService struct {
	storage Storage
	log     *slog.Logger
}

func (u *UserService) NewUser(ctx context.Context, id int64, username string) error {
	const op = "User.NewUser"
	log := u.log.With(slog.String("op", op))

	log.Info("creating new user")
	if err := u.NewUser(ctx, id, username); err != nil {
		if errors.Is(err, storage.ErrUserAlreadyExists) {
			log.Error("user already exists")
			return fmt.Errorf("%s:%w", op, storage.ErrUserAlreadyExists)
		}
		log.Error("failed to create new user")
		return fmt.Errorf("%s:%w", op, err)
	}
	return nil
}

func (u *UserService) GetUser(ctx context.Context, username string) (User, error) {
	const op = "User.GetUser"
	log := u.log.With(slog.String("op", op))

	log.Info("getting user")
	user, err := u.storage.Get(ctx, username)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Error("user not found")
			return User{}, fmt.Errorf("%s:%w", op, storage.ErrUserNotFound)
		}
		log.Error("failed to get user")
		return User{}, fmt.Errorf("%s:%w", op, err)
	}
	return user, nil
}
