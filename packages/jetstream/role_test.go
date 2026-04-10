package jetstream

import "testing"

func TestRoleHasPermission(t *testing.T) {
	role := &Role{
		Key:         "admin",
		Name:        "Administrator",
		Permissions: []string{"create", "read", "update", "delete"},
	}

	if !role.HasPermission("create") {
		t.Fatal("admin should have create permission")
	}

	if role.HasPermission("deploy") {
		t.Fatal("admin should not have deploy permission")
	}
}

func TestRoleRegistryDefine(t *testing.T) {
	rr := NewRoleRegistry()

	admin := rr.Define("admin", "Administrator", "Full access", []string{"*"})

	if admin.Key != "admin" {
		t.Fatalf("expected key admin, got %s", admin.Key)
	}

	editor := rr.Define("editor", "Editor", "Can edit", []string{"read", "update"})

	if editor.Name != "Editor" {
		t.Fatalf("expected name Editor, got %s", editor.Name)
	}

	all := rr.All()

	if len(all) != 2 {
		t.Fatalf("expected 2 roles, got %d", len(all))
	}
}

func TestRoleRegistryFind(t *testing.T) {
	rr := NewRoleRegistry()
	rr.Define("admin", "Admin", "", []string{"*"})

	found := rr.Find("admin")

	if found == nil {
		t.Fatal("expected to find admin role")
	}

	notFound := rr.Find("nonexistent")

	if notFound != nil {
		t.Fatal("expected nil for nonexistent role")
	}
}

func TestRoleRegistryDefault(t *testing.T) {
	rr := NewRoleRegistry()
	rr.Define("member", "Member", "", []string{"read"})
	rr.SetDefault("member")

	if rr.Default() != "member" {
		t.Fatalf("expected default member, got %s", rr.Default())
	}
}

func TestRoleRegistryEmptyDefault(t *testing.T) {
	rr := NewRoleRegistry()

	if rr.Default() != "" {
		t.Fatalf("expected empty default, got %s", rr.Default())
	}
}
