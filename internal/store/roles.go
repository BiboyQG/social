package store

import (
	"context"
	"database/sql"
	"errors"
)

type Role struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Level       int    `json:"level"`
}

type RoleStore struct {
	db *sql.DB
}

func (s *RoleStore) GetByName(ctx context.Context, name string) (*Role, error) {
	query := `SELECT id, name, description, level FROM roles WHERE name = $1`
	row := s.db.QueryRowContext(ctx, query, name)

	role := &Role{}
	err := row.Scan(&role.ID, &role.Name, &role.Description, &role.Level)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("role not found")
		}
		return nil, err
	}
	return role, nil
}
