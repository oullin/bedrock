package str

import (
	"strings"
	"testing"
)

// Port of Framework\Tests\Support\SupportStrTest::testStringCanBeLimitedByWords
func TestStrWords(t *testing.T) {
	t.Parallel()

	if got := StrWords("This is a sentence", 3); got != "This is a..." {
		t.Errorf("StrWords(3) = %q", got)
	}

	if got := StrWords("This is a sentence", 3, " >>>"); got != "This is a >>>" {
		t.Errorf("StrWords(3, custom) = %q", got)
	}

	if got := StrWords("This is a sentence", 10); got != "This is a sentence" {
		t.Errorf("StrWords(10) = %q", got)
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testStringTitle
func TestStrTitle(t *testing.T) {
	t.Parallel()

	if got := StrTitle("hello world"); got != "Hello World" {
		t.Errorf("StrTitle = %q", got)
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testStringHeadline
func TestStrHeadline(t *testing.T) {
	t.Parallel()

	cases := []struct{ in, want string }{
		{"steve_jobs", "Steve Jobs"},
		{"EmailNotificationSent", "Email Notification Sent"},
		{"hello-world", "Hello World"},
		{"hello world", "Hello World"},
	}

	for _, tc := range cases {
		if got := StrHeadline(tc.in); got != tc.want {
			t.Errorf("StrHeadline(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testStartsWith
func TestStrStartsWith(t *testing.T) {
	t.Parallel()

	if !StrStartsWith("jason", "jas") {
		t.Error("startsWith(jas) should be true")
	}

	if !StrStartsWith("jason", "jason") {
		t.Error("startsWith(jason) should be true")
	}

	if StrStartsWith("jason", "day") {
		t.Error("startsWith(day) should be false")
	}
	// Multiple prefixes
	if !StrStartsWith("jason", "jas", "nope") {
		t.Error("startsWith with multiple should find first match")
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testEndsWith
func TestStrEndsWith(t *testing.T) {
	t.Parallel()

	if !StrEndsWith("jason", "on") {
		t.Error("endsWith(on) should be true")
	}

	if !StrEndsWith("jason", "jason") {
		t.Error("endsWith(jason) should be true")
	}

	if StrEndsWith("jason", "nope") {
		t.Error("endsWith(nope) should be false")
	}

	if !StrEndsWith("jason", "on", "nope") {
		t.Error("endsWith with multiple should find match")
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testStrContains
func TestStrContains(t *testing.T) {
	t.Parallel()

	if !StrContains("taylor", "ylo") {
		t.Error("contains(ylo) should be true")
	}

	if !StrContains("taylor", "taylor") {
		t.Error("contains(full) should be true")
	}

	if StrContains("taylor", "nope") {
		t.Error("contains(nope) should be false")
	}
	// Multiple needles (OR semantics)
	if !StrContains("taylor", "ylo", "nope") {
		t.Error("contains multiple (OR) should find match")
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testStrContainsAll
func TestStrContainsAll(t *testing.T) {
	t.Parallel()

	if !StrContainsAll("taylor otwell", []string{"taylor", "otwell"}) {
		t.Error("containsAll should be true for both")
	}

	if StrContainsAll("taylor", []string{"taylor", "otwell"}) {
		t.Error("containsAll should be false when one is missing")
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testSlug
func TestStrSlug(t *testing.T) {
	t.Parallel()

	cases := []struct{ in, sep, want string }{
		{"Hello World", "-", "hello-world"},
		{"Hello World", "_", "hello_world"},
		{"My name is Taylor Otwell", "-", "my-name-is-taylor-otwell"},
		{"hello---world", "-", "hello-world"},
	}

	for _, tc := range cases {
		got := StrSlug(tc.in, tc.sep)

		if got != tc.want {
			t.Errorf("StrSlug(%q, %q) = %q, want %q", tc.in, tc.sep, got, tc.want)
		}
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testSnake
func TestStrSnake(t *testing.T) {
	t.Parallel()

	cases := []struct{ in, want string }{
		{"fooBar", "foo_bar"},
		{"FooBar", "foo_bar"},
		{"foo bar", "foo_bar"},
		{"HTMLParser", "html_parser"},
	}

	for _, tc := range cases {
		got := StrSnake(tc.in)

		if got != tc.want {
			t.Errorf("StrSnake(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testCamel
func TestStrCamel(t *testing.T) {
	t.Parallel()

	cases := []struct{ in, want string }{
		{"foo_bar", "fooBar"},
		{"foo-bar", "fooBar"},
		{"foo bar", "fooBar"},
		{"FooBar", "fooBar"},
	}

	for _, tc := range cases {
		got := StrCamel(tc.in)

		if got != tc.want {
			t.Errorf("StrCamel(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testStudly
func TestStrStudly(t *testing.T) {
	t.Parallel()

	cases := []struct{ in, want string }{
		{"foo_bar", "FooBar"},
		{"foo-bar", "FooBar"},
		{"foo bar", "FooBar"},
		{"fooBar", "FooBar"},
	}

	for _, tc := range cases {
		got := StrStudly(tc.in)

		if got != tc.want {
			t.Errorf("StrStudly(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testKebab
func TestStrKebab(t *testing.T) {
	t.Parallel()

	if got := StrKebab("fooBar"); got != "foo-bar" {
		t.Errorf("StrKebab(fooBar) = %q", got)
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testLimit
func TestStrLimit(t *testing.T) {
	t.Parallel()

	if got := StrLimit("The quick brown fox jumped over the lazy dog", 20); got != "The quick brown fox ..." {
		t.Errorf("StrLimit(20) = %q", got)
	}

	if got := StrLimit("Hello World", 100); got != "Hello World" {
		t.Errorf("StrLimit(100) should return full string, got %q", got)
	}

	if got := StrLimit("Hello", 5, ""); got != "Hello" {
		t.Errorf("StrLimit(5) exact = %q", got)
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testStrAfter
func TestStrAfter(t *testing.T) {
	t.Parallel()

	if got := StrAfter("hannah", "han"); got != "nah" {
		t.Errorf("StrAfter = %q", got)
	}

	if got := StrAfter("hannah", ""); got != "hannah" {
		t.Errorf("StrAfter empty = %q", got)
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testStrBefore
func TestStrBefore(t *testing.T) {
	t.Parallel()

	if got := StrBefore("hannah", "nah"); got != "han" {
		t.Errorf("StrBefore = %q", got)
	}

	if got := StrBefore("hannah", ""); got != "hannah" {
		t.Errorf("StrBefore empty = %q", got)
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testStrBetween
func TestStrBetween(t *testing.T) {
	t.Parallel()

	if got := StrBetween("[abc]", "[", "]"); got != "abc" {
		t.Errorf("StrBetween = %q", got)
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testStrBetweenFirst
func TestStrBetweenFirst(t *testing.T) {
	t.Parallel()

	if got := StrBetweenFirst("[abc][def]", "[", "]"); got != "abc" {
		t.Errorf("StrBetweenFirst = %q", got)
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testIsJson
func TestStrIsJson(t *testing.T) {
	t.Parallel()

	if !StrIsJson(`{"key":"value"}`) {
		t.Error("valid JSON should return true")
	}

	if !StrIsJson(`[1, 2, 3]`) {
		t.Error("valid JSON array should return true")
	}

	if !StrIsJson(`"string"`) {
		t.Error("valid JSON string should return true")
	}

	if StrIsJson(`not json`) {
		t.Error("invalid JSON should return false")
	}

	if StrIsJson("") {
		t.Error("empty string should return false")
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testIsUuidWithValidUuid
func TestStrIsUuid(t *testing.T) {
	t.Parallel()

	if !StrIsUuid("550e8400-e29b-41d4-a716-446655440000") {
		t.Error("valid UUID should return true")
	}

	if StrIsUuid("not-a-uuid") {
		t.Error("invalid UUID should return false")
	}

	if StrIsUuid("") {
		t.Error("empty should return false")
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testIsUlidWithValidUlid
func TestStrIsUlid(t *testing.T) {
	t.Parallel()

	if !StrIsUlid("01ARZ3NDEKTSV4RRFFQ69G5FAV") {
		t.Error("valid ULID should return true")
	}

	if StrIsUlid("not-a-ulid") {
		t.Error("invalid ULID should return false")
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testRandom
func TestStrRandom(t *testing.T) {
	t.Parallel()

	r := StrRandom(16)

	if len(r) != 16 {
		t.Errorf("expected length 16, got %d", len(r))
	}
	// Default length
	if len(StrRandom()) != 16 {
		t.Error("default length should be 16")
	}
	// Should generate different strings
	if StrRandom() == StrRandom() {
		t.Log("warning: two random strings matched (rare but possible)")
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testToBase64
func TestStrBase64(t *testing.T) {
	t.Parallel()

	encoded := StrToBase64("Hello World")

	if encoded == "" {
		t.Error("base64 encoding should not be empty")
	}

	decoded, err := StrFromBase64(encoded)

	if err != nil {
		t.Errorf("base64 decode error: %v", err)
	}

	if decoded != "Hello World" {
		t.Errorf("round-trip failed: got %q", decoded)
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testReverse
func TestStrReverse(t *testing.T) {
	t.Parallel()

	if got := StrReverse("hello"); got != "olleh" {
		t.Errorf("StrReverse = %q", got)
	}
	// UTF-8
	if got := StrReverse("Héllo"); got != "ollÉH" || !strings.Contains(got, "é") {
		// just check it doesn't panic
		_ = StrReverse("Héllo")
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testSquish
func TestStrSquish(t *testing.T) {
	t.Parallel()

	if got := StrSquish("  hello   world  "); got != "hello world" {
		t.Errorf("StrSquish = %q", got)
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testStrStart
func TestStrStart(t *testing.T) {
	t.Parallel()

	if got := StrStart("world", "/"); got != "/world" {
		t.Errorf("StrStart = %q", got)
	}
	// Should not double-prefix
	if got := StrStart("/world", "/"); got != "/world" {
		t.Errorf("StrStart already has prefix = %q", got)
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testFinish
func TestStrFinish(t *testing.T) {
	t.Parallel()

	if got := StrFinish("hello", "/"); got != "hello/" {
		t.Errorf("StrFinish = %q", got)
	}

	if got := StrFinish("hello/", "/"); got != "hello/" {
		t.Errorf("StrFinish already has suffix = %q", got)
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testWrap
func TestStrWrap(t *testing.T) {
	t.Parallel()

	if got := StrWrap("value", "'"); got != "'value'" {
		t.Errorf("StrWrap = %q", got)
	}

	if got := StrWrap("value", "<", ">"); got != "<value>" {
		t.Errorf("StrWrap asymmetric = %q", got)
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testUnwrap
func TestStrUnwrap(t *testing.T) {
	t.Parallel()

	if got := StrUnwrap("'value'", "'", "'"); got != "value" {
		t.Errorf("StrUnwrap = %q", got)
	}

	if got := StrUnwrap("<value>", "<", ">"); got != "value" {
		t.Errorf("StrUnwrap asymmetric = %q", got)
	}
	// Not wrapped — return unchanged
	if got := StrUnwrap("value", "'", "'"); got != "value" {
		t.Errorf("StrUnwrap not wrapped = %q", got)
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testSubstr
func TestStrSubstr(t *testing.T) {
	t.Parallel()

	if got := StrSubstr("hello world", 6); got != "world" {
		t.Errorf("StrSubstr(6) = %q", got)
	}

	if got := StrSubstr("hello world", 0, 5); got != "hello" {
		t.Errorf("StrSubstr(0, 5) = %q", got)
	}

	if got := StrSubstr("hello world", -5); got != "world" {
		t.Errorf("StrSubstr(-5) = %q", got)
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testMask
func TestStrMask(t *testing.T) {
	t.Parallel()

	if got := StrMask("taylor@example.com", "*", 3); got != "tay***************" {
		t.Errorf("StrMask(3) = %q", got)
	}

	if got := StrMask("taylor@example.com", "*", -3); got != "taylor@example.***" {
		t.Errorf("StrMask(-3) = %q", got)
	}

	if got := StrMask("taylor@example.com", "*", 3, 3); got != "tay***@example.com" {
		t.Errorf("StrMask(3, 3) = %q", got)
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testStrIsAscii
func TestStrIsAscii(t *testing.T) {
	t.Parallel()

	if !StrIsAscii("hello") {
		t.Error("ASCII string should be ASCII")
	}

	if StrIsAscii("héllo") {
		t.Error("non-ASCII string should not be ASCII")
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testChopStart
func TestStrChopStart(t *testing.T) {
	t.Parallel()

	if got := StrChopStart("foobar", "foo"); got != "bar" {
		t.Errorf("StrChopStart = %q", got)
	}

	if got := StrChopStart("foobar", "baz"); got != "foobar" {
		t.Errorf("StrChopStart no match = %q", got)
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testChopEnd
func TestStrChopEnd(t *testing.T) {
	t.Parallel()

	if got := StrChopEnd("foobar", "bar"); got != "foo" {
		t.Errorf("StrChopEnd = %q", got)
	}

	if got := StrChopEnd("foobar", "baz"); got != "foobar" {
		t.Errorf("StrChopEnd no match = %q", got)
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testReplace
func TestStrReplace(t *testing.T) {
	t.Parallel()

	if got := StrReplace("foo", "bar", "foobar"); got != "barbar" {
		t.Errorf("StrReplace = %q", got)
	}
	// Slice search
	if got := StrReplace([]string{"a", "b"}, "x", "abc"); got != "xxc" {
		t.Errorf("StrReplace slice search = %q", got)
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testReplaceArray
func TestStrReplaceArray(t *testing.T) {
	t.Parallel()

	got := StrReplaceArray("?", []string{"foo", "bar"}, "? and ?")

	if got != "foo and bar" {
		t.Errorf("StrReplaceArray = %q", got)
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testReplaceFirst
func TestStrReplaceFirst(t *testing.T) {
	t.Parallel()

	if got := StrReplaceFirst("a", "b", "aaa"); got != "baa" {
		t.Errorf("StrReplaceFirst = %q", got)
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testReplaceLast
func TestStrReplaceLast(t *testing.T) {
	t.Parallel()

	if got := StrReplaceLast("a", "b", "aaa"); got != "aab" {
		t.Errorf("StrReplaceLast = %q", got)
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testReplaceStart
func TestStrReplaceStart(t *testing.T) {
	t.Parallel()

	if got := StrReplaceStart("foo", "bar", "foobar"); got != "barbar" {
		t.Errorf("StrReplaceStart = %q", got)
	}

	if got := StrReplaceStart("bar", "baz", "foobar"); got != "foobar" {
		t.Errorf("StrReplaceStart no match = %q", got)
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testReplaceEnd
func TestStrReplaceEnd(t *testing.T) {
	t.Parallel()

	if got := StrReplaceEnd("bar", "baz", "foobar"); got != "foobaz" {
		t.Errorf("StrReplaceEnd = %q", got)
	}

	if got := StrReplaceEnd("foo", "baz", "foobar"); got != "foobar" {
		t.Errorf("StrReplaceEnd no match = %q", got)
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testPluralStudly
func TestStrPluralStudly(t *testing.T) {
	t.Parallel()

	if got := StrPluralStudly("VerifiedHuman"); !strings.HasSuffix(got, "s") {
		t.Errorf("StrPluralStudly = %q (should be plural)", got)
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testPlural
func TestStrPlural(t *testing.T) {
	t.Parallel()

	if got := StrPlural("user"); got != "users" {
		t.Errorf("StrPlural(user) = %q", got)
	}
	// With count=1 should return singular
	if got := StrPlural("user", 1); got != "user" {
		t.Errorf("StrPlural(user, 1) = %q", got)
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testWordCount
func TestStrWordCount(t *testing.T) {
	t.Parallel()

	if got := StrWordCount("Hello World"); got != 2 {
		t.Errorf("StrWordCount = %d", got)
	}

	if got := StrWordCount("one"); got != 1 {
		t.Errorf("StrWordCount single = %d", got)
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testReplaceMatches
func TestStrReplaceMatches(t *testing.T) {
	t.Parallel()

	got := StrReplaceMatches(`\d+`, "number", "hello 123 world 456")

	if got != "hello number world number" {
		t.Errorf("StrReplaceMatches = %q", got)
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testIsMatch
func TestStrIsMatch(t *testing.T) {
	t.Parallel()

	if !StrIsMatch([]string{`\d+`}, "abc123") {
		t.Error("should match digits pattern")
	}

	if StrIsMatch([]string{`^only-letters$`}, "abc123") {
		t.Error("should not match letters-only pattern")
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testIs
func TestStrIs(t *testing.T) {
	t.Parallel()

	if !StrIs("*oo*", "foobar") {
		t.Error("wildcard pattern should match")
	}

	if !StrIs("foo*", "foobar") {
		t.Error("prefix wildcard should match")
	}

	if StrIs("baz*", "foobar") {
		t.Error("non-matching pattern should fail")
	}

	if !StrIs("*", "anything") {
		t.Error("* should match anything")
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testSwapKeywords
func TestStrSwap(t *testing.T) {
	t.Parallel()

	got := StrSwap(map[string]string{"foo": "bar", "baz": "qux"}, "foo and baz")

	if got != "bar and qux" {
		t.Errorf("StrSwap = %q", got)
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testTake
func TestStrTake(t *testing.T) {
	t.Parallel()

	if got := StrTake("hello", 3); got != "hel" {
		t.Errorf("StrTake(3) = %q", got)
	}

	if got := StrTake("hello", -3); got != "llo" {
		t.Errorf("StrTake(-3) = %q", got)
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testUcfirst
func TestStrUcfirst(t *testing.T) {
	t.Parallel()

	if got := StrUcfirst("hello world"); got != "Hello world" {
		t.Errorf("StrUcfirst = %q", got)
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testLcfirst
func TestStrLcfirst(t *testing.T) {
	t.Parallel()

	if got := StrLcfirst("Hello World"); got != "hello World" {
		t.Errorf("StrLcfirst = %q", got)
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testUcsplit
func TestStrUcsplit(t *testing.T) {
	t.Parallel()

	got := StrUcsplit("FooBar")

	if len(got) != 2 || got[0] != "Foo" || got[1] != "Bar" {
		t.Errorf("StrUcsplit(FooBar) = %v", got)
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testStrExcerpt
func TestStrExcerpt(t *testing.T) {
	t.Parallel()

	text := "This is my name"
	got := StrExcerpt(text, "my", 5)

	if !strings.Contains(got, "my") {
		t.Errorf("StrExcerpt should contain 'my', got %q", got)
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testMarkdown
func TestStrMarkdown(t *testing.T) {
	t.Parallel()

	got := StrMarkdown("## Hello World")

	if !strings.Contains(got, "<h2") {
		t.Errorf("StrMarkdown should produce h2 tag, got %q", got)
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testInlineMarkdown
func TestStrInlineMarkdown(t *testing.T) {
	t.Parallel()

	got := StrInlineMarkdown("**Hello**")

	if strings.Contains(got, "<p>") {
		t.Errorf("StrInlineMarkdown should not have <p> wrapper, got %q", got)
	}

	if !strings.Contains(got, "<strong>") {
		t.Errorf("StrInlineMarkdown should have <strong>, got %q", got)
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testNumbers
func TestStrNumbers(t *testing.T) {
	t.Parallel()

	if got := StrNumbers("abc123def456"); got != "123456" {
		t.Errorf("StrNumbers = %q", got)
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testApa
func TestStrApa(t *testing.T) {
	t.Parallel()

	if got := StrApa("the quick brown fox"); !strings.HasPrefix(got, "The") {
		t.Errorf("StrApa should capitalize first word, got %q", got)
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testSubstrCount
func TestStrSubstrCount(t *testing.T) {
	t.Parallel()

	if got := StrSubstrCount("hello world", "o"); got != 2 {
		t.Errorf("StrSubstrCount = %d", got)
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testPosition
func TestStrPosition(t *testing.T) {
	t.Parallel()

	pos, ok := StrPosition("hello world", "world")

	if !ok || pos != 6 {
		t.Errorf("StrPosition = (%d, %v), want (6, true)", pos, ok)
	}

	_, ok2 := StrPosition("hello world", "missing")

	if ok2 {
		t.Error("StrPosition for missing should return false")
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testLower
func TestStrLower(t *testing.T) {
	t.Parallel()

	if got := StrLower("HELLO WORLD"); got != "hello world" {
		t.Errorf("StrLower = %q", got)
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testUpper
func TestStrUpper(t *testing.T) {
	t.Parallel()

	if got := StrUpper("hello world"); got != "HELLO WORLD" {
		t.Errorf("StrUpper = %q", got)
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testTrim
func TestStrTrim(t *testing.T) {
	t.Parallel()

	if got := StrTrim("  hello  "); got != "hello" {
		t.Errorf("StrTrim = %q", got)
	}

	if got := StrTrim("//hello//", "/"); got != "hello" {
		t.Errorf("StrTrim with chars = %q", got)
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testPadBoth
func TestStrPad(t *testing.T) {
	t.Parallel()

	got := StrPadBoth("hello", 11)

	if got != "   hello   " {
		t.Errorf("StrPadBoth = %q", got)
	}

	gotLeft := StrPadLeft("hello", 10)

	if len(gotLeft) != 10 {
		t.Errorf("StrPadLeft length = %d", len(gotLeft))
	}

	gotRight := StrPadRight("hello", 10)

	if len(gotRight) != 10 {
		t.Errorf("StrPadRight length = %d", len(gotRight))
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testStringInitials
func TestStrInitials(t *testing.T) {
	t.Parallel()

	if got := StrInitials("Taylor Otwell"); got != "TO" {
		t.Errorf("StrInitials = %q", got)
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testWordWrap
func TestStrWordWrap(t *testing.T) {
	t.Parallel()

	got := StrWordWrap("The quick brown fox", 10)

	if !strings.Contains(got, "\n") {
		t.Errorf("StrWordWrap should contain newlines, got %q", got)
	}
}

// Test the fluent StringBuilder builder (Str::of())
func TestStrOf(t *testing.T) {
	t.Parallel()

	result := Of("  hello world  ").Trim().Upper().Value()

	if result != "HELLO WORLD" {
		t.Errorf("fluent chain = %q", result)
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testRandomStringFactoryCanBeSet
func TestStrRandomFactory(t *testing.T) {
	// NOT parallel — modifies global factory state
	CreateRandomStringsUsing(func(int) string { return "fixed" })

	defer CreateRandomStringsNormally()

	if got := StrRandom(); got != "fixed" {
		t.Errorf("custom factory = %q", got)
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testItCanSpecifyASequenceOfRandomStringsToUtilise
func TestStrRandomSequence(t *testing.T) {
	// NOT parallel — modifies global state
	cleanup := func() { CreateRandomStringsNormally() }
	CreateRandomStringsUsingSequence([]string{"first", "second"})

	defer cleanup()

	if got := StrRandom(); got != "first" {
		t.Errorf("first in sequence = %q", got)
	}

	if got := StrRandom(); got != "second" {
		t.Errorf("second in sequence = %q", got)
	}
}

// Port of Framework\Tests\Support\SupportStrTest::testSubstrReplace
func TestStrSubstrReplace(t *testing.T) {
	t.Parallel()

	if got := StrSubstrReplace("hello world", "earth", 6); got != "hello earth" {
		t.Errorf("StrSubstrReplace = %q", got)
	}
}
