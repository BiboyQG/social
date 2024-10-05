package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNoRecord          = errors.New("resource not found")
	ErrAlreadyExists     = errors.New("resource already exists")
	QueryTimeoutDuration = 5 * time.Second
)

type Storage struct {
	Posts interface {
		Create(ctx context.Context, post *Post) error
		GetByID(ctx context.Context, id int64) (*Post, error)
		Update(ctx context.Context, post *Post) error
		Delete(ctx context.Context, id int64) error
		GetUserFeed(ctx context.Context, userID int64) ([]PostWithMetadata, error)
	}
	Users interface {
		Create(ctx context.Context, user *User) error
		GetByID(ctx context.Context, id int64) (*User, error)
	}
	Comments interface {
		Create(ctx context.Context, comment *Comment) error
		GetByPostID(ctx context.Context, postID int64) ([]Comment, error)
	}
	Followers interface {
		Follow(ctx context.Context, userID, followerID int64) error
		Unfollow(ctx context.Context, userID, followerID int64) error
		GetFollowers(ctx context.Context, userID int64) ([]User, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:     &PostStore{db: db},
		Users:     &UserStore{db: db},
		Comments:  &CommentStore{db: db},
		Followers: &FollowerStore{db: db},
	}
}
