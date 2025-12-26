package storage

import (
	"database/sql"
	"errors"
)

type GetDTO struct {
	Id        int64
	Username  string
	FirstName sql.NullString // Може бути NULL
	LastName  sql.NullString // Може бути NULL
	BirthDay  sql.NullTime
	Phone     sql.NullString
}

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)
