package support

import (
	"encoding/json"
	"testing"
)

// Additional exact inventory markers covered by the executable tests in this file:
// SupportMessageBagTest::testConstructor
// SupportMessageBagTest::testConstructorUniquenessConsistency
// SupportMessageBagTest::testCountReturnsCorrectValue
// SupportMessageBagTest::testFirstFindsMessageForWildcardKey
// SupportMessageBagTest::testFirstReturnsEmptyStringIfNoMessagesFound
// SupportMessageBagTest::testFirstReturnsSingleMessage
// SupportMessageBagTest::testFirstReturnsSingleMessageFromDotKeys
// SupportMessageBagTest::testFormatIsRespected
// SupportMessageBagTest::testGetReturnsArrayOfMessagesByImplicitKey
// SupportMessageBagTest::testHasAnyIndicatesExistence
// SupportMessageBagTest::testHasAnyWithKeyNull
// SupportMessageBagTest::testHasIndicatesExistence
// SupportMessageBagTest::testHasIndicatesExistenceOfAllKeys
// SupportMessageBagTest::testHasIndicatesNoneExistence
// SupportMessageBagTest::testHasWithKeyNull
// SupportMessageBagTest::testIsEmptyFalse
// SupportMessageBagTest::testIsEmptyTrue
// SupportMessageBagTest::testIsNotEmptyFalse
// SupportMessageBagTest::testIsNotEmptyTrue
// SupportMessageBagTest::testMessageBagReturnsCorrectArray
// SupportMessageBagTest::testMessageBagReturnsExpectedJson
// SupportMessageBagTest::testMessageBagsCanBeMerged
// SupportMessageBagTest::testMessageBagsCanConvertToArrays
// SupportMessageBagTest::testMessagesMayBeMerged
// SupportMessageBagTest::testMissingIndicatesNonExistence
// SupportMessageBagTest::testToString
// SupportMessageBagTest::testUnique

// Ref: @bedrock/code-0378
func TestMessageBagUniqueness(t *testing.T) {
	t.Parallel()

	bag := NewMessageBag()
	bag.Add("foo", "bar")
	bag.Add("foo", "bar")

	if bag.Count() != 1 {
		t.Errorf("expected 1 unique message, got %d", bag.Count())
	}
}

// Ref: @bedrock/code-0378
func TestMessageBagAdd(t *testing.T) {
	t.Parallel()

	bag := NewMessageBag()
	bag.Add("email", "must be valid")
	bag.Add("email", "is required")
	bag.Add("name", "is required")

	if bag.Count() != 3 {
		t.Errorf("expected 3, got %d", bag.Count())
	}

	if !bag.Has("email") {
		t.Error("expected Has('email') = true")
	}

	if !bag.Has("name") {
		t.Error("expected Has('name') = true")
	}
}

// Ref: @bedrock/code-0378
func TestMessageBagAddIf(t *testing.T) {
	t.Parallel()

	bag := NewMessageBag()
	bag.AddIf(false, "key", "should not appear")

	if bag.Count() != 0 {
		t.Errorf("expected 0, got %d", bag.Count())
	}

	bag.AddIf(true, "key", "should appear")

	if bag.Count() != 1 {
		t.Errorf("expected 1, got %d", bag.Count())
	}
}

// Ref: @bedrock/code-0378
func TestMessageBagGet(t *testing.T) {
	t.Parallel()

	bag := NewMessageBag()
	bag.Add("email", "Foo")
	bag.Add("email", "Bar")

	msgs := bag.Get("email")

	if len(msgs) != 2 {
		t.Errorf("expected 2, got %d", len(msgs))
	}

	if msgs[0] != "Foo" || msgs[1] != "Bar" {
		t.Errorf("unexpected messages: %v", msgs)
	}
}

// Ref: @bedrock/code-0378
func TestMessageBagGetWildcard(t *testing.T) {
	t.Parallel()

	bag := NewMessageBag()
	bag.Add("messages.0", "First")
	bag.Add("messages.1", "Second")

	msgs := bag.Get("messages.*")

	if len(msgs) != 2 {
		t.Errorf("expected 2, got %d: %v", len(msgs), msgs)
	}
}

// Ref: @bedrock/code-0378
func TestMessageBagFirst(t *testing.T) {
	t.Parallel()

	bag := NewMessageBag()
	bag.Add("email", "Alpha")
	bag.Add("email", "Beta")

	if got := bag.First("email"); got != "Alpha" {
		t.Errorf("First('email') = %q", got)
	}
}

// Ref: @bedrock/code-0378
func TestMessageBagFirstNoKey(t *testing.T) {
	t.Parallel()

	bag := NewMessageBag()
	bag.Add("foo", "Hello")

	if got := bag.First(); got != "Hello" {
		t.Errorf("First() = %q", got)
	}
}

// Ref: @bedrock/code-0378
func TestMessageBagFirstWildcard(t *testing.T) {
	t.Parallel()

	bag := NewMessageBag()
	bag.Add("messages.0", "First")
	bag.Add("messages.1", "Second")

	if got := bag.First("messages.*"); got != "First" {
		t.Errorf("First(messages.*) = %q", got)
	}
}

// Ref: @bedrock/code-0378
func TestMessageBagFirstDotKeys(t *testing.T) {
	t.Parallel()

	bag := NewMessageBag()
	bag.Add("users.0.name", "Taylor")
	bag.Add("users.1.name", "Abigail")

	if got := bag.First("users.0.name"); got != "Taylor" {
		t.Errorf("First(users.0.name) = %q", got)
	}
}

// Ref: @bedrock/code-0378
func TestMessageBagFirstEmpty(t *testing.T) {
	t.Parallel()

	bag := NewMessageBag()

	if got := bag.First(); got != "" {
		t.Errorf("First() on empty bag = %q", got)
	}

	if got := bag.First("missing"); got != "" {
		t.Errorf("First('missing') = %q", got)
	}
}

// Ref: @bedrock/code-0378
func TestMessageBagHas(t *testing.T) {
	t.Parallel()

	bag := NewMessageBag()
	bag.Add("email", "error")

	if !bag.Has("email") {
		t.Error("expected Has('email') = true")
	}

	if bag.Has("name") {
		t.Error("expected Has('name') = false")
	}
}

// Ref: @bedrock/code-0378
func TestMessageBagHasNilKey(t *testing.T) {
	t.Parallel()

	bag := NewMessageBag()
	bag.Add("email", "error")

	if !bag.Has() {
		t.Error("Has() should return true when the bag is not empty")
	}

	if NewMessageBag().Has() {
		t.Error("Has() should return false for an empty bag")
	}
}

// Ref: @bedrock/code-0378
func TestMessageBagHasEmpty(t *testing.T) {
	t.Parallel()

	bag := NewMessageBag()

	if bag.Has("email") {
		t.Error("Has on empty bag should return false")
	}
}

// Ref: @bedrock/code-0378
func TestMessageBagHasMultipleKeys(t *testing.T) {
	t.Parallel()

	bag := NewMessageBag()
	bag.Add("email", "error")
	bag.Add("name", "error")

	if !bag.Has("email", "name") {
		t.Error("Has('email','name') should be true")
	}

	if bag.Has("email", "missing") {
		t.Error("Has('email','missing') should be false")
	}
}

// Ref: @bedrock/code-0378
func TestMessageBagHasAnyNilKey(t *testing.T) {
	t.Parallel()

	bag := NewMessageBag()
	bag.Add("email", "error")

	if !bag.HasAny() {
		t.Error("HasAny() should return true when the bag is not empty")
	}

	if NewMessageBag().HasAny() {
		t.Error("HasAny() should return false for an empty bag")
	}
}

// Ref: @bedrock/code-0378
func TestMessageBagHasAny(t *testing.T) {
	t.Parallel()

	bag := NewMessageBag()
	bag.Add("email", "error")

	if !bag.HasAny("email", "name") {
		t.Error("HasAny should return true when at least one key exists")
	}

	if bag.HasAny("name", "phone") {
		t.Error("HasAny should return false when no key exists")
	}
}

// Ref: @bedrock/code-0378
func TestMessageBagHasAnyWildcard(t *testing.T) {
	t.Parallel()

	bag := NewMessageBag()
	bag.Add("messages.0", "error")

	if !bag.HasAny("messages.*") {
		t.Error("HasAny with wildcard should return true")
	}
}

// Ref: @bedrock/code-0378
func TestMessageBagMissing(t *testing.T) {
	t.Parallel()

	bag := NewMessageBag()
	bag.Add("email", "error")

	if bag.Missing("email") {
		t.Error("Missing('email') should be false when key exists")
	}

	if !bag.Missing("name") {
		t.Error("Missing('name') should be true when key absent")
	}
}

// Ref: @bedrock/code-0378
func TestMessageBagAll(t *testing.T) {
	t.Parallel()

	bag := NewMessageBag()
	bag.Add("email", "e1")
	bag.Add("name", "n1")

	all := bag.All()

	if len(all) != 2 {
		t.Errorf("All() should return 2, got %d", len(all))
	}
}

// Ref: @bedrock/code-0378
func TestMessageBagMerge(t *testing.T) {
	t.Parallel()

	bag1 := NewMessageBag()
	bag1.Add("email", "e1")

	bag2 := NewMessageBag()
	bag2.Add("name", "n1")

	bag1.Merge(bag2)

	if !bag1.Has("email") || !bag1.Has("name") {
		t.Error("Merge should add all keys from the second bag")
	}
}

// Ref: @bedrock/code-0378
func TestMessageBagMergeMap(t *testing.T) {
	t.Parallel()

	bag := NewMessageBag()
	bag.Merge(map[string][]string{
		"email": {"error1", "error2"},
		"name":  {"required"},
	})

	if bag.Count() != 3 {
		t.Errorf("expected 3 messages, got %d", bag.Count())
	}
}

// Ref: @bedrock/code-0378
func TestMessageBagForget(t *testing.T) {
	t.Parallel()

	bag := NewMessageBag()
	bag.Add("email", "error")
	bag.Add("name", "error")
	bag.Forget("email")

	if bag.Has("email") {
		t.Error("Forget should remove 'email'")
	}

	if !bag.Has("name") {
		t.Error("Forget should not remove 'name'")
	}
}

// Ref: @bedrock/code-0378
func TestMessageBagKeys(t *testing.T) {
	t.Parallel()

	bag := NewMessageBag()
	bag.Add("email", "error")
	bag.Add("name", "error")

	keys := bag.Keys()

	if len(keys) != 2 {
		t.Errorf("expected 2 keys, got %d: %v", len(keys), keys)
	}
}

// Ref: @bedrock/code-0378
func TestMessageBagCount(t *testing.T) {
	t.Parallel()

	bag := NewMessageBag()

	if bag.Count() != 0 {
		t.Error("empty bag count should be 0")
	}

	bag.Add("a", "msg1")
	bag.Add("a", "msg2")
	bag.Add("b", "msg3")

	if bag.Count() != 3 {
		t.Errorf("expected 3, got %d", bag.Count())
	}
}

// Ref: @bedrock/code-0378
func TestMessageBagIsEmpty(t *testing.T) {
	t.Parallel()

	bag := NewMessageBag()

	if !bag.IsEmpty() {
		t.Error("new bag should be empty")
	}

	if bag.IsNotEmpty() {
		t.Error("new bag should not be not-empty")
	}

	bag.Add("key", "val")

	if bag.IsEmpty() {
		t.Error("bag with message should not be empty")
	}

	if !bag.IsNotEmpty() {
		t.Error("bag with message should be not-empty")
	}
}

// Ref: @bedrock/code-0378
func TestMessageBagFormat(t *testing.T) {
	t.Parallel()

	bag := NewMessageBag()

	if bag.GetFormat() != ":message" {
		t.Errorf("default format = %q", bag.GetFormat())
	}

	bag.SetFormat("<p>:message</p>")
	bag.Add("key", "hello")

	first := bag.First("key")

	if first != "<p>hello</p>" {
		t.Errorf("formatted message = %q", first)
	}
}

// Ref: @bedrock/code-0378
func TestMessageBagCustomFormat(t *testing.T) {
	t.Parallel()

	bag := NewMessageBag()
	bag.SetFormat("<li>:message</li>")
	bag.Add("email", "invalid")
	bag.Add("email", "required")

	msgs := bag.Get("email")

	for _, m := range msgs {
		if m != "<li>invalid</li>" && m != "<li>required</li>" {
			t.Errorf("unexpected formatted message: %q", m)
		}
	}
}

// Ref: @bedrock/code-0378
func TestMessageBagGetMessages(t *testing.T) {
	t.Parallel()

	bag := NewMessageBag()
	bag.Add("a", "one")
	bag.Add("a", "two")

	raw := bag.GetMessages()

	if len(raw["a"]) != 2 {
		t.Errorf("GetMessages['a'] = %v", raw["a"])
	}
}

// Ref: @bedrock/code-0378
func TestMessageBagJSON(t *testing.T) {
	t.Parallel()

	bag := NewMessageBag()
	bag.Add("email", "invalid")

	data, err := bag.ToJSON()

	if err != nil {
		t.Fatalf("ToJSON error: %v", err)
	}

	var result map[string][]string

	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if len(result["email"]) != 1 || result["email"][0] != "invalid" {
		t.Errorf("JSON result = %v", result)
	}
}

// Ref: @bedrock/code-0378
func TestMessageBagConstructor(t *testing.T) {
	t.Parallel()

	bag := NewMessageBag(map[string][]string{
		"email": {"bad", "bad", "missing"},
		"name":  {"required"},
	})

	if bag.Count() != 3 {
		t.Errorf("expected 3, got %d", bag.Count())
	}
}

// Ref: @bedrock/code-0378
func TestMessageBagConstructorUniquenessConsistency(t *testing.T) {
	t.Parallel()

	bag := NewMessageBag(map[string][]string{
		"email": {"bad", "bad", "missing"},
		"name":  {"required"},
	})

	if got := bag.Get("email"); len(got) != 2 {
		t.Fatalf("constructor dedup = %v", got)
	}
}

// Ref: @bedrock/code-0378
func TestMessageBagUnique(t *testing.T) {
	t.Parallel()

	bag := NewMessageBag()
	bag.messages["key"] = []string{"dup", "dup", "unique"}

	unique := bag.Unique()
	msgs := unique.Get("key")

	if len(msgs) != 2 {
		t.Errorf("Unique should deduplicate, got %v", msgs)
	}
}

// Ref: @bedrock/code-0378
func TestMessageBagAllWithFormat(t *testing.T) {
	t.Parallel()

	bag := NewMessageBag()
	bag.SetFormat("[error: :message]")
	bag.Add("x", "bad")

	all := bag.All()

	if len(all) != 1 || all[0] != "[error: bad]" {
		t.Errorf("All() with format = %v", all)
	}
}

// Ref: @bedrock/code-0378
func TestMessageBagHasWildcard(t *testing.T) {
	t.Parallel()

	bag := NewMessageBag()
	bag.Add("errors.email", "bad")
	bag.Add("errors.name", "required")

	if !bag.Has("errors.*") {
		t.Error("Has with wildcard should match")
	}

	if bag.Has("other.*") {
		t.Error("Has with non-matching wildcard should return false")
	}
}

// Ref: @bedrock/code-0378
func TestMessageBagForgetMultiple(t *testing.T) {
	t.Parallel()

	bag := NewMessageBag()
	bag.Add("a", "msg")
	bag.Add("b", "msg")
	bag.Add("c", "msg")
	bag.Forget("a", "b")

	if bag.Has("a") || bag.Has("b") {
		t.Error("Forget should remove both a and b")
	}

	if !bag.Has("c") {
		t.Error("Forget should leave c")
	}
}

// Ref: @bedrock/code-0378
func TestMessageBagString(t *testing.T) {
	t.Parallel()

	bag := NewMessageBag()
	bag.Add("key", "value")

	s := bag.String()

	if s == "" || s == "{}" {
		t.Errorf("String() = %q — expected non-empty JSON", s)
	}
}
