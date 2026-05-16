package config_test

import (
	"testing"

	bconfig "github.com/bedrock/packages/config"
	"github.com/bedrock/packages/workflow"
	"github.com/bedrock/packages/workflow/config"
	"github.com/bedrock/packages/workflow/events"
	"github.com/bedrock/packages/workflow/store"
)

type Subscription struct{ State string }

func TestLoad_BuildsDefinitionFromRepository(t *testing.T) {
	repo := bconfig.New(map[string]any{
		"workflow": map[string]any{
			"name":    "subscription",
			"places":  []any{"trial", "active", "cancelled"},
			"initial": []any{"trial"},
			"transitions": []any{
				map[string]any{"name": "activate", "from": []any{"trial"}, "to": []any{"active"}},
				map[string]any{"name": "cancel", "from": []any{"active"}, "to": []any{"cancelled"}},
			},
			"metadata": map[string]any{
				"purpose": "billing",
			},
			"transitions_metadata": map[string]any{
				"activate": map[string]any{"audit_level": "high"},
			},
		},
	})

	def, err := config.Load(repo)

	if err != nil {
		t.Fatalf("load: %v", err)
	}

	if len(def.Places) != 3 {
		t.Fatalf("expected 3 places, got %d", len(def.Places))
	}

	if got, _ := def.MetadataValue("purpose"); got != "billing" {
		t.Fatalf("metadata purpose = %v, want billing", got)
	}

	if got, _ := def.TransitionMetadataValue("activate", "audit_level"); got != "high" {
		t.Fatalf("transition metadata = %v, want high", got)
	}

	// Make sure the produced Definition drives Apply.
	sm, err := workflow.NewStateMachine("subscription", def, &store.SingleState[*Subscription]{
		Getter: func(s *Subscription) string { return s.State },
		Setter: func(s *Subscription, place string) { s.State = place },
	}, events.NewDispatcher[*Subscription]())

	if err != nil {
		t.Fatalf("new state machine: %v", err)
	}

	sub := &Subscription{State: "trial"}

	if _, err := sm.Apply(sub, "activate", nil); err != nil {
		t.Fatalf("apply activate: %v", err)
	}

	if sub.State != "active" {
		t.Fatalf("expected state active, got %q", sub.State)
	}
}

func TestLoad_MissingPlacesReturnsError(t *testing.T) {
	repo := bconfig.New(map[string]any{
		"workflow": map[string]any{"name": "broken"},
	})

	if _, err := config.Load(repo); err == nil {
		t.Fatal("expected error for missing places")
	}
}
