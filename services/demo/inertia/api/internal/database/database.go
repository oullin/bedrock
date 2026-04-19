package database

import (
	"database/sql"
	"fmt"
	"math"
	"time"

	_ "modernc.org/sqlite"
)

func Open(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path+"?_pragma=foreign_keys(1)")

	if err != nil {
		return nil, fmt.Errorf("database open: %w", err)
	}

	db.SetMaxOpenConns(1)

	if err := db.Ping(); err != nil {
		db.Close()

		return nil, fmt.Errorf("database ping: %w", err)
	}

	if err := migrate(db); err != nil {
		db.Close()

		return nil, fmt.Errorf("database migrate: %w", err)
	}

	return db, nil
}

func migrate(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id          INTEGER PRIMARY KEY AUTOINCREMENT,
			name        TEXT NOT NULL,
			email       TEXT NOT NULL UNIQUE,
			password    TEXT NOT NULL,
			verified_at DATETIME,
			created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS organizations (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			name       TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS contacts (
			id              INTEGER PRIMARY KEY AUTOINCREMENT,
			organization_id INTEGER,
			first_name      TEXT NOT NULL,
			last_name       TEXT NOT NULL,
			email           TEXT NOT NULL,
			phone           TEXT NOT NULL DEFAULT '',
			is_favorite     INTEGER NOT NULL DEFAULT 0,
			created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (organization_id) REFERENCES organizations (id)
		);

		CREATE TABLE IF NOT EXISTS notes (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			contact_id INTEGER NOT NULL,
			user_id    INTEGER NOT NULL,
			body       TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (contact_id) REFERENCES contacts (id),
			FOREIGN KEY (user_id) REFERENCES users (id)
		);

		CREATE TABLE IF NOT EXISTS invites (
			id         TEXT PRIMARY KEY,
			name       TEXT NOT NULL,
			email      TEXT NOT NULL,
			role       TEXT NOT NULL,
			status     TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS uploads (
			id         TEXT PRIMARY KEY,
			label      TEXT NOT NULL,
			filename   TEXT NOT NULL,
			size       TEXT NOT NULL,
			status     TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS approvals (
			id         TEXT PRIMARY KEY,
			label      TEXT NOT NULL,
			status     TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS counters (
			key   TEXT PRIMARY KEY,
			value INTEGER NOT NULL DEFAULT 0
		);

		CREATE INDEX IF NOT EXISTS idx_contacts_organization_id ON contacts (organization_id);
		CREATE INDEX IF NOT EXISTS idx_notes_contact_id ON notes (contact_id);
		CREATE INDEX IF NOT EXISTS idx_notes_user_id ON notes (user_id);
	`)

	return err
}

func Truncate(db *sql.DB) error {
	_, err := db.Exec(`
		DELETE FROM notes;
		DELETE FROM contacts;
		DELETE FROM organizations;
		DELETE FROM users;
		DELETE FROM invites;
		DELETE FROM uploads;
		DELETE FROM approvals;
		DELETE FROM counters;
		DELETE FROM sqlite_sequence;
	`)

	return err
}

func ListInvites(db *sql.DB) []map[string]any {
	rows, err := db.Query("SELECT id, name, email, role, status, created_at FROM invites ORDER BY created_at DESC")

	if err != nil {
		return nil
	}

	defer rows.Close()

	return scanRows(rows, func(scan func(...any) error) (map[string]any, error) {
		var id, name, email, role, status string

		var createdAt time.Time

		if err := scan(&id, &name, &email, &role, &status, &createdAt); err != nil {
			return nil, err
		}

		return map[string]any{
			"id":     id,
			"name":   name,
			"email":  email,
			"role":   role,
			"status": status,
			"time":   TimeAgo(createdAt),
		}, nil
	})
}

func CreateInvite(db *sql.DB, id, name, email, role, status string) error {
	_, err := db.Exec(
		"INSERT INTO invites (id, name, email, role, status) VALUES (?, ?, ?, ?, ?)",
		id, name, email, role, status,
	)

	return err
}

func CreateInviteAt(db *sql.DB, id, name, email, role, status string, createdAt time.Time) error {
	_, err := db.Exec(
		"INSERT INTO invites (id, name, email, role, status, created_at) VALUES (?, ?, ?, ?, ?, ?)",
		id, name, email, role, status, createdAt,
	)

	return err
}

func InviteCount(db *sql.DB) int {
	var count int

	db.QueryRow("SELECT COUNT(*) FROM invites").Scan(&count)

	return count
}

func ListUploads(db *sql.DB) []map[string]any {
	rows, err := db.Query("SELECT id, label, filename, size, status, created_at FROM uploads ORDER BY created_at DESC")

	if err != nil {
		return nil
	}

	defer rows.Close()

	return scanRows(rows, func(scan func(...any) error) (map[string]any, error) {
		var id, label, filename, size, status string

		var createdAt time.Time

		if err := scan(&id, &label, &filename, &size, &status, &createdAt); err != nil {
			return nil, err
		}

		return map[string]any{
			"id":       id,
			"label":    label,
			"filename": filename,
			"size":     size,
			"status":   status,
			"time":     TimeAgo(createdAt),
		}, nil
	})
}

func CreateUpload(db *sql.DB, id, label, filename, size, status string) error {
	_, err := db.Exec(
		"INSERT INTO uploads (id, label, filename, size, status) VALUES (?, ?, ?, ?, ?)",
		id, label, filename, size, status,
	)

	return err
}

func CreateUploadAt(db *sql.DB, id, label, filename, size, status string, createdAt time.Time) error {
	_, err := db.Exec(
		"INSERT INTO uploads (id, label, filename, size, status, created_at) VALUES (?, ?, ?, ?, ?, ?)",
		id, label, filename, size, status, createdAt,
	)

	return err
}

func UploadCount(db *sql.DB) int {
	var count int

	db.QueryRow("SELECT COUNT(*) FROM uploads").Scan(&count)

	return count
}

func ListApprovals(db *sql.DB) []map[string]any {
	rows, err := db.Query("SELECT id, label, status, created_at FROM approvals ORDER BY created_at DESC")

	if err != nil {
		return nil
	}

	defer rows.Close()

	return scanRows(rows, func(scan func(...any) error) (map[string]any, error) {
		var id, label, status string

		var createdAt time.Time

		if err := scan(&id, &label, &status, &createdAt); err != nil {
			return nil, err
		}

		return map[string]any{
			"id":     id,
			"label":  label,
			"status": status,
			"time":   TimeAgo(createdAt),
		}, nil
	})
}

func CreateApproval(db *sql.DB, id, label, status string) error {
	_, err := db.Exec(
		"INSERT INTO approvals (id, label, status) VALUES (?, ?, ?)",
		id, label, status,
	)

	return err
}

func CreateApprovalAt(db *sql.DB, id, label, status string, createdAt time.Time) error {
	_, err := db.Exec(
		"INSERT INTO approvals (id, label, status, created_at) VALUES (?, ?, ?, ?)",
		id, label, status, createdAt,
	)

	return err
}

func ApprovalCount(db *sql.DB) int {
	var count int

	db.QueryRow("SELECT COUNT(*) FROM approvals").Scan(&count)

	return count
}

func GetCounter(db *sql.DB, key string) int {
	var value int

	db.QueryRow("SELECT value FROM counters WHERE key = ?", key).Scan(&value)

	return value
}

func IncrementCounter(db *sql.DB, key string) error {
	_, err := db.Exec(`
		INSERT INTO counters (key, value) VALUES (?, 1)
		ON CONFLICT(key) DO UPDATE SET value = value + 1
	`, key)

	return err
}

func SetCounter(db *sql.DB, key string, value int) error {
	_, err := db.Exec(`
		INSERT INTO counters (key, value) VALUES (?, ?)
		ON CONFLICT(key) DO UPDATE SET value = ?
	`, key, value, value)

	return err
}

func TimeAgo(t time.Time) string {
	d := time.Since(t)

	switch {
	case d < time.Minute:
		return "Just now"
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	default:
		days := int(math.Floor(d.Hours() / 24))

		if days == 1 {
			return "1d ago"
		}

		return fmt.Sprintf("%dd ago", days)
	}
}

func scanRows[T any](rows *sql.Rows, fn func(scan func(...any) error) (T, error)) []T {
	var result []T

	for rows.Next() {
		item, err := fn(rows.Scan)

		if err != nil {
			continue
		}

		result = append(result, item)
	}

	return result
}
