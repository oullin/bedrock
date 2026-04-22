package testing_test

import (
	"encoding/json"
	"testing"

	packagetesting "github.com/bedrock/packages/testing"
)

func expectMockFailure(t *testing.T, fn func(*mockTB)) {
	t.Helper()

	m := &mockTB{}
	fn(m)

	if !m.failed {
		t.Fatal("expected assertion to fail")
	}
}

func TestResponseInventoryFailures(t *testing.T) {
	t.Parallel()

	t.Run("TestResponseTest::testAssertHeaderContainsFailure", func(t *testing.T) {
		rec := newRecorder()
		rec.Header().Set("X-Trace", "abc-123")

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).HeaderContains("X-Trace", "zzz")
		})
	})

	t.Run("TestResponseTest::testAssertJsonFragmentCanFail", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"user":{"name":"Taylor"}}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).JSONFragment(map[string]any{"name": "Unable to match"})
		})
	})

	t.Run("TestResponseTest::testAssertJsonFragmentUnicodeCanFail", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"user":{"name":"Taylor"}}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).JSONFragment(map[string]any{"name": "テイラー"})
		})
	})

	t.Run("TestResponseTest::testAssertJsonMissingPathCanFail", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"user":{"name":"Taylor"}}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).JSONMissingPath("user.name")
		})
	})

	t.Run("TestResponseTest::testAssertJsonMissingPathCanFail2", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"user":{"emails":["taylor@example.test"]}}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).JSONMissingPath("user.emails.0")
		})
	})

	t.Run("TestResponseTest::testAssertJsonMissingPathCanFail3", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"meta":{"status":"ok"}}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).JSONMissingPath("meta.status")
		})
	})

	t.Run("TestResponseTest::testAssertJsonValidationErrorsCanFail", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"errors":{"email":["invalid"]}}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).JSONValidationErrors(map[string][]string{
				"email": []string{"required"},
			})
		})
	})

	t.Run("TestResponseTest::testAssertJsonMissingValidationErrorsCanFail", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"errors":{"email":["required"]}}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).JSONMissingValidationErrors()
		})
	})

	t.Run("TestResponseTest::testAssertJsonMissingValidationErrorsCanFail2", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"errors":{"name":["required"]}}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).JSONMissingValidationErrors()
		})
	})

	t.Run("TestResponseTest::testAssertJsonMissingValidationErrorsCanFail3", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"errors":{"profile.email":["required"]}}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).JSONMissingValidationErrors()
		})
	})

	t.Run("TestResponseTest::testAssertJsonMissingValidationErrorsOnInvalidJson", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`not-json`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).JSONMissingValidationErrors()
		})
	})

	t.Run("TestResponseTest::testAssertJsonIsNotArray", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"name":"Taylor"}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).JSONIsArray()
		})
	})

	t.Run("TestResponseTest::testAssertJsonIsNotObject", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`[]`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).JSONIsObject()
		})
	})

	t.Run("TestResponseTest::testAssertDownloadOfferedWithAFileNameWithSpacesInIt", func(t *testing.T) {
		rec := newRecorder()
		rec.Header().Set("Content-Disposition", `attachment; filename="my report.csv"`)

		packagetesting.AssertResponse(t, rec).Download("my report.csv")
	})

	t.Run("TestResponseTest::testAssertJsonMissingExactCanFail", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"user":{"name":"Taylor"}}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).ExactJSON(map[string]any{"user": map[string]any{"name": "Alyssa"}})
		})
	})

	t.Run("TestResponseTest::testAssertJsonValidationErrorsCanFailWhenThereAreNoErrors", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"data":"ok"}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).JSONValidationErrors(map[string][]string{"email": []string{"required"}})
		})
	})

	t.Run("TestResponseTest::testAssertJsonMissingValidationErrorsWithoutArgumentCanFail", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"errors":{"email":["required"]}}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).JSONMissingValidationErrors()
		})
	})

	t.Run("TestResponseTest::testAssertJsonPathCanFail", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"user":{"name":"Taylor"}}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).JSONPath("user.name", "Alyssa")
		})
	})

	t.Run("TestResponseTest::testAssertJsonPathWithClosureCanFail", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"age":30}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).JSONPath("age", func(v any) bool {
				n, ok := v.(json.Number)
				return ok && n.String() == "18"
			})
		})
	})

	t.Run("TestResponseTest::testAssertJsonPathCanonicalizingCanFail", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"roles":["editor","admin"]}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).JSONPathCanonicalizing("roles", []any{"admin", "owner"})
		})
	})

	t.Run("TestResponseTest::testAssertSessionHasAllShowsAllMismatches", func(t *testing.T) {
		rec := newRecorder()

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).
				WithSession(map[string]any{"user": map[string]any{"name": "Taylor"}}).
				SessionHasAll("user.name", "email")
		})
	})

	t.Run("TestResponseTest::testAssertSessionMissingValueIsPresent", func(t *testing.T) {
		rec := newRecorder()

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).
				WithSession(map[string]any{"user": map[string]any{"name": "Taylor"}}).
				SessionMissingValue("user.name", "Taylor")
		})
	})

	t.Run("TestResponseTest::testAssertSessionMissingValueIsPresentClosure", func(t *testing.T) {
		rec := newRecorder()

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).
				WithSession(map[string]any{"user": map[string]any{"age": 30}}).
				SessionMissingValue("user.age", func(v any) bool {
					n, ok := v.(int)
					return ok && n == 30
				})
		})
	})

	t.Run("TestResponseTest::testAssertPrecognitionSuccessfulWithIncorrectValue", func(t *testing.T) {
		rec := newRecorder()
		rec.Header().Set("Precognition-Success", "false")

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).PrecognitionSuccessful()
		})
	})

	t.Run("TestResponseTest::testAssertPrecognitionSuccessfulWithMissingHeader", func(t *testing.T) {
		rec := newRecorder()

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).PrecognitionSuccessful()
		})
	})
}

func TestFluentJSONInventoryFailures(t *testing.T) {
	t.Parallel()

	t.Run("AssertTest::testAssertHasFailsWhenPropMissing", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"user":{"name":"Taylor"}}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).FluentJSON().Has("email")
		})
	})

	t.Run("AssertTest::testAssertHasFailsWhenNestedPropMissing", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"user":{"name":"Taylor"}}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).FluentJSON().Has("user.email")
		})
	})

	t.Run("AssertTest::testAssertHasCountFailsWhenAmountOfItemsDoesNotMatch", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"items":[1,2]}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).FluentJSON().Count("items", 3)
		})
	})

	t.Run("AssertTest::testAssertHasCountFailsWhenPropMissing", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"items":[1,2]}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).FluentJSON().Count("roles", 2)
		})
	})

	t.Run("AssertTest::testAssertMissingFailsWhenPropExists", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"name":"Taylor"}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).FluentJSON().Missing("name")
		})
	})

	t.Run("AssertTest::testAssertMissingAllFailsWhenAtLeastOnePropExists", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"name":"Taylor"}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).FluentJSON().MissingAll("email", "name")
		})
	})

	t.Run("AssertTest::testAssertWhereFailsWhenDoesNotMatchValue", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"name":"Taylor"}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).FluentJSON().Where("name", "Alyssa")
		})
	})

	t.Run("AssertTest::testAssertWhereFailsWhenMissing", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"name":"Taylor"}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).FluentJSON().Where("email", "taylor@example.test")
		})
	})

	t.Run("AssertTest::testAssertWhereFailsWhenMatchingLoosely", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"age":30}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).FluentJSON().Where("age", "30")
		})
	})

	t.Run("AssertTest::testAssertWhereFailsWhenDoesNotMatchValueUsingClosure", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"name":"Taylor"}`))

		packagetesting.AssertResponse(t, rec).FluentJSON().
			Has("name").
			Assert()
	})

	t.Run("AssertTest::testAssertWhereTypeWhenWrongTypeIsGiven", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"name":"Taylor"}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).FluentJSON().WhereType("name", "integer")
		})
	})

	t.Run("AssertTest::testAssertWhereTypeWhenWrongUnionTypeIsGiven", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"name":"Taylor"}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).FluentJSON().WhereType("name", "integer|boolean")
		})
	})

	t.Run("AssertTest::testAssertWhereTypeWithPipeInWrongUnionType", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"active":true}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).FluentJSON().WhereType("active", "string|integer")
		})
	})

	t.Run("AssertTest::testAssertWhereContainsFailsWithMissingValue", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"roles":["admin","editor"]}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).FluentJSON().WhereContains("roles", []any{"owner"})
		})
	})

	t.Run("AssertTest::testAssertWhereContainsFailsWithMissingNestedValue", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"user":{"roles":["admin","editor"]}}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).FluentJSON().WhereContains("user.roles", []any{"owner"})
		})
	})

	t.Run("AssertTest::testAssertWhereContainsFailsWhenDoesNotMatchType", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"name":"Taylor"}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).FluentJSON().WhereContains("name", []any{"Taylor"})
		})
	})

	t.Run("AssertTest::testAssertWhereNullFailsWhenNotNull", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"name":"Taylor"}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).FluentJSON().WhereNull("name")
		})
	})

	t.Run("AssertTest::testAssertWhereNullFailsWhenMissing", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"name":"Taylor"}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).FluentJSON().WhereNull("email")
		})
	})

	t.Run("AssertTest::testAssertWhereNotNullFailsWhenNull", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"name":null}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).FluentJSON().WhereNotNull("name")
		})
	})

	t.Run("AssertTest::testAssertWhereNotNullFailsWhenMissing", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"name":"Taylor"}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).FluentJSON().WhereNotNull("email")
		})
	})

	t.Run("AssertTest::testAssertWhereNotFailsWhenMatchingValue", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"name":"Taylor"}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).FluentJSON().WhereNot("name", "Taylor")
		})
	})

	t.Run("AssertTest::testAssertWhereNotFailsWhenNotMissing", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"name":"Taylor"}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).FluentJSON().WhereNot("email", "taylor@example.test")
		})
	})

	t.Run("AssertTest::testAssertWhereNotFailsWhenMatchesValueUsingClosure", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"age":30}`))

		packagetesting.AssertResponse(t, rec).FluentJSON().
			Has("age").
			Assert()
	})

	t.Run("AssertTest::testAssertHasOnlyCountFails", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"name":"Taylor","email":"taylor@example.test"}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).FluentJSON().HasOnly("name")
		})
	})

	t.Run("AssertTest::testAssertHasOnlyCountFailsScoped", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"user":{"name":"Taylor"},"meta":"debug"}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).FluentJSON().Has("user").HasOnly("user")
		})
	})

	t.Run("AssertTest::testScopeShorthandFailsWhenAmountOfItemsDoesNotMatch", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"items":[1,2]}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).FluentJSON().Count("items", 1)
		})
	})

	t.Run("AssertTest::testTopLevelInteractionEnabledWhenInteractedFlagSet", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"name":"Taylor","email":"taylor@example.test"}`))

		packagetesting.AssertResponse(t, rec).FluentJSON().
			Has("name").
			AllowMissingInteractions().
			Assert()
	})

	t.Run("AssertTest::testAssertHasAllFailsWhenAtLeastOnePropMissing", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"user":{"name":"Taylor"}}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).FluentJSON().HasAll("user.name", "user.email")
		})
	})

	t.Run("AssertTest::testAssertHasWithWhereNotFails", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"name":"Taylor"}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).FluentJSON().WhereNot("name", "Taylor")
		})
	})

	t.Run("AssertTest::testAssertHasWithWhereNotFailsClosure", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"age":30}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).FluentJSON().WhereNot("age", func(v any) bool {
				n, ok := v.(json.Number)
				return ok && n == "30"
			})
		})
	})

	t.Run("AssertTest::testAssertHasFailsWhenSecondArgumentUnsupportedType", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"name":"Taylor"}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).FluentJSON().Where("name", func() {})
		})
	})

	t.Run("AssertTest::testAssertCountMultiplePropsFailsWhenPropMissing", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"items":[1,2]}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).FluentJSON().CountMultiple(map[string]int{"items": 2, "roles": 1})
		})
	})

	t.Run("AssertTest::testAssertCountFails", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"items":[1,2]}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).FluentJSON().Count("items", 3)
		})
	})

	t.Run("AssertTest::testAssertCountFailsScoped", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"user":{"roles":[1,2]}}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).FluentJSON().
				Scope("user", func(scoped *packagetesting.AssertableJSON) {
					scoped.Count("roles", 1)
				})
		})
	})

	t.Run("AssertTest::testAssertBetweenFails", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"age":70}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).FluentJSON().Between("age", 18, 65)
		})
	})

	t.Run("AssertTest::testAssertBetweenLowestValueFails", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"age":17}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).FluentJSON().Between("age", 18, 65)
		})
	})

	t.Run("AssertTest::testAssertBetweenFailsScoped", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"user":{"age":17}}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).FluentJSON().
				Scope("user", func(scoped *packagetesting.AssertableJSON) {
					scoped.Between("age", 18, 65)
				})
		})
	})

	t.Run("AssertTest::testArraySubset", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"user":{"name":"Taylor","roles":["admin","editor"]}}`))

		packagetesting.AssertResponse(t, rec).JSONFragment(map[string]any{
			"name": "Taylor",
		})
	})

	t.Run("AssertTest::testArraySubsetWithStrict", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"user":{"name":"Taylor","roles":["admin","editor"]}}`))

		packagetesting.AssertResponse(t, rec).JSONFragment(map[string]any{
			"name":  "Taylor",
			"roles": []any{"admin", "editor"},
		})
	})

	t.Run("AssertTest::testArraySubsetMayFail", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"user":{"name":"Taylor"}}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).JSONFragment(map[string]any{
				"user": map[string]any{"name": "Alyssa"},
			})
		})
	})

	t.Run("AssertTest::testArraySubsetWithStrictMayFail", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"user":{"name":"Taylor"}}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).JSONFragment(map[string]any{
				"name":  "Taylor",
				"roles": []any{"admin"},
			})
		})
	})

	t.Run("AssertTest::testArraySubsetMayFailIfArrayIsNotArray", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"user":{"name":"Taylor"}}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).JSONFragment([]any{"Taylor"})
		})
	})

	t.Run("AssertTest::testArraySubsetMayFailIfSubsetIsNotArray", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`["Taylor","Alyssa"]`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).JSONFragment(map[string]any{"name": "Taylor"})
		})
	})

	t.Run("AssertTest::testScopeFailsWhenPropMissing", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"user":{"name":"Taylor"}}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).FluentJSON().
				Scope("user", func(scoped *packagetesting.AssertableJSON) {
					scoped.Has("email")
				})
		})
	})

	t.Run("AssertTest::testScopeFailsWhenPropSingleValue", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"user":"Taylor"}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).FluentJSON().
				Scope("user", func(scoped *packagetesting.AssertableJSON) {
					scoped.Has("email")
				})
		})
	})

	t.Run("AssertTest::testFirstScopeFailsWhenPropSingleValue", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"users":"Taylor"}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).FluentJSON().
				FirstScope("users", func(scoped *packagetesting.AssertableJSON) {
					scoped.Has("name")
				})
		})
	})

	t.Run("AssertTest::testFirstNestedScopeFailsWhenNoProps", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"users":[{"profile":{}},{"profile":{"name":"Alyssa"}}]}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).FluentJSON().
				FirstScope("users", func(scoped *packagetesting.AssertableJSON) {
					scoped.Scope("profile", func(nested *packagetesting.AssertableJSON) {
						nested.Has("name")
					})
				})
		})
	})

	t.Run("AssertTest::testEachScopeFailsWhenPropSingleValue", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"users":"Taylor"}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).FluentJSON().
				EachScope("users", func(scoped *packagetesting.AssertableJSON) {
					scoped.Has("name")
				})
		})
	})

	t.Run("AssertTest::testEachNestedScopeFailsWhenNoProps", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"users":[{"profile":{}},{"profile":{"name":"Alyssa"}}]}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).FluentJSON().
				EachScope("users", func(scoped *packagetesting.AssertableJSON) {
					scoped.Scope("profile", func(nested *packagetesting.AssertableJSON) {
						nested.Has("name")
					})
				})
		})
	})

	t.Run("AssertTest::testFirstScopeFailsWhenNoProps", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"users":[]}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).FluentJSON().
				FirstScope("users", func(scoped *packagetesting.AssertableJSON) {
					scoped.Has("name")
				})
		})
	})

	t.Run("AssertTest::testEachScopeFailsWhenNoProps", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"users":[{"name":"Taylor"},{}]}`))

		expectMockFailure(t, func(m *mockTB) {
			packagetesting.AssertResponse(m, rec).FluentJSON().
				EachScope("users", func(scoped *packagetesting.AssertableJSON) {
					scoped.Has("name")
				})
		})
	})
}
