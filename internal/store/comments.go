package store

import (
	"context"
	"database/sql"
)

type Comment struct {
	ID        int64    `json:"id"`
	PostID    int64    `json:"post_id"`
	UserID    int64    `json:"user_id"`
	Content   string   `json:"content"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
	User      User     `json:"user"`
}

type CommentStore struct {
	db *sql.DB
}

func (s *CommentStore) Create(ctx context.Context, comment *Comment) (int64, error) {
	query := `
		INSERT INTO comments (post_id, user_id, content)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var id int64

	err := s.db.QueryRowContext(ctx, query, comment.PostID, comment.UserID, comment.Content).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *CommentStore) GetByPostID(ctx context.Context, postID int64) ([]Comment, error) {
	query := `
		SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, c.updated_at, u.id, u.username, u.email, u.created_at
		FROM comments c
		LEFT JOIN users u ON c.user_id = u.id
		WHERE c.post_id = $1
		ORDER BY c.created_at DESC
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment

	for rows.Next() {
		var comment Comment
		comment.User = User{}
		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt, &comment.UpdatedAt, &comment.User.ID, &comment.User.Username, &comment.User.Email, &comment.User.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}





