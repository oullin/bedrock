package resources_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/httpx/resources"
)

type User struct {
	ID    int
	Name  string
	Email string
	Admin bool
}

func userMapper(u User, r *http.Request) map[string]any {
	return map[string]any{
		"id":    u.ID,
		"name":  u.Name,
		"email": u.Email,
	}
}

func TestJsonResourceToMap(t *testing.T) {
	t.Parallel()

	user := User{ID: 1, Name: "Taylor", Email: "taylor@example.com"}
	resource := resources.NewResource(user, userMapper)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	result := resource.ToMap(req)

	if result["name"] != "Taylor" {
		t.Fatalf("expected Taylor, got %v", result["name"])
	}

	if result["id"] != 1 {
		t.Fatalf("expected 1, got %v", result["id"])
	}
}

func TestJsonResourceToJSON(t *testing.T) {
	t.Parallel()

	user := User{ID: 1, Name: "Taylor", Email: "taylor@example.com"}
	resource := resources.NewResource(user, userMapper)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	b, err := resource.ToJSON(req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var data map[string]any

	json.Unmarshal(b, &data)

	// Default wrap key is "data".
	inner, ok := data["data"].(map[string]any)

	if !ok {
		t.Fatalf("expected data wrapper, got %v", data)
	}

	if inner["name"] != "Taylor" {
		t.Fatalf("expected Taylor, got %v", inner["name"])
	}
}

func TestJsonResourceWithoutWrapping(t *testing.T) {
	t.Parallel()

	user := User{ID: 1, Name: "Taylor"}
	resource := resources.NewResource(user, userMapper)
	resource.WrapKey = ""

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	b, _ := resource.ToJSON(req)

	var data map[string]any

	json.Unmarshal(b, &data)

	if data["name"] != "Taylor" {
		t.Fatal("expected unwrapped resource")
	}

	if _, ok := data["data"]; ok {
		t.Fatal("should not have data wrapper")
	}
}

func TestJsonResourceWithAdditionalData(t *testing.T) {
	t.Parallel()

	user := User{ID: 1, Name: "Taylor"}
	resource := resources.NewResource(user, userMapper)
	resource.With = map[string]any{"version": "1.0"}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	b, _ := resource.ToJSON(req)

	var data map[string]any

	json.Unmarshal(b, &data)

	if data["version"] != "1.0" {
		t.Fatalf("expected version 1.0, got %v", data["version"])
	}
}

func TestJsonResourceResponse(t *testing.T) {
	t.Parallel()

	user := User{ID: 1, Name: "Taylor"}
	resource := resources.NewResource(user, userMapper)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	err := resource.Response(rec, req, http.StatusOK)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	if rec.Header().Get("Content-Type") != "application/json" {
		t.Fatal("expected application/json content type")
	}
}

func TestJsonResourceConditionalValue(t *testing.T) {
	t.Parallel()

	user := User{ID: 1, Name: "Taylor", Admin: true}

	resource := resources.NewResource(user, func(u User, r *http.Request) map[string]any {
		return map[string]any{
			"id":       u.ID,
			"name":     u.Name,
			"is_admin": resources.When(u.Admin, true),
			"secret":   resources.When(false, "hidden"),
		}
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	result := resource.ToMap(req)

	if result["is_admin"] != true {
		t.Fatal("expected is_admin to be included")
	}

	if _, ok := result["secret"]; ok {
		t.Fatal("expected secret to be omitted")
	}
}

func TestJsonResourceUnless(t *testing.T) {
	t.Parallel()

	resource := resources.NewResource(User{ID: 1}, func(u User, r *http.Request) map[string]any {
		return map[string]any{
			"id":     u.ID,
			"hidden": resources.Unless(true, "should hide"),
			"shown":  resources.Unless(false, "should show"),
		}
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	result := resource.ToMap(req)

	if _, ok := result["hidden"]; ok {
		t.Fatal("expected hidden to be omitted")
	}

	if result["shown"] != "should show" {
		t.Fatal("expected shown to be included")
	}
}

func TestJsonResourceMergeWhen(t *testing.T) {
	t.Parallel()

	resource := resources.NewResource(User{ID: 1, Admin: true}, func(u User, r *http.Request) map[string]any {
		return map[string]any{
			"id": u.ID,
			"_admin": resources.MergeWhen(u.Admin, map[string]any{
				"role":        "admin",
				"permissions": 42,
			}),
		}
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	result := resource.ToMap(req)

	if result["role"] != "admin" {
		t.Fatal("expected role to be merged")
	}

	if result["permissions"] != 42 {
		t.Fatal("expected permissions to be merged")
	}

	if _, ok := result["_admin"]; ok {
		t.Fatal("expected _admin key to be consumed by merge")
	}
}

func TestJsonResourceMergeWhenFalse(t *testing.T) {
	t.Parallel()

	resource := resources.NewResource(User{ID: 1}, func(u User, r *http.Request) map[string]any {
		return map[string]any{
			"id":     u.ID,
			"_admin": resources.MergeWhen(false, map[string]any{"role": "admin"}),
		}
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	result := resource.ToMap(req)

	if _, ok := result["role"]; ok {
		t.Fatal("expected role to NOT be merged when false")
	}
}
