package postgres

import (
	"context"
	"errors"
	"fmt"
	"user_service/internal/service"
	"user_service/internal/storage"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	db *pgxpool.Pool
}

func (p *Postgres) Get(ctx context.Context, username string) (service.User, error) {
	const op = "postgres.Get"
	sql := "SELECT (user_id,username,first_name,last_name,birth_day,phone_number) FROM profiles WHERE username = $1"

	var getDTO storage.GetDTO

	err := p.db.QueryRow(ctx, sql, username).Scan(&getDTO.Id, &getDTO.Username, &getDTO.FirstName, &getDTO.LastName, &getDTO.BirthDay, &getDTO.Phone)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "02000" {
			return service.User{}, fmt.Errorf("%s:%w ", op, storage.ErrUserNotFound)
		}
		return service.User{}, fmt.Errorf("%s:%w ", op, err)
	}
	return service.User{
		Username:    getDTO.Username,
		Id:          getDTO.Id,
		LastName:    getDTO.LastName.String,
		FirstName:   getDTO.FirstName.String,
		Birthday:    getDTO.BirthDay.Time,
		PhoneNumber: getDTO.Phone.String,
	}, nil
}

func (p *Postgres) New(ctx context.Context, username string, id int64) error {
	const op = "postgres.New"
	sql := "INSERT INTO profiles (user_id,username) VALUES ($1,$2)"

	if _, err := p.db.Exec(ctx, sql, id, username); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("%s:%w", op, storage.ErrUserAlreadyExists)
		}
		return fmt.Errorf("%s:%w", op, err)
	}
	return nil
}
