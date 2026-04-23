package rules

import "testing"

// ValidationArrayRuleTest::testArrayValidation
// ValidationInArrayKeysTest::testInArrayKeysValidation
// ValidationRuleContainsTest::testContainsValidation
// ValidationRuleDoesntContainTest::testDoesntContainValidation
func TestArrayValidation(t *testing.T) {
	t.Parallel()

	if !validateArray("profile", map[string]any{
		"name":  "Taylor",
		"email": "user@example.com",
	}, []string{"name", "email"}, nil) {
		t.Fatal("expected array rule to accept permitted keys")
	}

	if validateArray("profile", map[string]any{
		"name":  "Taylor",
		"email": "user@example.com",
		"admin": true,
	}, []string{"name", "email"}, nil) {
		t.Fatal("expected array rule to reject disallowed keys")
	}

	if !validateRequiredArrayKeys("profile", map[string]any{
		"name":  "Taylor",
		"email": "user@example.com",
	}, []string{"name", "email"}, nil) {
		t.Fatal("expected required_array_keys to accept complete maps")
	}

	if validateRequiredArrayKeys("profile", map[string]any{
		"name": "Taylor",
	}, []string{"name", "email"}, nil) {
		t.Fatal("expected required_array_keys to reject missing keys")
	}

	if !validateInArrayKeys("profile", map[string]any{
		"name":  "Taylor",
		"email": "user@example.com",
	}, []string{"name", "email"}, nil) {
		t.Fatal("expected in_array_keys to accept permitted keys")
	}

	if validateInArrayKeys("profile", map[string]any{
		"name":  "Taylor",
		"email": "user@example.com",
		"admin": true,
	}, []string{"name", "email"}, nil) {
		t.Fatal("expected in_array_keys to reject disallowed keys")
	}

	if !validateContains("tags", []any{"alpha", "beta"}, []string{"alpha", "beta"}, nil) {
		t.Fatal("expected contains to accept listed values")
	}

	if validateContains("tags", []any{"alpha"}, []string{"alpha", "beta"}, nil) {
		t.Fatal("expected contains to reject missing values")
	}

	if !validateDoesntContain("tags", []any{"alpha", "beta"}, []string{"secret", "classified"}, nil) {
		t.Fatal("expected doesnt_contain to accept values outside the deny list")
	}

	if validateDoesntContain("tags", []any{"alpha", "secret"}, []string{"secret", "classified"}, nil) {
		t.Fatal("expected doesnt_contain to reject denied values")
	}
}
