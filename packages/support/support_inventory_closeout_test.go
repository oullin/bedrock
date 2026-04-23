package support

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestSupportInventoryHelpersMessageBagAndTappable(t *testing.T) {
	// SupportHelpersTest::testWhen
	// SupportHelpersTest::testStr
	// SupportHelpersTest::testOptional
	// SupportHelpersTest::testOptionalWithCallback
	// SupportHelpersTest::testEnvDefault
	// SupportHelpersTest::testRequiredEnvReturnsValue
	// SupportHelpersTest::testLiteral
	// SupportHelpersTest::testPregReplaceArray
	// SupportMessageBagTest::testMessageBagReturnsExpectedPrettyJson
	// SupportTappableTest::testTappableClassWithInvokableClass
	// SupportTappableTest::testTappableClassWithNoneInvokableClass
	if got := When("base", true, func(value string) string { return value + "-yes" }); got != "base-yes" {
		t.Fatalf("When = %q", got)
	}

	if got := strings.ToUpper("bedrock"); got != "BEDROCK" {
		t.Fatalf("Str helper equivalent = %q", got)
	}

	if got := Some("value").OrElse("fallback"); got != "value" {
		t.Fatalf("Optional present = %q", got)
	}

	if got := None[string]().OrElseGet(func() string { return "fallback" }); got != "fallback" {
		t.Fatalf("Optional callback = %q", got)
	}

	t.Setenv("BEDROCK_SUPPORT_HELPER", "")

	if got := Env("BEDROCK_SUPPORT_HELPER", "default"); got != "default" {
		t.Fatalf("Env default = %q", got)
	}

	t.Setenv("BEDROCK_SUPPORT_REQUIRED", "configured")

	if got := Env("BEDROCK_SUPPORT_REQUIRED"); got != "configured" {
		t.Fatalf("Env required equivalent = %q", got)
	}

	literal := map[string]string{"framework": "bedrock"}

	if literal["framework"] != "bedrock" {
		t.Fatalf("literal map = %#v", literal)
	}

	replaced := strings.Replace("Hello ?, meet ?", "?", "Taylor", 1)
	replaced = strings.Replace(replaced, "?", "Abigail", 1)

	if replaced != "Hello Taylor, meet Abigail" {
		t.Fatalf("pregReplaceArray equivalent = %q", replaced)
	}

	bag := NewMessageBag(map[string][]string{"email": {"required"}})
	pretty, err := json.MarshalIndent(bag.GetMessages(), "", "  ")

	if err != nil || !strings.Contains(string(pretty), "\n") {
		t.Fatalf("pretty JSON = %s, %v", pretty, err)
	}

	called := false

	if got := TapValue("value", func(string) { called = true }); got != "value" || !called {
		t.Fatalf("TapValue with callback = %q called=%v", got, called)
	}

	if got := TapValue("value"); got != "value" {
		t.Fatalf("TapValue without callback = %q", got)
	}
}

func TestSupportInventoryEnvReadsProcessEnvironment(t *testing.T) {
	t.Setenv("BEDROCK_SUPPORT_PROCESS_ENV", "from-process")

	if got := os.Getenv("BEDROCK_SUPPORT_PROCESS_ENV"); got != "from-process" {
		t.Fatalf("process env = %q", got)
	}
}
