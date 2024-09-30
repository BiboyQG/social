package store

import (
	"context"
	"database/sql"
)

type Storage struct {
	Posts interface {
		Create(ctx context.Context, post *Post) error
	}
	Users interface {
		Create(ctx context.Context, user *User) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts: &PostStore{db: db},
		Users: &UserStore{db: db},
	}
}
