package seed

import (
	"database/sql"
	"time"

	"github.com/bedrock/services/demo/inertia/api/internal/database"
)

func seedUsers(db *sql.DB, now time.Time) error {
	users := []struct {
		name, email string
	}{
		{"Demo User", "test@example.com"},
		{"Alice Manager", "alice@example.com"},
		{"Bob Analyst", "bob@example.com"},
		{"Carol Reviewer", "carol@example.com"},
	}

	for _, u := range users {
		if _, err := database.CreateUser(db, u.name, u.email, "password", &now); err != nil {
			return err
		}
	}

	return nil
}
