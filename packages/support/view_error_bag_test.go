package support

import "testing"

// Exact inventory markers covered by the executable tests in this file:
// SupportViewErrorBagTest::testHasBagTrue
// SupportViewErrorBagTest::testHasBagFalse
// SupportViewErrorBagTest::testGet
// SupportViewErrorBagTest::testGetBagWithNew
// SupportViewErrorBagTest::testGetBags
// SupportViewErrorBagTest::testPut
// SupportViewErrorBagTest::testAnyTrue
// SupportViewErrorBagTest::testAnyFalse
// SupportViewErrorBagTest::testAnyFalseWithEmptyErrorBag
// SupportViewErrorBagTest::testCount
// SupportViewErrorBagTest::testCountWithNoMessagesInMessageBag
// SupportViewErrorBagTest::testCountWithNoMessageBags
// SupportViewErrorBagTest::testToString

func TestViewErrorBagStorage(t *testing.T) {
	t.Parallel()

	errors := NewViewErrorBag()
	bag := NewMessageBag(map[string][]string{"email": {"required"}})

	if errors.HasBag("default") {
		t.Fatal("expected missing default bag")
	}

	errors.Put("default", bag)

	if !errors.HasBag("default") {
		t.Fatal("expected default bag")
	}

	if errors.Get("default").First("email") != "required" {
		t.Fatalf("Get(default).First(email) = %q", errors.Get("default").First("email"))
	}

	if !errors.HasBag("new") {
		created := errors.Get("new")
		if created == nil || !errors.HasBag("new") {
			t.Fatal("expected Get to create missing bag")
		}
	}

	if len(errors.Bags()) != 2 {
		t.Fatalf("Bags count = %d", len(errors.Bags()))
	}
}

func TestViewErrorBagAnyCountAndString(t *testing.T) {
	t.Parallel()

	empty := NewViewErrorBag()
	if empty.Any() {
		t.Fatal("empty ViewErrorBag should not have any messages")
	}

	if empty.Count() != 0 {
		t.Fatalf("empty Count = %d", empty.Count())
	}

	withEmptyBag := NewViewErrorBag().Put("default", NewMessageBag())
	if withEmptyBag.Any() {
		t.Fatal("ViewErrorBag with empty MessageBag should not have any messages")
	}

	if withEmptyBag.Count() != 0 {
		t.Fatalf("empty MessageBag Count = %d", withEmptyBag.Count())
	}

	bag := NewMessageBag(map[string][]string{"email": {"required"}, "name": {"required"}})
	errors := NewViewErrorBag().Put("default", bag)

	if !errors.Any() {
		t.Fatal("expected messages")
	}

	if errors.Count() != 2 {
		t.Fatalf("Count = %d", errors.Count())
	}

	if errors.String() != bag.String() {
		t.Fatalf("String = %q want %q", errors.String(), bag.String())
	}
}
