package crm

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/bedrock/packages/inertia/props"
	"github.com/bedrock/packages/inertia/protocol"
	"github.com/bedrock/services/inertia-demo/api/internal/httputil"
)

func (a app) dashboardHandler(w http.ResponseWriter, r *http.Request) {
	activity, err := a.repo.ListRecentNotes(10)

	if err != nil {
		slog.Error("recent activity", "error", err)

		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}

	a.container.Render(w, r, "Crm/Dashboard", protocol.Props{
		"recentActivity": recentActivityProps(activity),
		"totalContacts": props.Defer(func() any {
			if httputil.SleepCtx(r.Context(), 150*time.Millisecond) != nil {
				return nil
			}

			n, err := a.repo.CountContacts()

			if err != nil {
				slog.Error("count contacts", "error", err)

				return nil
			}

			return n
		}, "stats"),
		"totalOrganizations": props.Defer(func() any {
			if httputil.SleepCtx(r.Context(), 150*time.Millisecond) != nil {
				return nil
			}

			n, err := a.repo.CountOrganizations()

			if err != nil {
				slog.Error("count organizations", "error", err)

				return nil
			}

			return n
		}, "stats"),
		"recentNotesCount": props.Defer(func() any {
			if httputil.SleepCtx(r.Context(), 150*time.Millisecond) != nil {
				return nil
			}

			n, err := a.repo.CountNotes()

			if err != nil {
				slog.Error("count notes", "error", err)

				return nil
			}

			return n
		}, "stats"),
	})
}
