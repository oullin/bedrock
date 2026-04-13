package notifications_test

import (
	"context"
	"testing"

	"github.com/bedrock/packages/notifications"
)

func TestAnonymousNotifiableRoute(t *testing.T) {
	t.Parallel()

	n := notifications.NewAnonymousNotifiable().
		Route("mail", "test@example.com").
		Route("slack", "#general")

	ctx := context.Background()

	if route := n.RouteNotificationFor(ctx, "mail"); route != "test@example.com" {
		t.Fatalf("mail route = %v, want %q", route, "test@example.com")
	}

	if route := n.RouteNotificationFor(ctx, "slack"); route != "#general" {
		t.Fatalf("slack route = %v, want %q", route, "#general")
	}
}

func TestAnonymousNotifiableRouteIgnoresDatabase(t *testing.T) {
	t.Parallel()

	n := notifications.NewAnonymousNotifiable().
		Route("database", "some-value")

	ctx := context.Background()

	if route := n.RouteNotificationFor(ctx, "database"); route != nil {
		t.Fatalf("database route = %v, want nil", route)
	}
}

func TestAnonymousNotifiableRouteNotFoundReturnsNil(t *testing.T) {
	t.Parallel()

	n := notifications.NewAnonymousNotifiable()
	ctx := context.Background()

	if route := n.RouteNotificationFor(ctx, "mail"); route != nil {
		t.Fatalf("unknown route = %v, want nil", route)
	}
}

func TestAnonymousNotifiableGetKey(t *testing.T) {
	t.Parallel()

	n := notifications.NewAnonymousNotifiable()

	if key := n.GetKey(); key != "" {
		t.Fatalf("GetKey() = %q, want empty", key)
	}
}

func TestAnonymousNotifiableChaining(t *testing.T) {
	t.Parallel()

	n := notifications.NewAnonymousNotifiable().
		Route("mail", "a@b.com").
		Route("sms", "+1234567890")

	// Verify both routes exist.
	ctx := context.Background()

	if n.RouteNotificationFor(ctx, "mail") == nil {
		t.Fatal("expected mail route")
	}

	if n.RouteNotificationFor(ctx, "sms") == nil {
		t.Fatal("expected sms route")
	}
}
