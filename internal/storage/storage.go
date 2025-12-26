package storage

import (
	"errors"
	"time"
)

type GetDTO struct {
	Id        int64
	Username  string
	FirstName string
	LastName  string
	BirthDay  time.Time
	Phone     string
}

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)
