package redis

import "context"

// Pipeline executes fn against a non-transactional pipeline and returns
// the list of result Cmders in order. Parity with
// PhpRedisConnection::pipeline.
func (c *Connection) Pipeline(ctx context.Context, fn func(Pipeliner) error) ([]Cmder, error) {
	p := c.client.Pipeline()

	if err := fn(p); err != nil {
		p.Discard()

		return nil, err
	}

	return p.Exec(ctx)
}

// Transaction executes fn inside a MULTI/EXEC transaction. Parity with
// PhpRedisConnection::transaction.
func (c *Connection) Transaction(ctx context.Context, fn func(Pipeliner) error) ([]Cmder, error) {
	p := c.client.TxPipeline()

	if err := fn(p); err != nil {
		p.Discard()

		return nil, err
	}

	return p.Exec(ctx)
}
