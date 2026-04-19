package seed

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/bedrock/services/demo/inertia/api/internal/database"
)

func seedUploads(db *sql.DB, now time.Time) error {
	uploads := []struct {
		id, label, filename, size, status string
		ago                               time.Duration
	}{
		{"upload_s00", "March collections export", "collections-march.csv", "84 KB", "Processed", 12 * time.Minute},
		{"upload_s01", "Sales handoff packet", "handoff-q2.pdf", "1.2 MB", "Ready", 2 * time.Hour},
		{"upload_s02", "Q1 revenue reconciliation", "revenue-q1.xlsx", "320 KB", "Processed", 3 * time.Hour},
		{"upload_s03", "Partner channel attribution", "channel-attr.csv", "156 KB", "Queued", 4 * time.Hour},
		{"upload_s04", "Northwind contract redlines", "northwind-redlines.pdf", "2.4 MB", "Processed", 6 * time.Hour},
		{"upload_s05", "Atlas compliance artifacts", "atlas-compliance.zip", "4.1 MB", "Ready", 8 * time.Hour},
		{"upload_s06", "Cedar renewal terms", "cedar-renewal.docx", "92 KB", "Queued", 24 * time.Hour},
		{"upload_s07", "Helio onboarding checklist", "helio-onboarding.pdf", "540 KB", "Processed", 48 * time.Hour},
		{"upload_s08", "Juniper SLA benchmarks", "juniper-sla.xlsx", "210 KB", "Ready", 3 * 24 * time.Hour},
		{"upload_s09", "Mistral pricing model", "mistral-pricing.json", "18 KB", "Processed", 4 * 24 * time.Hour},
		{"upload_s10", "Summit invoice batch", "summit-invoices-q1.csv", "1.8 MB", "Queued", 5 * 24 * time.Hour},
		{"upload_s11", "Polar integration spec", "polar-api-spec.pdf", "3.2 MB", "Ready", 6 * 24 * time.Hour},
		{"upload_s12", "Orchid deal room export", "orchid-deal-room.zip", "7.6 MB", "Processed", 8 * 24 * time.Hour},
		{"upload_s13", "Bluebird NDA packet", "bluebird-nda.pdf", "420 KB", "Queued", 10 * 24 * time.Hour},
		{"upload_s14", "Beacon ops runbook", "beacon-runbook.docx", "1.1 MB", "Processed", 14 * 24 * time.Hour},
	}

	for _, upl := range uploads {
		if err := database.CreateUploadAt(db, upl.id, upl.label, upl.filename, upl.size, upl.status, now.Add(-upl.ago)); err != nil {
			return fmt.Errorf("upload %s: %w", upl.id, err)
		}
	}

	return nil
}
