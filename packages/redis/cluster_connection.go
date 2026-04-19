package redis

import "context"

// ClusterAware is optionally implemented by a Client to expose
// cluster-wide fan-out operations. The go-redis ClusterClient adapter
// implements it; the in-memory fake does not (tests treat a fake cluster
// as a single node).
type ClusterAware interface {
	ForEachMaster(ctx context.Context, fn func(Client) error) error
}

// ClusterScan fans out a SCAN across master nodes if the underlying
// client is cluster-aware. Otherwise it behaves like Scan.
//
// Parity with PhpRedisClusterConnection::scan.
func (c *Connection) ClusterScan(ctx context.Context, cursor uint64, match string, count int64) (ScanResult, error) {
	if ca, ok := c.client.(ClusterAware); ok {
		var all []string
		err := ca.ForEachMaster(ctx, func(cl Client) error {
			shadow := NewConnection(c.name, cl)
			// walk each master from cursor 0 until complete
			var cur uint64

			for {
				res, err := shadow.Scan(ctx, cur, match, count)

				if err != nil {
					return err
				}

				all = append(all, res.Values...)

				if res.Cursor == 0 {
					return nil
				}

				cur = res.Cursor
			}
		})

		return ScanResult{Cursor: 0, Values: all}, err
	}

	return c.Scan(ctx, cursor, match, count)
}

// Keys returns all keys matching the given pattern. When the underlying
// client is cluster-aware the lookup fans out across masters.
func (c *Connection) Keys(ctx context.Context, pattern string) ([]string, error) {
	if ca, ok := c.client.(ClusterAware); ok {
		var all []string
		err := ca.ForEachMaster(ctx, func(cl Client) error {
			v, err := cl.Do(ctx, "KEYS", pattern)

			if err != nil {
				return err
			}

			s, _ := toStringSlice(v)
			all = append(all, s...)

			return nil
		})

		return all, err
	}

	v, err := c.Command(ctx, "KEYS", pattern)

	if err != nil {
		return nil, err
	}

	return toStringSlice(v)
}

// ClusterFlushDB flushes every master when cluster-aware; otherwise
// delegates to FlushDB.
func (c *Connection) ClusterFlushDB(ctx context.Context) error {
	if ca, ok := c.client.(ClusterAware); ok {
		return ca.ForEachMaster(ctx, func(cl Client) error {
			_, err := cl.Do(ctx, "FLUSHDB")

			return err
		})
	}

	return c.FlushDB(ctx)
}
