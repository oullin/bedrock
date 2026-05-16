package workflow_test

import (
	"errors"
	"testing"

	"github.com/bedrock/packages/workflow"
	"github.com/bedrock/packages/workflow/audit"
	"github.com/bedrock/packages/workflow/events"
	"github.com/bedrock/packages/workflow/registry"
	"github.com/bedrock/packages/workflow/store"
)

// Subscription is the canonical state-machine test subject: trial -> active -> cancelled.
type Subscription struct {
	ID    string
	State string
}

func subscriptionDef(t *testing.T) *workflow.Definition {
	t.Helper()

	def, err := workflow.NewDefinitionBuilder().
		AddPlace("trial").
		AddPlace("active").
		AddPlace("cancelled").
		SetInitialPlaces("trial").
		AddTransition("activate", []string{"trial"}, []string{"active"}).
		AddTransition("cancel", []string{"active"}, []string{"cancelled"}).
		Build()

	if err != nil {
		t.Fatalf("build definition: %v", err)
	}

	return def
}

func subscriptionStore() *store.SingleState[*Subscription] {
	return &store.SingleState[*Subscription]{
		Getter: func(s *Subscription) string { return s.State },
		Setter: func(s *Subscription, place string) { s.State = place },
	}
}

func TestStateMachine_ActivateAdvancesMarking(t *testing.T) {
	def := subscriptionDef(t)

	dispatcher := events.NewDispatcher[*Subscription]()

	var enteredCount int

	dispatcher.On(workflow.EventNameEntered("subscription"), func(events.Event[*Subscription]) {
		enteredCount++
	})

	sm, err := workflow.NewStateMachine("subscription", def, subscriptionStore(), dispatcher)

	if err != nil {
		t.Fatalf("new state machine: %v", err)
	}

	sub := &Subscription{ID: "s1", State: "trial"}

	marking, err := sm.Apply(sub, "activate", nil)

	if err != nil {
		t.Fatalf("apply activate: %v", err)
	}

	if !marking.Has("active") {
		t.Fatalf("expected active place, got %v", marking.ActivePlaces())
	}

	if sub.State != "active" {
		t.Fatalf("expected subject state 'active', got %q", sub.State)
	}

	if enteredCount == 0 {
		t.Fatal("expected at least one Entered event")
	}
}

func TestStateMachine_GuardVeto(t *testing.T) {
	def := subscriptionDef(t)

	dispatcher := events.NewDispatcher[*Subscription]()

	dispatcher.OnGuard(workflow.EventNameGuardNamed("subscription", "activate"), func(g *events.GuardEvent[*Subscription]) {
		g.SetBlocked(true, "billing not configured")
	})

	sm, _ := workflow.NewStateMachine("subscription", def, subscriptionStore(), dispatcher)

	sub := &Subscription{ID: "s2", State: "trial"}

	_, err := sm.Apply(sub, "activate", nil)

	if err == nil {
		t.Fatal("expected guard veto to return error")
	}

	var te *workflow.TransitionError

	if !errors.As(err, &te) {
		t.Fatalf("expected *TransitionError, got %T: %v", err, err)
	}

	if len(te.Blockers) == 0 || te.Blockers[0].Message != "billing not configured" {
		t.Fatalf("expected blocker message, got %#v", te.Blockers)
	}

	if sub.State != "trial" {
		t.Fatalf("subject state should be unchanged on veto, got %q", sub.State)
	}
}

func TestStateMachine_UnknownTransition(t *testing.T) {
	def := subscriptionDef(t)

	sm, _ := workflow.NewStateMachine("subscription", def, subscriptionStore(), nil)

	sub := &Subscription{State: "trial"}

	_, err := sm.Apply(sub, "bogus", nil)

	if !errors.Is(err, workflow.ErrTransitionNotFound) {
		t.Fatalf("expected ErrTransitionNotFound, got %v", err)
	}
}

func TestRegistry_LookupBySubjectAndName(t *testing.T) {
	def := subscriptionDef(t)
	sm, _ := workflow.NewStateMachine("subscription", def, subscriptionStore(), nil)

	reg := registry.New[*Subscription]()

	reg.Add(registry.Entry[*Subscription]{
		Name:     "subscription",
		Workflow: sm,
		Supports: func(s *Subscription) bool { return s != nil },
	})

	got, err := reg.Get(&Subscription{ID: "x"}, "subscription")

	if err != nil {
		t.Fatalf("registry get: %v", err)
	}

	if got.Name() != "subscription" {
		t.Fatalf("expected subscription workflow, got %q", got.Name())
	}

	if _, err := reg.Get(&Subscription{}, "nope"); err == nil {
		t.Fatal("expected error for unknown workflow name")
	}
}

func TestAudit_TrailRecordsCompletedTransitions(t *testing.T) {
	def := subscriptionDef(t)

	dispatcher := events.NewDispatcher[*Subscription]()

	trail := &audit.Trail[*Subscription]{}
	trail.Attach("subscription", dispatcher)

	sm, _ := workflow.NewStateMachine("subscription", def, subscriptionStore(), dispatcher)

	sub := &Subscription{ID: "s3", State: "trial"}

	if _, err := sm.Apply(sub, "activate", map[string]any{"reason": "trial-end"}); err != nil {
		t.Fatalf("activate: %v", err)
	}

	if _, err := sm.Apply(sub, "cancel", nil); err != nil {
		t.Fatalf("cancel: %v", err)
	}

	if len(trail.Entries) != 2 {
		t.Fatalf("expected 2 trail entries, got %d", len(trail.Entries))
	}

	if trail.Entries[0].Transition != "activate" {
		t.Fatalf("entry[0].Transition = %q, want activate", trail.Entries[0].Transition)
	}

	if trail.Entries[1].Transition != "cancel" {
		t.Fatalf("entry[1].Transition = %q, want cancel", trail.Entries[1].Transition)
	}

	if got := trail.Entries[0].Context["reason"]; got != "trial-end" {
		t.Fatalf("expected context reason captured, got %v", got)
	}
}

func TestRegistry_ConcurrentAddAndGet(t *testing.T) {
	def := subscriptionDef(t)
	sm, _ := workflow.NewStateMachine("subscription", def, subscriptionStore(), nil)

	reg := registry.New[*Subscription]()
	done := make(chan struct{})

	go func() {
		for i := 0; i < 100; i++ {
			reg.Add(registry.Entry[*Subscription]{
				Name:     "subscription",
				Workflow: sm,
				Supports: func(*Subscription) bool { return true },
			})
		}

		close(done)
	}()

	for i := 0; i < 100; i++ {
		reg.Get(&Subscription{}, "subscription")
	}

	<-done
}
