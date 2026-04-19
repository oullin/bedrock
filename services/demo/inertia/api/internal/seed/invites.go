package seed

import (
	"database/sql"
	"time"

	"github.com/bedrock/services/demo/inertia/api/internal/database"
)

func seedInvites(db *sql.DB, now time.Time) error {
	invites := []struct {
		id, name, email, role, status string
		ago                           time.Duration
	}{
		{"invite_100", "Aria Lim", "aria@northstarhq.test", "Operator", "Accepted", 9 * time.Minute},
		{"invite_101", "Noah Chen", "noah@northstarhq.test", "Manager", "Pending", 41 * time.Minute},
		{"invite_102", "Priya Kapoor", "priya@northstarhq.test", "Executive", "Accepted", 1 * time.Hour},
		{"invite_103", "Leo Ruiz", "leo@northstarhq.test", "Analyst", "Pending", 2 * time.Hour},
		{"invite_104", "Sana Okafor", "sana@northstarhq.test", "Operator", "Accepted", 3 * time.Hour},
		{"invite_105", "Kai Tanaka", "kai@northstarhq.test", "Manager", "Expired", 5 * time.Hour},
		{"invite_106", "Dani Alves", "dani@northstarhq.test", "Analyst", "Accepted", 6 * time.Hour},
		{"invite_107", "Remi Dubois", "remi@northstarhq.test", "Operator", "Pending", 8 * time.Hour},
		{"invite_108", "Ines Moreno", "ines@northstarhq.test", "Executive", "Pending", 24 * time.Hour},
		{"invite_109", "Jules Park", "jules@northstarhq.test", "Manager", "Accepted", 48 * time.Hour},
		{"invite_110", "Maya Tan", "maya@northstarhq.test", "Executive", "Accepted", 3 * 24 * time.Hour},
		{"invite_111", "Ava Gomez", "ava@northstarhq.test", "Manager", "Accepted", 4 * 24 * time.Hour},
		{"invite_112", "Zara Hussain", "zara@northstarhq.test", "Analyst", "Expired", 5 * 24 * time.Hour},
		{"invite_113", "Tomás Reyes", "tomas@northstarhq.test", "Operator", "Accepted", 6 * 24 * time.Hour},
		{"invite_114", "Mila Novak", "mila@northstarhq.test", "Executive", "Pending", 7 * 24 * time.Hour},
		{"invite_115", "Ethan Osei", "ethan@northstarhq.test", "Analyst", "Accepted", 8 * 24 * time.Hour},
		{"invite_116", "Lena Vogt", "lena@northstarhq.test", "Operator", "Expired", 10 * 24 * time.Hour},
		{"invite_117", "Rohan Mehta", "rohan@northstarhq.test", "Manager", "Pending", 12 * 24 * time.Hour},
		{"invite_118", "Cleo Andersen", "cleo@northstarhq.test", "Analyst", "Accepted", 14 * 24 * time.Hour},
		{"invite_119", "Felix Braun", "felix@northstarhq.test", "Executive", "Accepted", 18 * 24 * time.Hour},
	}

	for _, inv := range invites {
		if err := database.CreateInviteAt(db, inv.id, inv.name, inv.email, inv.role, inv.status, now.Add(-inv.ago)); err != nil {
			return err
		}
	}

	return nil
}
