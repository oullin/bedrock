package database

import dbcontract "github.com/bedrock/packages/contracts/database"

// Expr is a raw SQL expression that will not be parameterised or quoted by the
// query builder. It implements the contracts/database.Expression interface.
type Expr struct {
	value string
}

var _ dbcontract.Expression = (*Expr)(nil)

// NewExpr creates a raw SQL expression.
func NewExpr(value string) *Expr {
	return &Expr{value: value}
}

// GetValue returns the raw SQL string.
func (e *Expr) GetValue() string {
	return e.value
}

// String implements fmt.Stringer.
func (e *Expr) String() string {
	return e.value
}
