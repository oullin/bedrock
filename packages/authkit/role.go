package authkit

import "sync"

// Role defines a named set of permissions within a team.
type Role struct {
	Key         string
	Name        string
	Description string
	Permissions []string
}

// HasPermission reports whether the role includes the given permission.

// RoleRegistry holds defined roles and the default role for new members.
type RoleRegistry struct {
	mu          sync.RWMutex
	roles       map[string]*Role
	defaultRole string
}

func (r *Role) HasPermission(permission string) bool {
	for _, p := range r.Permissions {
		if p == permission {
			return true
		}
	}

	return false
}

// NewRoleRegistry creates an empty role registry.
func NewRoleRegistry() *RoleRegistry {
	return &RoleRegistry{
		roles: make(map[string]*Role),
	}
}

// Define registers a role with the given key, name, description, and permissions.
func (rr *RoleRegistry) Define(key string, name string, description string, permissions []string) *Role {
	rr.mu.Lock()

	defer rr.mu.Unlock()

	role := &Role{
		Key:         key,
		Name:        name,
		Description: description,
		Permissions: permissions,
	}

	rr.roles[key] = role

	return role
}

// SetDefault sets the default role key assigned to new team members.
func (rr *RoleRegistry) SetDefault(key string) {
	rr.mu.Lock()

	defer rr.mu.Unlock()

	rr.defaultRole = key
}

// Default returns the default role key.
func (rr *RoleRegistry) Default() string {
	rr.mu.RLock()

	defer rr.mu.RUnlock()

	return rr.defaultRole
}

// Find returns the role for the given key, or nil if not found.
func (rr *RoleRegistry) Find(key string) *Role {
	rr.mu.RLock()

	defer rr.mu.RUnlock()

	return rr.roles[key]
}

// All returns all defined roles.
func (rr *RoleRegistry) All() []*Role {
	rr.mu.RLock()

	defer rr.mu.RUnlock()

	result := make([]*Role, 0, len(rr.roles))

	for _, r := range rr.roles {
		result = append(result, r)
	}

	return result
}
