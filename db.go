package main

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

// OpenTaskStore opens the SQLite database and ensures its schema exists.
func OpenTaskStore(path string) (*TaskStore, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	store := &TaskStore{db: db}
	if err := store.createSchema(); err != nil {
		db.Close()
		return nil, err
	}

	return store, nil

}

func (s *TaskStore) createSchema() error {
	const schema = `
		CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY,
			title TEXT NOT NULL,
			done INTEGER NOT NULL DEFAULT 0
		);
	`

	if _, err := s.db.Exec(schema); err != nil {
		return fmt.Errorf("create tasks table: %w", err)
	}

	return nil
}

// Close releases the database connection.
func (s *TaskStore) Close() error {
	return s.db.Close()
}
