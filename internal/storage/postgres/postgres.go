package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"
	"user_service/internal/service"
	"user_service/internal/storage"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Postgres {
	return &Postgres{
		db: db,
	}
}

func (p *Postgres) Get(ctx context.Context, username string) (service.User, error) {
	const op = "postgres.Get"
	sql := `SELECT user_id,username,
        COALESCE(first_name,''),
        COALESCE(last_name,''),
        COALESCE(birth_day,'0001-01-01'),
        COALESCE(phone_number,'') FROM profiles WHERE username = $1`

	var getDTO storage.GetDTO

	err := p.db.QueryRow(ctx, sql, username).Scan(&getDTO.Id, &getDTO.Username, &getDTO.FirstName, &getDTO.LastName, &getDTO.BirthDay, &getDTO.Phone)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return service.User{}, fmt.Errorf("%s:%w ", op, storage.ErrUserNotFound)
		}
		return service.User{}, fmt.Errorf("%s:%w ", op, err)
	}
	return service.User{
		Username:    getDTO.Username,
		Id:          getDTO.Id,
		LastName:    getDTO.LastName,
		FirstName:   getDTO.FirstName,
		Birthday:    getDTO.BirthDay,
		PhoneNumber: getDTO.Phone,
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

func (p *Postgres) Update(ctx context.Context, username string, firstname string, lastname string, birthDay time.Time, phone string) error {
	const op = "postgres.Update"

	sql := `
		UPDATE profiles 
		SET 
			first_name = CASE WHEN $1 <> '' THEN $1 ELSE first_name END,
			last_name = CASE WHEN $2 <> '' THEN $2 ELSE last_name END,
			phone_number = CASE WHEN $3 <> '' THEN $3 ELSE phone_number END,
			birth_day = CASE WHEN $4 > '0001-01-01'::date THEN $4 ELSE birth_day END
		WHERE username = $5`

	res, err := p.db.Exec(ctx, sql, firstname, lastname, phone, birthDay, username)
	if err != nil {
		return fmt.Errorf("%s:%w ", op, err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("%s:%w ", op, storage.ErrUserNotFound)
	}
	return nil
}
