package seed

import (
	"database/sql"
	"time"

	"github.com/bedrock/services/demo/inertia/api/internal/database"
)

func seedApprovals(db *sql.DB, now time.Time) error {
	approvals := []struct {
		id, label, status string
		ago               time.Duration
	}{
		{"approval_s00", "Northwind priority routing", "Approved", 16 * time.Minute},
		{"approval_s01", "Atlas expansion guardrail", "Queued", 47 * time.Minute},
		{"approval_s02", "Cedar payment override", "Approved", 1 * time.Hour},
		{"approval_s03", "Helio discount authorization", "Synced", 2 * time.Hour},
		{"approval_s04", "Juniper compliance waiver", "Queued", 4 * time.Hour},
		{"approval_s05", "Mistral seat expansion", "Approved", 5 * time.Hour},
		{"approval_s06", "Summit invoice fast-track", "Synced", 8 * time.Hour},
		{"approval_s07", "Polar onboarding exception", "Queued", 24 * time.Hour},
		{"approval_s08", "Orchid renewal acceleration", "Approved", 2 * 24 * time.Hour},
		{"approval_s09", "Bluebird SLA amendment", "Synced", 3 * 24 * time.Hour},
		{"approval_s10", "Beacon routing promotion", "Approved", 4 * 24 * time.Hour},
		{"approval_s11", "Northwind contract extension", "Queued", 5 * 24 * time.Hour},
		{"approval_s12", "Atlas budget reallocation", "Approved", 7 * 24 * time.Hour},
		{"approval_s13", "Cedar autopay enrollment", "Synced", 10 * 24 * time.Hour},
		{"approval_s14", "Helio partner tier upgrade", "Queued", 14 * 24 * time.Hour},
	}

	for _, apr := range approvals {
		if err := database.CreateApprovalAt(db, apr.id, apr.label, apr.status, now.Add(-apr.ago)); err != nil {
			return err
		}
	}

	return nil
}
