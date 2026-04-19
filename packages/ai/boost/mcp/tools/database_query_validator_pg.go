//go:build cgo

package tools

import (
	"fmt"

	pgquery "github.com/pganalyze/pg_query_go/v6"
)

func validatePostgresReadOnlySQL(query string) error {
	tree, err := pgquery.Parse(query)

	if err != nil {
		return fmt.Errorf("invalid SQL: %w", err)
	}

	stmts := tree.GetStmts()

	if len(stmts) != 1 {
		return fmt.Errorf("read-only validation requires a single statement")
	}

	stmt := stmts[0]

	if stmt == nil || stmt.GetStmt() == nil {
		return fmt.Errorf("invalid SQL: missing statement")
	}

	return validatePostgresReadOnlyNode(stmt.GetStmt())
}

func validatePostgresReadOnlyNode(node *pgquery.Node) error {
	switch {
	case node == nil:
		return fmt.Errorf("query is not read-only")
	case node.GetSelectStmt() != nil:
		return validatePostgresSelectStmt(node.GetSelectStmt())
	case node.GetExplainStmt() != nil:
		explain := node.GetExplainStmt()

		if explain.GetQuery() == nil {
			return fmt.Errorf("query is not read-only")
		}

		return validatePostgresReadOnlyNode(explain.GetQuery())
	default:
		return fmt.Errorf("query is not read-only")
	}
}

func validatePostgresSelectStmt(stmt *pgquery.SelectStmt) error {
	if stmt == nil {
		return fmt.Errorf("query is not read-only")
	}

	if stmt.GetIntoClause() != nil {
		return fmt.Errorf("query is not read-only")
	}

	if len(stmt.GetLockingClause()) > 0 {
		return fmt.Errorf("query is not read-only")
	}

	if err := validatePostgresWithClause(stmt.GetWithClause()); err != nil {
		return err
	}

	if stmt.GetLarg() != nil {
		if err := validatePostgresSelectStmt(stmt.GetLarg()); err != nil {
			return err
		}
	}

	if stmt.GetRarg() != nil {
		if err := validatePostgresSelectStmt(stmt.GetRarg()); err != nil {
			return err
		}
	}

	return nil
}

func validatePostgresWithClause(with *pgquery.WithClause) error {
	if with == nil {
		return nil
	}

	for _, cteNode := range with.GetCtes() {
		if cteNode == nil || cteNode.GetCommonTableExpr() == nil {
			return fmt.Errorf("query is not read-only")
		}

		cte := cteNode.GetCommonTableExpr()

		if cte.GetCtequery() == nil {
			return fmt.Errorf("query is not read-only")
		}

		if err := validatePostgresReadOnlyNode(cte.GetCtequery()); err != nil {
			return err
		}
	}

	return nil
}
