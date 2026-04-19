package seed

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/bedrock/services/demo/inertia/api/internal/database"
)

func Run(db *sql.DB) error {
	if err := database.Truncate(db); err != nil {
		return fmt.Errorf("seed truncate: %w", err)
	}

	now := time.Now()

	if err := seedUsers(db, now); err != nil {
		return fmt.Errorf("seed users: %w", err)
	}

	orgIDs, err := seedOrganizations(db)

	if err != nil {
		return fmt.Errorf("seed organizations: %w", err)
	}

	if err := seedContacts(db, orgIDs); err != nil {
		return fmt.Errorf("seed contacts: %w", err)
	}

	if err := seedNotes(db, now); err != nil {
		return fmt.Errorf("seed notes: %w", err)
	}

	if err := seedInvites(db, now); err != nil {
		return fmt.Errorf("seed invites: %w", err)
	}

	if err := seedUploads(db, now); err != nil {
		return fmt.Errorf("seed uploads: %w", err)
	}

	if err := seedApprovals(db, now); err != nil {
		return fmt.Errorf("seed approvals: %w", err)
	}

	if err := database.SetCounter(db, "priority_escalations", 18); err != nil {
		return fmt.Errorf("seed counter: %w", err)
	}

	return nil
}
