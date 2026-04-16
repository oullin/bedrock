package seeding

import (
	"context"

	dbcontract "github.com/bedrock/packages/contracts/database"
)

// Seeder populates database tables with data.
type Seeder interface {
	// Run executes the seeder.
	Run(ctx context.Context, conn dbcontract.Connection) error
}

// FuncSeeder implements Seeder using a function value.
type FuncSeeder struct {
	RunFunc func(ctx context.Context, conn dbcontract.Connection) error
}

// Run executes the seeder function.
func (s *FuncSeeder) Run(ctx context.Context, conn dbcontract.Connection) error {
	if s.RunFunc == nil {
		return nil
	}
	return s.RunFunc(ctx, conn)
}

// Runner executes a list of seeders in order.
type Runner struct {
	resolver   dbcontract.ConnectionResolver
	connection string
}

// NewRunner creates a new seeder Runner.
func NewRunner(resolver dbcontract.ConnectionResolver, connection string) *Runner {
	return &Runner{resolver: resolver, connection: connection}
}

// Run executes all given seeders.
func (r *Runner) Run(ctx context.Context, seeders ...Seeder) error {
	conn, err := r.resolver.Connection(ctx, r.connection)
	if err != nil {
		return err
	}

	for _, seeder := range seeders {
		if err := seeder.Run(ctx, conn); err != nil {
			return err
		}
	}

	return nil
}
