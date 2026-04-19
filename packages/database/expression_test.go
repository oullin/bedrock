package database_test

import (
	"testing"

	"github.com/bedrock/packages/database"
)

func TestExprGetValue(t *testing.T) {
	t.Parallel()

	expr := database.NewExpr("COUNT(*)")

	if expr.GetValue() != "COUNT(*)" {
		t.Fatalf("expected COUNT(*), got %s", expr.GetValue())
	}
}

func TestExprString(t *testing.T) {
	t.Parallel()

	expr := database.NewExpr("NOW()")

	if expr.String() != "NOW()" {
		t.Fatalf("expected NOW(), got %s", expr.String())
	}
}
