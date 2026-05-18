package grammars_test

import (
	"strings"
	"testing"

	"github.com/bedrock/packages/database/query"
	"github.com/bedrock/packages/database/query/grammars"
)

// These tests cover the happy-path (string, error) returns of
// CompileInsertGetId / CompileUpsert / CompileInsertOrIgnore for every
// non-ClickHouse grammar. ClickHouse's error path is covered by the
// dedicated tests in clickhouse_grammar_test.go.

type compileSiblingGrammar interface {
	CompileInsertGetId(b *query.Builder, values map[string]any, sequence string) (string, error)
	CompileUpsert(b *query.Builder, values []map[string]any, uniqueBy []string, update []string) (string, error)
	CompileInsertOrIgnore(b *query.Builder, values []map[string]any) (string, error)
}

func TestSiblingGrammars_CompileReturnsNoError(t *testing.T) {
	t.Parallel()

	type expect struct {
		insertGetIdContains    string
		upsertContains         string
		insertOrIgnoreContains string
	}

	cases := []struct {
		name    string
		grammar compileSiblingGrammar
		exp     expect
	}{
		{
			name:    "postgres",
			grammar: grammars.NewPostgresGrammar(),
			exp: expect{
				insertGetIdContains:    "returning",
				upsertContains:         "on conflict",
				insertOrIgnoreContains: "on conflict do nothing",
			},
		},
		{
			name:    "mysql",
			grammar: grammars.NewMySQLGrammar(),
			exp: expect{
				insertGetIdContains:    "insert into",
				upsertContains:         "on duplicate key update",
				insertOrIgnoreContains: "insert ignore",
			},
		},
		{
			name:    "mariadb",
			grammar: grammars.NewMariaDBGrammar(),
			exp: expect{
				insertGetIdContains:    "returning",
				upsertContains:         "on duplicate key update",
				insertOrIgnoreContains: "insert ignore",
			},
		},
		{
			name:    "sqlite",
			grammar: grammars.NewSQLiteGrammar(),
			exp: expect{
				insertGetIdContains:    "insert into",
				upsertContains:         "on conflict",
				insertOrIgnoreContains: "insert or ignore",
			},
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			b := query.NewBuilder(nil, tc.grammar.(query.Grammar), nil).From("events")

			sql, err := tc.grammar.CompileInsertGetId(b, map[string]any{"name": "x"}, "id")

			if err != nil {
				t.Fatalf("InsertGetId returned unexpected error: %v", err)
			}

			if !strings.Contains(strings.ToLower(sql), tc.exp.insertGetIdContains) {
				t.Errorf("InsertGetId: expected SQL to contain %q, got %q", tc.exp.insertGetIdContains, sql)
			}

			sql, err = tc.grammar.CompileUpsert(b,
				[]map[string]any{{"id": 1, "name": "x"}},
				[]string{"id"},
				[]string{"name"})

			if err != nil {
				t.Fatalf("Upsert returned unexpected error: %v", err)
			}

			if !strings.Contains(strings.ToLower(sql), tc.exp.upsertContains) {
				t.Errorf("Upsert: expected SQL to contain %q, got %q", tc.exp.upsertContains, sql)
			}

			sql, err = tc.grammar.CompileInsertOrIgnore(b, []map[string]any{{"id": 1}})

			if err != nil {
				t.Fatalf("InsertOrIgnore returned unexpected error: %v", err)
			}

			if !strings.Contains(strings.ToLower(sql), tc.exp.insertOrIgnoreContains) {
				t.Errorf("InsertOrIgnore: expected SQL to contain %q, got %q", tc.exp.insertOrIgnoreContains, sql)
			}
		})
	}
}
