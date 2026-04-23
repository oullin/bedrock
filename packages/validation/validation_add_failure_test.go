package validation_test

import (
	"testing"

	"github.com/bedrock/packages/validation"
)

// ValidationAddFailureTest::testAddFailureExists
func TestMessageBag_AddFailureExists(t *testing.T) {
	t.Parallel()

	b := validation.NewMessageBag()
	b.Add("email", "The email field is required.")

	if !b.Has("email") {
		t.Fatal("expected message bag to contain the failed field")
	}

	if got := b.First("email"); got != "The email field is required." {
		t.Fatalf("First(email) = %q, want %q", got, "The email field is required.")
	}
}

// ValidationAddFailureTest::testAddFailureIsFunctional
func TestMessageBag_AddFailureIsFunctional(t *testing.T) {
	t.Parallel()

	b := validation.NewMessageBag()
	b.Add("email", "The email field is required.")
	b.Add("email", "The email field is required.")
	b.Add("email", "The email field must be valid.")

	if got := b.Count(); got != 2 {
		t.Fatalf("Count() = %d, want 2", got)
	}

	msgs := b.Get("email")

	if len(msgs) != 2 {
		t.Fatalf("Get(email) = %v, want 2 messages", msgs)
	}

	if msgs[0] != "The email field is required." || msgs[1] != "The email field must be valid." {
		t.Fatalf("Get(email) = %v, want ordered unique messages", msgs)
	}
}
