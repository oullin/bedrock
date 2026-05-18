package query_test

import (
	"context"
	"errors"
	"testing"

	dbcontract "github.com/bedrock/packages/contracts/database"
	"github.com/bedrock/packages/database"
	"github.com/bedrock/packages/database/query"
	"github.com/bedrock/packages/database/query/grammars"
)

// stubBuilderConn satisfies query.ConnectionInterface for tests where the
// sentinel error must fire before any statement is issued.
type stubBuilderConn struct{}

func (stubBuilderConn) Select(context.Context, string, ...any) ([]map[string]any, error) {
	return nil, nil
}
func (stubBuilderConn) Insert(context.Context, string, ...any) (bool, error)  { return false, nil }
func (stubBuilderConn) Update(context.Context, string, ...any) (int64, error) { return 0, nil }
func (stubBuilderConn) Delete(context.Context, string, ...any) (int64, error) { return 0, nil }
func (stubBuilderConn) Statement(context.Context, string, ...any) (bool, error) {
	return false, nil
}
func (stubBuilderConn) AffectingStatement(context.Context, string, ...any) (int64, error) {
	return 0, nil
}
func (stubBuilderConn) Raw(string) dbcontract.Expression { return nil }
func (stubBuilderConn) GetTablePrefix() string           { return "" }

// TestBuilder_PropagatesClickHouseSentinels verifies that the
// `if err != nil { return ... }` branches in query.Builder.InsertGetId,
// .Upsert, and .InsertOrIgnore propagate the typed sentinel returned by
// the ClickHouse grammar instead of falling through to the connection.
func TestBuilder_PropagatesClickHouseSentinels(t *testing.T) {
	t.Parallel()

	g := grammars.NewClickHouseGrammar()
	conn := stubBuilderConn{}

	t.Run("InsertGetId", func(t *testing.T) {
		t.Parallel()

		b := query.NewBuilder(conn, g, nil).From("events")
		_, err := b.InsertGetId(context.Background(), map[string]any{"name": "x"})

		if !errors.Is(err, database.ErrInsertGetIdNotSupported) {
			t.Fatalf("expected ErrInsertGetIdNotSupported, got %v", err)
		}
	})

	t.Run("Upsert", func(t *testing.T) {
		t.Parallel()

		b := query.NewBuilder(conn, g, nil).From("events")
		_, err := b.Upsert(context.Background(),
			[]map[string]any{{"id": 1}},
			[]string{"id"},
			[]string{"id"})

		if !errors.Is(err, database.ErrUpsertNotSupported) {
			t.Fatalf("expected ErrUpsertNotSupported, got %v", err)
		}
	})

	t.Run("InsertOrIgnore", func(t *testing.T) {
		t.Parallel()

		b := query.NewBuilder(conn, g, nil).From("events")
		_, err := b.InsertOrIgnore(context.Background(), map[string]any{"id": 1})

		if !errors.Is(err, database.ErrInsertOrIgnoreNotSupported) {
			t.Fatalf("expected ErrInsertOrIgnoreNotSupported, got %v", err)
		}
	})
}
