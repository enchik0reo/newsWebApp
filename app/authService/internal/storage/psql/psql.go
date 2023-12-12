package psql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"newsWebApp/app/authService/internal/config"
	"newsWebApp/app/authService/internal/models"
	"newsWebApp/app/authService/internal/storage"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(cfg config.Postgres) (*Storage, error) {
	// dbdriver://username:password@host:port/dbname?param1=true&param2=false
	dsn := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Driver, cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)

	db, err := sql.Open(cfg.Driver, dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	s := Storage{db: db}

	return &s, nil
}

func (s *Storage) CloseConn() error {
	return s.db.Close()
}

func (s *Storage) SaveUser(ctx context.Context, email string, hashPass []byte) (int64, error) {
	stmt, err := s.db.Prepare("INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING user_id")
	if err != nil {
		return 0, fmt.Errorf("can't prepare statement: %w", err)
	}

	row := stmt.QueryRowContext(ctx, email, hashPass)

	if err := row.Err(); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, storage.ErrUserExists
		}
		return 0, fmt.Errorf("can't insert values: %w", err)
	}

	var id int64

	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("can't get last insert id: %w", err)
	}

	return id, nil
}

func (s *Storage) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	stmt, err := s.db.Prepare("SELECT user_id, email, password_hash FROM users WHERE email = $1")
	if err != nil {
		return nil, fmt.Errorf("can't prepare statement: %w", err)
	}

	row := stmt.QueryRowContext(ctx, email)

	if err := row.Err(); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrUserNotFound
		}
		return nil, fmt.Errorf("can't get user result: %w", err)
	}

	u := models.User{}

	if err := row.Scan(&u.ID, &u.Email, &u.PassHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrUserNotFound
		}
		return nil, fmt.Errorf("can't get result: %w", err)
	}

	return &u, nil
}
