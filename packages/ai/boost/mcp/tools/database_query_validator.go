package tools

import (
	"fmt"
	"strings"

	"vitess.io/vitess/go/vt/sqlparser"
)

func validateReadOnlySQL(driver, query string) error {
	normalized, ok := normalizeDatabaseQueryDriver(driver)

	if strings.TrimSpace(driver) == "" {
		return fmt.Errorf("driver is required for parser-verified read-only validation")
	}

	if !ok {
		return fmt.Errorf("unsupported driver %q for parser-verified read-only validation", driver)
	}

	switch normalized {
	case "postgres":
		return validatePostgresReadOnlySQL(query)
	case "mysql":
		return validateMySQLReadOnlySQL(query)
	default:
		return fmt.Errorf("unsupported driver %q for parser-verified read-only validation", driver)
	}
}

func normalizeDatabaseQueryDriver(driver string) (string, bool) {
	switch strings.ToLower(strings.TrimSpace(driver)) {
	case "postgres", "postgresql", "pgsql":
		return "postgres", true
	case "mysql", "mariadb":
		return "mysql", true
	default:
		return "", false
	}
}

func validateMySQLReadOnlySQL(query string) error {
	parser := sqlparser.NewTestParser()
	stmts, err := parser.ParseMultipleIgnoreEmpty(query)

	if err != nil {
		return fmt.Errorf("invalid SQL: %w", err)
	}

	if len(stmts) != 1 {
		return fmt.Errorf("read-only validation requires a single statement")
	}

	return validateMySQLReadOnlyStatement(stmts[0])
}

func validateMySQLReadOnlyStatement(stmt sqlparser.Statement) error {
	switch typed := stmt.(type) {
	case *sqlparser.Select:
		return validateMySQLReadOnlyTableStatement(typed)
	case *sqlparser.Union:
		return validateMySQLReadOnlyTableStatement(typed)
	case *sqlparser.ExplainStmt:
		if typed.Statement == nil {
			return fmt.Errorf("query is not read-only")
		}

		return validateMySQLReadOnlyStatement(typed.Statement)
	default:
		return fmt.Errorf("query is not read-only")
	}
}

func validateMySQLReadOnlyTableStatement(stmt sqlparser.TableStatement) error {
	switch typed := stmt.(type) {
	case *sqlparser.Select:
		if typed.Into != nil || typed.Lock != sqlparser.NoLock {
			return fmt.Errorf("query is not read-only")
		}

		if err := validateMySQLWithClause(typed.With); err != nil {
			return err
		}
	case *sqlparser.Union:
		if typed.Into != nil || typed.Lock != sqlparser.NoLock {
			return fmt.Errorf("query is not read-only")
		}

		if err := validateMySQLWithClause(typed.With); err != nil {
			return err
		}

		if err := validateMySQLReadOnlyTableStatement(typed.Left); err != nil {
			return err
		}

		if err := validateMySQLReadOnlyTableStatement(typed.Right); err != nil {
			return err
		}
	default:
		return fmt.Errorf("query is not read-only")
	}

	return nil
}

func validateMySQLWithClause(with *sqlparser.With) error {
	if with == nil {
		return nil
	}

	for _, cte := range with.CTEs {
		if cte == nil || cte.Subquery == nil {
			return fmt.Errorf("query is not read-only")
		}

		if err := validateMySQLReadOnlyTableStatement(cte.Subquery); err != nil {
			return err
		}
	}

	return nil
}
