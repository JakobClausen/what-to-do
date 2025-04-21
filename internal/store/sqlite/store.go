package sqlite

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"

	"what-to-do/internal/domain"

	_ "github.com/glebarez/go-sqlite"
)

type Store struct{ db *sql.DB }

func New() (*Store, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	appDir := filepath.Join(homeDir, ".what-to-do", "store")
	if err = os.MkdirAll(appDir, 0755); err != nil {
		return nil, err
	}

	dbPath := filepath.Join(appDir, "store.sqlite")

	db, err := sql.Open("sqlite", dbPath)

	if err != nil {
		return nil, err
	}

	s := &Store{db: db}

	return &Store{db: db}, s.migrate()
}

func (s *Store) migrate() error {
	const ddl = `
CREATE TABLE IF NOT EXISTS commands (
  id           INTEGER PRIMARY KEY AUTOINCREMENT,
  alias        TEXT NOT NULL,
  description  TEXT,
  command      TEXT,
  spotlighted  BOOLEAN DEFAULT FALSE,
  created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`
	_, err := s.db.Exec(ddl)
	return err
}

func (s *Store) Close() error { return s.db.Close() }

func (s *Store) CreateCommand(ctx context.Context, c *domain.Command) (int64, error) {
	res, err := s.db.ExecContext(
		ctx,
		`INSERT INTO commands (alias, description, command, spotlighted)
		 VALUES (?, ?, ?, ?)`,
		c.Alias, c.Description, c.Command, c.Spotlighted,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *Store) GetCommand(ctx context.Context, id int64) (*domain.Command, error) {
	var c domain.Command
	err := s.db.QueryRowContext(
		ctx,
		`SELECT id, alias, description, command, spotlighted,
		        created_at, updated_at
		   FROM commands WHERE id = ?`, id,
	).Scan(
		&c.ID, &c.Alias, &c.Description, &c.Command, &c.Spotlighted,
		&c.CreatedAt, &c.UpdatedAt,
	)
	return &c, err
}

func (s *Store) ListCommands(ctx context.Context) ([]*domain.Command, error) {
	rows, err := s.db.QueryContext(
		ctx,
		`SELECT id, alias, description, command, spotlighted,
		        created_at, updated_at
		   FROM commands
		  ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cmds []*domain.Command
	for rows.Next() {
		var c domain.Command
		if err := rows.Scan(
			&c.ID, &c.Alias, &c.Description, &c.Command, &c.Spotlighted,
			&c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, err
		}
		cmds = append(cmds, &c)
	}
	return cmds, rows.Err()
}

func (s *Store) UpdateCommand(ctx context.Context, c *domain.Command) error {
	_, err := s.db.ExecContext(
		ctx,
		`UPDATE commands SET alias = ?, description = ?, command = ?, spotlighted = ?, updated_at = CURRENT_TIMESTAMP
		  WHERE id = ?`,
		c.Alias, c.Description, c.Command, c.Spotlighted, c.ID,
	)
	return err
}

func (s *Store) DeleteCommand(ctx context.Context, id int64) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM commands WHERE id = ?`, id)
	return err
}
