package routing_test

import (
	"encoding/json"
	"testing"

	"github.com/bedrock/packages/routing"
)

func TestRegistryAdd(t *testing.T) {
	t.Parallel()

	reg := routing.NewRegistry()
	reg.Add("home", "GET", "/")
	reg.Add("users.show", "GET", "/users/{id}")

	if url := reg.URL("home", nil); url != "/" {
		t.Fatalf("expected /, got %q", url)
	}
}

func TestRegistryURLSubstitution(t *testing.T) {
	t.Parallel()

	reg := routing.NewRegistry()
	reg.Add("posts.show", "GET", "/posts/{id}/comments/{cid}")

	url := reg.URL("posts.show", map[string]string{"id": "5", "cid": "3"})

	if url != "/posts/5/comments/3" {
		t.Fatalf("unexpected url: %q", url)
	}
}

func TestRegistryUnknownRoute(t *testing.T) {
	t.Parallel()

	reg := routing.NewRegistry()
	url := reg.URL("no.such.route", nil)

	if url == "" || url[0] == '/' {
		t.Fatalf("expected fallback, got %q", url)
	}
}

func TestRegistryGroup(t *testing.T) {
	t.Parallel()

	reg := routing.NewRegistry()
	reg.Group("api", "/api", func(g *routing.RegistryGroup) {
		g.Add("users.index", "GET", "/users")
		g.Add("users.show", "GET", "/users/{id}")
	})

	url := reg.URL("api.users.show", map[string]string{"id": "42"})

	if url != "/api/users/42" {
		t.Fatalf("expected /api/users/42, got %q", url)
	}
}

func TestRegistryManifest(t *testing.T) {
	t.Parallel()

	reg := routing.NewRegistry()
	reg.Add("home", "GET", "/")
	reg.Add("about", "GET", "/about")

	m := reg.Manifest()

	if m["home"] != "/" || m["about"] != "/about" {
		t.Fatalf("unexpected manifest: %v", m)
	}
}

func TestRegistryToJSON(t *testing.T) {
	t.Parallel()

	reg := routing.NewRegistry()
	reg.Add("home", "GET", "/")

	b, err := reg.ToJSON()

	if err != nil {
		t.Fatal(err)
	}

	var routes []routing.Route

	if err := json.Unmarshal(b, &routes); err != nil {
		t.Fatal(err)
	}

	if len(routes) != 1 || routes[0].Name != "home" {
		t.Fatalf("unexpected routes: %v", routes)
	}
}

func TestRouteParams(t *testing.T) {
	t.Parallel()

	r := routing.Route{Name: "x", Method: "GET", Pattern: "/a/{id}/b/{slug}"}
	params := r.Params()

	if len(params) != 2 || params[0] != "id" || params[1] != "slug" {
		t.Fatalf("unexpected params: %v", params)
	}
}

func TestRegistryExportOrder(t *testing.T) {
	t.Parallel()

	reg := routing.NewRegistry()
	reg.Add("first", "GET", "/first")
	reg.Add("second", "GET", "/second")
	reg.Add("third", "GET", "/third")

	routes := reg.Export()

	if len(routes) != 3 || routes[0].Name != "first" || routes[2].Name != "third" {
		t.Fatalf("unexpected order: %v", routes)
	}
}

func TestRegistryLookup(t *testing.T) {
	t.Parallel()

	reg := routing.NewRegistry()
	reg.Add("users.show", "GET", "/users/{id}")

	r, ok := reg.Lookup("users.show")

	if !ok || r.Pattern != "/users/{id}" {
		t.Fatalf("unexpected lookup: %v %v", ok, r)
	}

	_, ok = reg.Lookup("missing")

	if ok {
		t.Fatal("expected not found")
	}
}
