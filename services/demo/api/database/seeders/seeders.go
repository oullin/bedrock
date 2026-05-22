package seeders

import (
	"database/sql"
	"fmt"
)

// Run inserts deterministic data matching the tiny upstream skeleton app.
func Run(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("demo seeders: nil database")
	}

	_, err := db.Exec(`
		INSERT INTO users (name, email, password, remember_token)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(email) DO UPDATE SET
			name = excluded.name,
			password = excluded.password,
			remember_token = excluded.remember_token,
			updated_at = CURRENT_TIMESTAMP
	`, "Taylor Otwell", "taylor@example.com", "$2y$12$bedrock.demo.password.hash", "demo-remember-token")

	if err != nil {
		return fmt.Errorf("demo seeders: %w", err)
	}

	return nil
}
