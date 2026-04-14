package validation_test

import (
	"testing"

	"github.com/bedrock/packages/validation"
)

func TestMessageBag_Add(t *testing.T) {
	t.Parallel()

	b := validation.NewMessageBag()
	b.Add("email", "required")
	b.Add("email", "must be valid")
	b.Add("name", "required")

	if b.Count() != 3 {
		t.Fatalf("Count: got %d, want 3", b.Count())
	}
}

func TestMessageBag_AddDeduplicates(t *testing.T) {
	t.Parallel()

	b := validation.NewMessageBag()
	b.Add("email", "required")
	b.Add("email", "required") // duplicate

	if b.Count() != 1 {
		t.Fatalf("Count after duplicate Add: got %d, want 1", b.Count())
	}
}

func TestMessageBag_Has(t *testing.T) {
	t.Parallel()

	b := validation.NewMessageBag()
	b.Add("email", "required")

	if !b.Has("email") {
		t.Error("Has(email): want true")
	}

	if b.Has("name") {
		t.Error("Has(name): want false")
	}
}

func TestMessageBag_Get(t *testing.T) {
	t.Parallel()

	b := validation.NewMessageBag()
	b.Add("email", "required")
	b.Add("email", "invalid")

	msgs := b.Get("email")
	if len(msgs) != 2 {
		t.Fatalf("Get: got %d messages, want 2", len(msgs))
	}
}

func TestMessageBag_First(t *testing.T) {
	t.Parallel()

	b := validation.NewMessageBag()
	b.Add("email", "required")
	b.Add("email", "invalid")

	if got := b.First("email"); got != "required" {
		t.Errorf("First(email): got %q, want %q", got, "required")
	}

	if got := b.First(); got != "required" {
		t.Errorf("First(): got %q, want %q", got, "required")
	}
}

func TestMessageBag_All(t *testing.T) {
	t.Parallel()

	b := validation.NewMessageBag()
	b.Add("age", "min")
	b.Add("email", "required")

	all := b.All()

	// Keys are sorted: age comes before email
	if len(all) != 2 {
		t.Fatalf("All: got %d, want 2", len(all))
	}

	if all[0] != "min" {
		t.Errorf("All[0]: got %q, want %q", all[0], "min")
	}
}

func TestMessageBag_Merge(t *testing.T) {
	t.Parallel()

	b := validation.NewMessageBag()
	b.Add("email", "required")

	b.Merge(map[string][]string{
		"email": {"invalid"},
		"name":  {"required"},
	})

	if b.Count() != 3 {
		t.Fatalf("Count after Merge: got %d, want 3", b.Count())
	}
}

func TestMessageBag_IsEmpty(t *testing.T) {
	t.Parallel()

	b := validation.NewMessageBag()

	if !b.IsEmpty() {
		t.Error("IsEmpty on new bag: want true")
	}

	b.Add("x", "msg")

	if b.IsEmpty() {
		t.Error("IsEmpty after Add: want false")
	}
}

func TestMessageBag_ToMap(t *testing.T) {
	t.Parallel()

	b := validation.NewMessageBag()
	b.Add("email", "required")

	m := b.ToMap()
	if len(m["email"]) != 1 || m["email"][0] != "required" {
		t.Errorf("ToMap: got %v", m)
	}
}

func TestMessageBag_ToJSON(t *testing.T) {
	t.Parallel()

	b := validation.NewMessageBag()
	b.Add("email", "required")

	data, err := b.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON: %v", err)
	}

	if string(data) != `{"email":["required"]}` {
		t.Errorf("ToJSON: got %s", data)
	}
}
