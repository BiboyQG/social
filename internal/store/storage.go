package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNoRecord = errors.New("resource not found")
	QueryTimeoutDuration = 5 * time.Second
)

type Storage struct {
	Posts interface {
		Create(ctx context.Context, post *Post) error
		GetByID(ctx context.Context, id int64) (*Post, error)
		Update(ctx context.Context, post *Post) error
		Delete(ctx context.Context, id int64) error
	}
	Users interface {
		Create(ctx context.Context, user *User) error
		GetByID(ctx context.Context, id int64) (*User, error)
	}
	Comments interface {
		Create(ctx context.Context, comment *Comment) error
		GetByPostID(ctx context.Context, postID int64) ([]Comment, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostStore{db: db},
		Users:    &UserStore{db: db},
		Comments: &CommentStore{db: db},
	}
}
