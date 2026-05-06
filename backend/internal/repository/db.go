package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func Open(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path+"?_foreign_keys=on&_journal_mode=WAL")
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}
	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return db, nil
}

func migrate(db *sql.DB) error {
	schema := `
CREATE TABLE IF NOT EXISTS users (
  id         INTEGER PRIMARY KEY AUTOINCREMENT,
  name       TEXT    NOT NULL,
  email      TEXT    NOT NULL UNIQUE,
  password   TEXT    NOT NULL,
  role       TEXT    NOT NULL DEFAULT 'cliente',
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tables (
  id        INTEGER PRIMARY KEY AUTOINCREMENT,
  number    INTEGER NOT NULL,
  location  TEXT    NOT NULL,
  capacity  INTEGER NOT NULL,
  is_active INTEGER NOT NULL DEFAULT 1
);

CREATE TABLE IF NOT EXISTS availability (
  id           INTEGER PRIMARY KEY AUTOINCREMENT,
  date         TEXT    NOT NULL UNIQUE,
  is_open      INTEGER NOT NULL DEFAULT 1,
  auto_confirm INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS time_slots (
  id              INTEGER PRIMARY KEY AUTOINCREMENT,
  availability_id INTEGER NOT NULL REFERENCES availability(id) ON DELETE CASCADE,
  slot_time       TEXT    NOT NULL
);

CREATE TABLE IF NOT EXISTS table_exclusions (
  id       INTEGER PRIMARY KEY AUTOINCREMENT,
  table_id INTEGER NOT NULL REFERENCES tables(id) ON DELETE CASCADE,
  date     TEXT    NOT NULL,
  UNIQUE(table_id, date)
);

CREATE TABLE IF NOT EXISTS reservations (
  id         INTEGER  PRIMARY KEY AUTOINCREMENT,
  user_id    INTEGER  NOT NULL REFERENCES users(id),
  table_id   INTEGER  NOT NULL REFERENCES tables(id),
  slot_id    INTEGER  NOT NULL REFERENCES time_slots(id),
  date       TEXT     NOT NULL,
  party_size INTEGER  NOT NULL,
  status     TEXT     NOT NULL DEFAULT 'pending',
  notes      TEXT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);`
	_, err := db.Exec(schema)
	return err
}
