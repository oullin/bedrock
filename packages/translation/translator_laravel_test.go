package translation_test

// Ports of Illuminate\Tests\Translation\TranslationTranslatorTest.

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/bedrock/packages/translation"
)

// ── fixtures ──────────────────────────────────────────────────────────────

// seed loads messages into the translator's ArrayLoader via AddLines so that
// all subsequent Get calls hit the in-memory cache.

// ── has ───────────────────────────────────────────────────────────────────

// has returns false when Get returns the key itself.

// ── get ───────────────────────────────────────────────────────────────────

// Reset loaded cache so we can seed fresh messages.

// When locale == fallback the translator must not attempt two lookups
// for the same locale/group triple.

// Key is missing in every locale.

// may load "*" group and "messages" group once each

// countingLoader wraps ArrayLoader and counts Load invocations.
type countingLoader struct {
	*translation.ArrayLoader
	calls int
}

type countedItems struct{ n int }

// ── choice ────────────────────────────────────────────────────────────────

// cs has no translation; fallback to en is used for both fetch and plural form.
// English binary rule: n=10 → index 1 → "few".

// ── JSON flat lookups ─────────────────────────────────────────────────────

// :bar must NOT be substituted inside the value of :foo.

// ── Stringable ────────────────────────────────────────────────────────────

// fakeDate is a Go stand-in for Carbon (a Stringable with custom format).
type fakeDate struct{ unix int64 }

// Default fmt.Stringer returns Unix timestamp string.

// Without custom handler, String() is used.

// With a custom Stringable handler, the handler overrides String().

// ── Enum-style replacements ───────────────────────────────────────────────

// Go equivalent of PHP string-backed enum.
type month string

// Go equivalent of PHP int-backed enum.
type version int

// Go equivalent of PHP unit enum (no associated value).
type person struct{ name string }

func (c countedItems) Len() int { return c.n }

func newTranslator(locale string) *translation.Translator {
	return translation.NewTranslator(translation.NewArrayLoader(), locale)
}

func seed(t *translation.Translator, locale, ns, group string, msgs map[string]any) {
	t.GetLoader().(*translation.ArrayLoader).AddMessages(locale, group, msgs, ptrStr(ns))
}

func ptrStr(s string) *string { return &s }

// Port of Illuminate\Tests\Translation\TranslationTranslatorTest::testHasMethodReturnsFalseWhenReturnedTranslationIsNull
func TestHasMethodReturnsFalseWhenReturnedTranslationIsNull(t *testing.T) {
	t.Parallel()

	tr := newTranslator("en")
	tr.GetLoader().(*translation.ArrayLoader).AddMessages("en", "foo", map[string]any{"bar": nil}, nil)

	if tr.Has("foo.bar", ptrStr("en")) {
		t.Error("Has should return false when the loaded translation value is nil")
	}
}

func TestHasMethodReturnsFalseWhenTranslationEqualsKey(t *testing.T) {
	t.Parallel()

	tr := newTranslator("en")
	seed(tr, "en", "*", "foo", map[string]any{"bar": "foo.bar"})

	if tr.Has("foo.baz", ptrStr("en")) {
		t.Error("Has should return false when key is missing")
	}
}

// Port of Illuminate\Tests\Translation\TranslationTranslatorTest::testHasMethodReturnsTrueWhenReturnedTranslationIsNotNull
func TestHasMethodReturnsTrueWhenTranslationDiffersFromKey(t *testing.T) {
	t.Parallel()

	tr := newTranslator("en")
	seed(tr, "en", "*", "foo", map[string]any{"bar": "translated value"})

	if !tr.Has("foo.bar", ptrStr("en")) {
		t.Error("Has should return true when a translation exists")
	}
}

// Port of Illuminate\Tests\Translation\TranslationTranslatorTest::testHasMethodReturnsTrueWhenReturnedTranslationIsNotNullForLocale
func TestHasForLocaleReturnsTrueWhenKeyExists(t *testing.T) {
	t.Parallel()

	tr := newTranslator("en")
	seed(tr, "en", "*", "foo", map[string]any{"bar": "value"})

	if !tr.HasForLocale("foo.bar", ptrStr("en")) {
		t.Error("HasForLocale should return true when key exists in locale")
	}
}

// Port of Illuminate\Tests\Translation\TranslationTranslatorTest::testHasMethodReturnsFalseWhenReturnedTranslationIsNullForLocale
func TestHasForLocaleReturnsFalseWhenKeyMissing(t *testing.T) {
	t.Parallel()

	tr := newTranslator("en")
	seed(tr, "en", "*", "foo", map[string]any{})

	if tr.HasForLocale("foo.bar", ptrStr("en")) {
		t.Error("HasForLocale should return false when key is absent")
	}
}

// Port of Illuminate\Tests\Translation\TranslationTranslatorTest::testGetMethodProperlyLoadsAndRetrievesItem
func TestGetMethodProperlyLoadsAndRetrievesItem(t *testing.T) {
	t.Parallel()

	tr := newTranslator("en")
	seed(tr, "en", "foo", "bar", map[string]any{
		"foo": "foo",
		"baz": "breeze :foo",
		"qux": []any{"tree :foo", "breeze :foo"},
	})

	got := tr.Get("foo::bar.qux", map[string]any{"foo": "bar"}, ptrStr("en"))
	want := []any{"tree bar", "breeze bar"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v; want %v", got, want)
	}

	if tr.Get("foo::bar.baz", map[string]any{"foo": "bar"}, ptrStr("en")) != "breeze bar" {
		t.Errorf("expected 'breeze bar'")
	}

	if tr.Get("foo::bar.foo", nil, nil) != "foo" {
		t.Errorf("expected 'foo'")
	}
}

// Port of Illuminate\Tests\Translation\TranslationTranslatorTest::testGetMethodProperlyLoadsAndRetrievesArrayItem
func TestGetMethodProperlyLoadsAndRetrievesArrayItem(t *testing.T) {
	t.Parallel()

	tr := newTranslator("en")
	seed(tr, "en", "foo", "bar", map[string]any{
		"foo": "foo",
		"baz": "breeze :foo",
		"qux": []any{"tree :foo", "breeze :foo", map[string]any{"beep": map[string]any{"rock": "tree :foo"}}},
	})

	got := tr.Get("foo::bar", map[string]any{"foo": "bar"}, ptrStr("en"))
	want := map[string]any{
		"foo": "foo",
		"baz": "breeze bar",
		"qux": []any{"tree bar", "breeze bar", map[string]any{"beep": map[string]any{"rock": "tree bar"}}},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v; want %v", got, want)
	}
}

// Port of Illuminate\Tests\Translation\TranslationTranslatorTest::testGetMethodForNonExistingReturnsSameKey
func TestGetMethodForNonExistingReturnsSameKey(t *testing.T) {
	t.Parallel()

	tr := newTranslator("en")
	seed(tr, "en", "foo", "bar", map[string]any{"foo": "foo"})

	if tr.Get("foo::unknown", map[string]any{"foo": "bar"}, ptrStr("en")) != "foo::unknown" {
		t.Error("expected key itself for missing namespace")
	}

	if tr.Get("foo::bar.unknown", map[string]any{"foo": "bar"}, ptrStr("en")) != "foo::bar.unknown" {
		t.Error("expected key itself for missing item in known group")
	}

	if tr.Get("foo::unknown.bar", nil, nil) != "foo::unknown.bar" {
		t.Error("expected key itself for unknown group")
	}
}

// Port of Illuminate\Tests\Translation\TranslationTranslatorTest::testTransMethodProperlyLoadsAndRetrievesItemWithHTMLInTheMessage
func TestGetMethodWithHTMLInMessage(t *testing.T) {
	t.Parallel()

	tr := newTranslator("en")
	seed(tr, "en", "*", "foo", map[string]any{"bar": "breeze <p>test</p>"})

	got := tr.Get("foo.bar", nil, ptrStr("en"))

	if got != "breeze <p>test</p>" {
		t.Errorf("got %v", got)
	}
}

// Port of Illuminate\Tests\Translation\TranslationTranslatorTest::testGetMethodProperlyLoadsAndRetrievesItemWithCapitalization
func TestGetMethodProperlyLoadsAndRetrievesItemWithCapitalization(t *testing.T) {
	t.Parallel()

	tr := newTranslator("en")
	seed(tr, "en", "foo", "bar", map[string]any{
		"baz": "breeze :0 :Foo :BAR",
	})

	got := tr.Get("foo::bar.baz", map[string]any{"0": "john", "foo": "bar", "bar": "foo"}, ptrStr("en"))

	if got != "breeze john Bar FOO" {
		t.Errorf("got %q; want %q", got, "breeze john Bar FOO")
	}
}

// Port of Illuminate\Tests\Translation\TranslationTranslatorTest::testGetMethodProperlyLoadsAndRetrievesItemWithLongestReplacementsFirst
func TestGetMethodProperlyLoadsAndRetrievesItemWithLongestReplacementsFirst(t *testing.T) {
	t.Parallel()

	tr := newTranslator("en")
	seed(tr, "en", "foo", "bar", map[string]any{
		"baz": "breeze :foo :foobar",
	})

	got := tr.Get("foo::bar.baz", map[string]any{"foo": "bar", "foobar": "taylor"}, ptrStr("en"))

	if got != "breeze bar taylor" {
		t.Errorf("got %q; want %q", got, "breeze bar taylor")
	}

	tr2 := newTranslator("en")
	seed(tr2, "en", "foo", "bar", map[string]any{
		"baz": "breeze :foo :foobar",
	})
	got2 := tr2.Get("foo::bar.baz", map[string]any{"foo": "foo bar baz", "foobar": "taylor"}, ptrStr("en"))

	if got2 != "breeze foo bar baz taylor" {
		t.Errorf("got %q; want %q", got2, "breeze foo bar baz taylor")
	}
}

// Port of Illuminate\Tests\Translation\TranslationTranslatorTest::testGetMethodProperlyLoadsAndRetrievesItemForFallback
func TestGetMethodProperlyLoadsAndRetrievesItemForFallback(t *testing.T) {
	t.Parallel()

	tr := newTranslator("en")
	tr.SetFallback("lv")
	seed(tr, "en", "foo", "bar", map[string]any{})
	seed(tr, "lv", "foo", "bar", map[string]any{"baz": "breeze :foo"})

	got := tr.Get("foo::bar.baz", map[string]any{"foo": "bar"}, ptrStr("en"))

	if got != "breeze bar" {
		t.Errorf("got %q; want %q", got, "breeze bar")
	}
}

// Port of Illuminate\Tests\Translation\TranslationTranslatorTest::testGetDoesNotCallGetLineTwiceForMissingKeyWhenLocaleMatchesFallback
func TestGetDoesNotCallLoaderTwiceWhenLocaleMatchesFallback(t *testing.T) {
	t.Parallel()

	loader := &countingLoader{ArrayLoader: translation.NewArrayLoader()}
	tr := translation.NewTranslator(loader, "en")
	tr.SetFallback("en")

	result := tr.Get("messages.test", nil, ptrStr("en"))

	if result != "messages.test" {
		t.Errorf("expected key itself, got %v", result)
	}

	if loader.calls > 2 {
		t.Errorf("loader called %d times; expected ≤2", loader.calls)
	}
}

func (c *countingLoader) Load(locale, group string, namespace *string) map[string]any {
	c.calls++

	return c.ArrayLoader.Load(locale, group, namespace)
}

// Port of Illuminate\Tests\Translation\TranslationTranslatorTest::testGetMethodProperlyLoadsAndRetrievesItemForGlobalNamespace
func TestGetMethodProperlyLoadsAndRetrievesItemForGlobalNamespace(t *testing.T) {
	t.Parallel()

	tr := newTranslator("en")
	seed(tr, "en", "*", "foo", map[string]any{"bar": "breeze :foo"})

	got := tr.Get("foo.bar", map[string]any{"foo": "bar"}, nil)

	if got != "breeze bar" {
		t.Errorf("got %q; want %q", got, "breeze bar")
	}
}

// Port of Illuminate\Tests\Translation\TranslationTranslatorTest::testChoiceMethodProperlyLoadsAndRetrievesItemForAnInt
func TestChoiceMethodProperlyLoadsAndRetrievesItemForAnInt(t *testing.T) {
	t.Parallel()

	tr := newTranslator("en")
	tr.AddLines(map[string]any{"foo": "one item|many items"}, "en")

	got := tr.Choice("foo", 10, nil, ptrStr("en"))

	if got != "many items" {
		t.Errorf("got %q", got)
	}
}

// Port of Illuminate\Tests\Translation\TranslationTranslatorTest::testChoiceMethodProperlyLoadsAndRetrievesItemForAFloat
func TestChoiceMethodProperlyLoadsAndRetrievesItemForAFloat(t *testing.T) {
	t.Parallel()

	tr := newTranslator("en")
	tr.AddLines(map[string]any{"foo": "one item|many items"}, "en")

	got := tr.Choice("foo", 1.2, nil, ptrStr("en"))

	if got != "many items" {
		t.Errorf("got %q", got)
	}
}

// Port of Illuminate\Tests\Translation\TranslationTranslatorTest::testChoiceMethodProperlyCountsCollectionsAndLoadsAndRetrievesItem
func TestChoiceMethodProperlyCountsCollectionsAndLoadsAndRetrievesItem(t *testing.T) {
	t.Parallel()

	tr := newTranslator("en")
	tr.AddLines(map[string]any{"foo": "{1} one|[2,*] many"}, "en")

	got := tr.Choice("foo", countedItems{n: 3}, nil, ptrStr("en"))

	if got != "many" {
		t.Errorf("got %q", got)
	}
}

func TestChoiceMethodProperlyCountsSlices(t *testing.T) {
	t.Parallel()

	tr := newTranslator("en")
	tr.AddLines(map[string]any{"foo": "{1} one|[2,*] many"}, "en")

	values := []string{"a", "b", "c"}
	got := tr.Choice("foo", values, nil, ptrStr("en"))

	if got != "many" {
		t.Errorf("got %q", got)
	}
}

// Port of Illuminate\Tests\Translation\TranslationTranslatorTest::testChoiceMethodProperlySelectsLocaleForChoose
func TestChoiceMethodProperlySelectsLocaleForChoose(t *testing.T) {
	t.Parallel()

	tr := newTranslator("cs")
	tr.SetFallback("en")
	seed(tr, "en", "*", "*", map[string]any{"foo": "one|few|many"})

	got := tr.Choice("foo", 10, nil, nil)

	if got != "few" {
		t.Errorf("got %q; want %q", got, "few")
	}
}

// Port of Illuminate\Tests\Translation\TranslationTranslatorTest::testChoiceMethodProperlyUsesCustomCountReplacement
func TestChoiceMethodProperlyUsesCustomCountReplacement(t *testing.T) {
	t.Parallel()

	tr := newTranslator("en")
	tr.AddLines(map[string]any{":count foos": "{1} :count foos|[2,*] :count foos"}, "en")

	got := tr.Choice(":count foos", 1234, map[string]any{"count": "1,234"}, ptrStr("en"))

	if got != "1,234 foos" {
		t.Errorf("got %q; want %q", got, "1,234 foos")
	}
}

// Port of Illuminate\Tests\Translation\TranslationTranslatorTest::testGetJson
func TestGetJson(t *testing.T) {
	t.Parallel()

	tr := newTranslator("en")
	tr.AddLines(map[string]any{"foo": "one"}, "en")

	got := tr.Get("foo", nil, nil)

	if got != "one" {
		t.Errorf("got %v", got)
	}
}

// Port of Illuminate\Tests\Translation\TranslationTranslatorTest::testGetJsonReplaces
func TestGetJsonReplaces(t *testing.T) {
	t.Parallel()

	tr := newTranslator("en")
	tr.AddLines(map[string]any{"foo :i:c :u": "bar :i:c :u"}, "en")

	got := tr.Get("foo :i:c :u", map[string]any{"i": "one", "c": "two", "u": "three"}, nil)

	if got != "bar onetwo three" {
		t.Errorf("got %q", got)
	}
}

// Port of Illuminate\Tests\Translation\TranslationTranslatorTest::testGetJsonHasAtomicReplacements
func TestGetJsonHasAtomicReplacements(t *testing.T) {
	t.Parallel()

	tr := newTranslator("en")
	tr.AddLines(map[string]any{"Hello :foo!": "Hello :foo!"}, "en")

	got := tr.Get("Hello :foo!", map[string]any{"foo": "baz:bar", "bar": "abcdef"}, nil)

	if got != "Hello baz:bar!" {
		t.Errorf("got %q; want %q", got, "Hello baz:bar!")
	}
}

// Port of Illuminate\Tests\Translation\TranslationTranslatorTest::testGetJsonReplacesForAssociativeInput
func TestGetJsonReplacesForAssociativeInput(t *testing.T) {
	t.Parallel()

	tr := newTranslator("en")
	tr.AddLines(map[string]any{"foo :i :c": "bar :i :c"}, "en")

	got := tr.Get("foo :i :c", map[string]any{"i": "eye", "c": "see"}, nil)

	if got != "bar eye see" {
		t.Errorf("got %q", got)
	}
}

// Port of Illuminate\Tests\Translation\TranslationTranslatorTest::testGetJsonPreservesOrder
func TestGetJsonPreservesOrder(t *testing.T) {
	t.Parallel()

	tr := newTranslator("en")
	tr.AddLines(map[string]any{"to :name I give :greeting": ":greeting :name"}, "en")

	got := tr.Get("to :name I give :greeting", map[string]any{"name": "David", "greeting": "Greetings"}, nil)

	if got != "Greetings David" {
		t.Errorf("got %q", got)
	}
}

// Port of Illuminate\Tests\Translation\TranslationTranslatorTest::testGetJsonForNonExistingJsonKeyLooksForRegularKeys
func TestGetJsonForNonExistingJsonKeyLooksForRegularKeys(t *testing.T) {
	t.Parallel()

	tr := newTranslator("en")
	seed(tr, "en", "*", "foo", map[string]any{"bar": "one"})

	got := tr.Get("foo.bar", nil, nil)

	if got != "one" {
		t.Errorf("got %v", got)
	}
}

// Port of Illuminate\Tests\Translation\TranslationTranslatorTest::testGetJsonForNonExistingJsonKeyLooksForRegularKeysAndReplace
func TestGetJsonForNonExistingJsonKeyLooksForRegularKeysAndReplace(t *testing.T) {
	t.Parallel()

	tr := newTranslator("en")
	seed(tr, "en", "*", "foo", map[string]any{"bar": "one :message"})

	got := tr.Get("foo.bar", map[string]any{"message": "two"}, nil)

	if got != "one two" {
		t.Errorf("got %q", got)
	}
}

// Port of Illuminate\Tests\Translation\TranslationTranslatorTest::testGetJsonForNonExistingReturnsSameKey
func TestGetJsonForNonExistingReturnsSameKey(t *testing.T) {
	t.Parallel()

	tr := newTranslator("en")

	got := tr.Get("Foo that bar", nil, nil)

	if got != "Foo that bar" {
		t.Errorf("got %v", got)
	}
}

// Port of Illuminate\Tests\Translation\TranslationTranslatorTest::testGetJsonForNonExistingReturnsSameKeyAndReplaces
func TestGetJsonForNonExistingReturnsSameKeyAndReplaces(t *testing.T) {
	t.Parallel()

	tr := newTranslator("en")

	got := tr.Get("foo :message", map[string]any{"message": "baz"}, nil)

	if got != "foo baz" {
		t.Errorf("got %q", got)
	}
}

// Port of Illuminate\Tests\Translation\TranslationTranslatorTest::testEmptyFallbacks
func TestEmptyFallbacks(t *testing.T) {
	t.Parallel()

	tr := newTranslator("en")

	got := tr.Get("foo :message", map[string]any{"message": nil}, nil)

	if got != "foo " {
		t.Errorf("got %q; want %q", got, "foo ")
	}
}

func (d fakeDate) String() string {

	return fmt.Sprintf("%d", d.unix)
}

// Port of Illuminate\Tests\Translation\TranslationTranslatorTest::testGetJsonReplacesWithStringable
func TestGetJsonReplacesWithStringable(t *testing.T) {
	t.Parallel()

	tr := newTranslator("en")
	tr.AddLines(map[string]any{"test": "the date is :date"}, "en")

	date := fakeDate{unix: 0}

	got1 := tr.Get("test", map[string]any{"date": date}, nil)

	if got1 != "the date is 0" {
		t.Errorf("got %q", got1)
	}

	tr.Stringable(reflect.TypeOf(fakeDate{}), func(v any) string {
		d := v.(fakeDate)

		return fmt.Sprintf("1st Jan %d", 1970+int(d.unix))
	})

	got2 := tr.Get("test", map[string]any{"date": date}, nil)

	if got2 != "the date is 1st Jan 1970" {
		t.Errorf("got %q", got2)
	}
}

func (m month) String() string { return string(m) }

const february month = "February"

func (v version) String() string { return fmt.Sprintf("%d", int(v)) }

const thirteen version = 13

func (p person) String() string { return p.name }

var hosni = person{name: "Hosni"}

// Port of Illuminate\Tests\Translation\TranslationTranslatorTest::testGetJsonReplacesWithEnums
func TestGetJsonReplacesWithEnums(t *testing.T) {
	t.Parallel()

	tr := newTranslator("en")
	tr.AddLines(map[string]any{
		"string_backed_enum": "Laravel 12 was released in :month 2025",
		"int_backed_enum":    "Stay tuned for Laravel v:version",
		"unit_enum":          ":person gets excited about every new Laravel release",
	}, "en")

	got1 := tr.Get("string_backed_enum", map[string]any{"month": february}, nil)

	if got1 != "Laravel 12 was released in February 2025" {
		t.Errorf("got %q", got1)
	}

	got2 := tr.Get("int_backed_enum", map[string]any{"version": thirteen}, nil)

	if got2 != "Stay tuned for Laravel v13" {
		t.Errorf("got %q", got2)
	}

	got3 := tr.Get("unit_enum", map[string]any{"person": hosni}, nil)

	if got3 != "Hosni gets excited about every new Laravel release" {
		t.Errorf("got %q", got3)
	}
}

// ── Tag replacements ──────────────────────────────────────────────────────

// Port of Illuminate\Tests\Translation\TranslationTranslatorTest::testTagReplacements
func TestTagReplacements(t *testing.T) {
	t.Parallel()

	tr := newTranslator("en")

	got := tr.Get(
		"We have some nice <docs-link>documentation</docs-link>",
		map[string]any{
			"docs-link": func(children string) string {
				return `<a href="https://laravel.com/docs">` + children + `</a>`
			},
		},
		nil,
	)

	want := `We have some nice <a href="https://laravel.com/docs">documentation</a>`

	if got != want {
		t.Errorf("got %q; want %q", got, want)
	}
}

// Port of Illuminate\Tests\Translation\TranslationTranslatorTest::testTagReplacementsHandleMultipleOfSameTag
func TestTagReplacementsHandleMultipleOfSameTag(t *testing.T) {
	t.Parallel()

	tr := newTranslator("en")

	got := tr.Get(
		"<bold-this>bold</bold-this> something else <bold-this>also bold</bold-this>",
		map[string]any{
			"bold-this": func(children string) string {
				return "<b>" + children + "</b>"
			},
		},
		nil,
	)

	want := "<b>bold</b> something else <b>also bold</b>"

	if got != want {
		t.Errorf("got %q; want %q", got, want)
	}
}

// ── DetermineLocalesUsing ─────────────────────────────────────────────────

// Port of Illuminate\Tests\Translation\TranslationTranslatorTest::testDetermineLocalesUsingMethod
func TestDetermineLocalesUsingMethod(t *testing.T) {
	t.Parallel()

	tr := newTranslator("en")
	tr.DetermineLocalesUsing(func(locales []string) []string {
		// Append a custom locale to the chain.
		return append(locales, "lz")
	})

	seed(tr, "lz", "*", "foo", map[string]any{"bar": "lz-value"})

	got := tr.Get("foo.bar", nil, nil)

	if got != "lz-value" {
		t.Errorf("got %v; want %q", got, "lz-value")
	}
}
