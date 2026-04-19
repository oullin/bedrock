package features

import (
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bedrock/packages/inertia/props"
	"github.com/bedrock/packages/inertia/protocol"
	"github.com/bedrock/services/demo/inertia/api/internal/database"
	"github.com/bedrock/services/demo/inertia/api/internal/httputil"
)

const contactsPerPage = 15

func (a app) deferredPropsHandler(w http.ResponseWriter, r *http.Request) {
	a.container.Render(w, r, "Features/DataLoading/DeferredProps", protocol.Props{
		"quickStat": 42,
		"slowStats": props.Defer(func() any {
			if httputil.SleepCtx(r.Context(), 800*time.Millisecond) != nil {
				return nil
			}

			totalContacts, err := database.CountContacts(a.container.DB)

			if err != nil {
				slog.Error("count contacts", "error", err)
			}

			return map[string]any{
				"totalContacts":  totalContacts,
				"totalFavorites": 12,
			}
		}, "slow"),
		"heavyData": props.Defer(func() any {
			if httputil.SleepCtx(r.Context(), 1500*time.Millisecond) != nil {
				return nil
			}

			items := make([]map[string]any, 20)

			for i := range items {
				items[i] = map[string]any{"id": i + 1, "name": fmt.Sprintf("Item %d", i+1)}
			}

			return items
		}, "heavy"),
	})
}

func (a app) partialReloadsHandler(w http.ResponseWriter, r *http.Request) {
	contacts, err := database.ListRecentContacts(a.container.DB, 5)

	if err != nil {
		slog.Error("list contacts", "error", err)

		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}

	users := make([]map[string]any, 0, len(contacts))

	for _, c := range contacts {
		users = append(users, map[string]any{
			"id":    c.ID,
			"name":  c.FirstName + " " + c.LastName,
			"email": c.Email,
		})
	}

	totalContacts, _ := database.CountContacts(a.container.DB)
	totalOrgs, _ := database.CountOrganizations(a.container.DB)

	a.container.Render(w, r, "Features/DataLoading/PartialReloads", protocol.Props{
		"users":        users,
		"stats":        map[string]any{"totalContacts": totalContacts, "totalOrganizations": totalOrgs},
		"timestamp":    time.Now().Format(time.RFC3339),
		"randomNumber": rand.Intn(1000),
	})
}

func (a app) infiniteScrollHandler(w http.ResponseWriter, r *http.Request) {
	var cursor *string

	if c := r.URL.Query().Get("cursor"); strings.TrimSpace(c) != "" {
		cursor = &c
	}

	page, err := database.ListContactsPaginated(a.container.DB, "", false, cursor, "next", contactsPerPage)

	if err != nil {
		slog.Error("list contacts paginated", "error", err)

		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}

	items := make([]map[string]any, 0, len(page.Data))

	for _, c := range page.Data {
		item := map[string]any{
			"id":          c.ID,
			"first_name":  c.FirstName,
			"last_name":   c.LastName,
			"email":       c.Email,
			"is_favorite": c.IsFavorite,
		}

		if c.OrganizationID != nil {
			item["organization"] = map[string]any{"name": c.OrganizationName}
		}

		items = append(items, item)
	}

	a.container.Render(w, r, "Features/DataLoading/InfiniteScroll", protocol.Props{
		"contacts": map[string]any{
			"data":        items,
			"next_cursor": page.NextCursor,
		},
	})
}

func (a app) whenVisibleHandler(w http.ResponseWriter, r *http.Request) {
	a.container.Render(w, r, "Features/DataLoading/WhenVisible", protocol.Props{
		"section1": props.Optional(func() any {
			if httputil.SleepCtx(r.Context(), 500*time.Millisecond) != nil {
				return nil
			}

			count, _ := database.CountContacts(a.container.DB)

			return map[string]any{"title": "Contacts", "content": fmt.Sprintf("%d contacts in the database.", count)}
		}),
		"section2": props.Optional(func() any {
			if httputil.SleepCtx(r.Context(), 800*time.Millisecond) != nil {
				return nil
			}

			count, _ := database.CountOrganizations(a.container.DB)

			return map[string]any{"title": "Organizations", "content": fmt.Sprintf("%d organizations tracked.", count)}
		}),
		"section3": props.Optional(func() any {
			if httputil.SleepCtx(r.Context(), 1000*time.Millisecond) != nil {
				return nil
			}

			count, _ := database.CountNotes(a.container.DB)

			return map[string]any{"title": "Notes", "content": fmt.Sprintf("%d notes recorded.", count)}
		}),
	})
}

func (a app) pollingHandler(w http.ResponseWriter, r *http.Request) {
	contactCount, err := database.CountContacts(a.container.DB)

	if err != nil {
		slog.Error("count contacts", "error", err)
	}

	a.container.Render(w, r, "Features/DataLoading/Polling", protocol.Props{
		"currentTime":  time.Now().Format(time.RFC3339),
		"randomNumber": rand.Intn(1000),
		"contactCount": contactCount,
	})
}

func (a app) propMergingHandler(w http.ResponseWriter, r *http.Request) {
	notes, err := database.ListRecentNotes(a.container.DB, 3)

	if err != nil {
		slog.Error("list recent notes", "error", err)
	}

	notifications := make([]map[string]any, 0, len(notes))

	for _, n := range notes {
		notifications = append(notifications, map[string]any{
			"id":    n.ID,
			"title": n.Body,
			"body":  n.ContactName,
		})
	}

	activityNotes, err := database.ListRecentNotes(a.container.DB, 2)

	if err != nil {
		slog.Error("list recent activity", "error", err)
	}

	activities := make([]map[string]any, 0, len(activityNotes))

	for _, n := range activityNotes {
		activities = append(activities, map[string]any{
			"id":          n.ID,
			"type":        "note",
			"description": n.Body,
		})
	}

	a.container.Render(w, r, "Features/DataLoading/PropMerging", protocol.Props{
		"notifications": props.Merge(notifications),
		"activities":    props.Merge(activities),
		"timestamp":     time.Now().Format(time.RFC3339),
	})
}

func (a app) optionalPropsHandler(w http.ResponseWriter, r *http.Request) {
	a.container.Render(w, r, "Features/DataLoading/OptionalProps", protocol.Props{
		"regularData": map[string]any{"message": "Always loaded"},
		"optionalData": props.Optional(func() any {
			return map[string]any{"message": "Loaded on demand"}
		}),
		"deferredData": props.Defer(func() any {
			if httputil.SleepCtx(r.Context(), 500*time.Millisecond) != nil {
				return nil
			}

			return map[string]any{"message": "Loaded asynchronously"}
		}),
	})
}

func (a app) oncePropsHandler(w http.ResponseWriter, r *http.Request) {
	pageNum, _ := strconv.Atoi(r.PathValue("page"))

	if pageNum < 1 {
		pageNum = 1
	}

	a.container.Render(w, r, "Features/DataLoading/OnceProps", protocol.Props{
		"page":       pageNum,
		"staticData": props.Once(map[string]any{"cached": true, "message": "This won't change on reload"}),
		"freshData":  map[string]any{"timestamp": time.Now().Format(time.RFC3339)},
		"dynamicData": map[string]any{
			"page":      pageNum,
			"generated": time.Now().Format(time.RFC3339),
		},
	})
}
