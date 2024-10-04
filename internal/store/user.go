package store

import (
	"context"
	"database/sql"
	"errors"
)

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	CreatedAt string `json:"created_at"`
}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (username, password, email)
		VALUES ($1, $2, $3) RETURNING id, created_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	
	row := s.db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Password,
		user.Email,
	)
	err := row.Scan(
		&user.ID,
		&user.CreatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserStore) GetByID(ctx context.Context, id int64) (*User, error) {
	query := `
		SELECT id, username, email, created_at, password
		FROM users
		WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	row := s.db.QueryRowContext(ctx, query, id)
	user := &User{}
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.Password,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	return user, nil
}

