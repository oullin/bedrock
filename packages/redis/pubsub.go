package redis

import "context"

// Subscribe listens on the given channels, invoking fn for every incoming
// message. It returns when ctx is cancelled or Close is called on the
// returned cancel func via the context.
//
// Parity with PhpRedisConnection::subscribe.
func (c *Connection) Subscribe(ctx context.Context, channels []string, fn func(channel, payload string)) error {
	sub := c.client.Subscribe(ctx, channels...)
	defer sub.Close()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case m, ok := <-sub.Channel():
			if !ok {
				return nil
			}
			fn(m.Channel, m.Payload)
		}
	}
}

// PSubscribe listens on channel patterns, invoking fn for every match.
// Parity with PhpRedisConnection::psubscribe.
func (c *Connection) PSubscribe(ctx context.Context, patterns []string, fn func(pattern, channel, payload string)) error {
	sub := c.client.PSubscribe(ctx, patterns...)
	defer sub.Close()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case m, ok := <-sub.Channel():
			if !ok {
				return nil
			}
			fn(m.Pattern, m.Channel, m.Payload)
		}
	}
}
