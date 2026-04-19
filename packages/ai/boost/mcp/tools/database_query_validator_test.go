package tools

import (
	"database/sql"
	"strings"
	"testing"
)

func TestValidateReadOnlySQL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		driver  string
		query   string
		wantErr string
	}{
		{
			name:   "postgres plain select",
			driver: "postgres",
			query:  "SELECT 1",
		},
		{
			name:   "postgres read only cte",
			driver: "postgres",
			query:  "WITH x AS (SELECT 1) SELECT * FROM x",
		},
		{
			name:   "postgres explain select",
			driver: "postgres",
			query:  "EXPLAIN SELECT 1",
		},
		{
			name:   "postgres explain cte select",
			driver: "postgres",
			query:  "EXPLAIN WITH x AS (SELECT 1) SELECT * FROM x",
		},
		{
			name:   "pgsql alias plain select",
			driver: "pgsql",
			query:  "SELECT 1",
		},
		{
			name:   "mysql plain select",
			driver: "mysql",
			query:  "SELECT 1",
		},
		{
			name:    "postgres writable cte",
			driver:  "postgres",
			query:   "WITH x AS (DELETE FROM users RETURNING id) SELECT * FROM x",
			wantErr: "read-only",
		},
		{
			name:    "postgres with insert",
			driver:  "postgres",
			query:   "WITH x AS (INSERT INTO users (name) VALUES ('a') RETURNING id) SELECT * FROM x",
			wantErr: "read-only",
		},
		{
			name:    "postgres explain delete",
			driver:  "postgres",
			query:   "EXPLAIN DELETE FROM users",
			wantErr: "read-only",
		},
		{
			name:    "postgres explain analyze delete",
			driver:  "postgres",
			query:   "EXPLAIN ANALYZE DELETE FROM users",
			wantErr: "read-only",
		},
		{
			name:    "mysql explain delete",
			driver:  "mysql",
			query:   "EXPLAIN DELETE FROM users",
			wantErr: "read-only",
		},
		{
			name:    "multi statement input",
			driver:  "postgres",
			query:   "SELECT 1; DELETE FROM users",
			wantErr: "single statement",
		},
		{
			name:    "unsupported driver",
			driver:  "sqlite",
			query:   "SELECT 1",
			wantErr: "unsupported driver",
		},
		{
			name:    "comment prefixed mutating statement",
			driver:  "postgres",
			query:   "/* leading */ DELETE FROM users",
			wantErr: "read-only",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := validateReadOnlySQL(tt.driver, tt.query)

			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("validateReadOnlySQL(%q, %q) error = %v, want nil", tt.driver, tt.query, err)
				}

				return
			}

			if err == nil {
				t.Fatalf("validateReadOnlySQL(%q, %q) error = nil, want %q", tt.driver, tt.query, tt.wantErr)
			}

			if !strings.Contains(strings.ToLower(err.Error()), strings.ToLower(tt.wantErr)) {
				t.Fatalf("validateReadOnlySQL(%q, %q) error = %q, want substring %q", tt.driver, tt.query, err.Error(), tt.wantErr)
			}
		})
	}
}

func TestDatabaseQueryHandleRejectsUnsafeSQLBeforeExecution(t *testing.T) {
	t.Parallel()

	tool := &DatabaseQuery{
		DB:     &sql.DB{},
		Driver: "postgres",
	}

	resp, err := tool.Handle(McpRequest{
		Args: map[string]any{
			"query": "EXPLAIN DELETE FROM users",
		},
	})

	if err != nil {
		t.Fatalf("Handle: %v", err)
	}

	if !resp.IsError {
		t.Fatal("expected error response for unsafe SQL")
	}

	if len(resp.Content) != 1 {
		t.Fatalf("expected one response content item, got %d", len(resp.Content))
	}

	if !strings.Contains(resp.Content[0].Text, "read-only") {
		t.Fatalf("response text = %q, want read-only validation error", resp.Content[0].Text)
	}
}
