package redis

import (
	"context"
	"strconv"

	goredis "github.com/redis/go-redis/v9"
)

// goRedisClient adapts a go-redis UniversalClient to the Client interface.
type goRedisClient struct {
	uc goredis.UniversalClient
}

// NewGoRedisClient wraps a go-redis client so it can be used with this
// package. Both *redis.Client and *redis.ClusterClient are
// UniversalClients, so this adapter works for single-node, cluster, and
// failover setups.
func NewGoRedisClient(uc goredis.UniversalClient) Client {
	return &goRedisClient{uc: uc}
}

func (c *goRedisClient) Do(ctx context.Context, args ...any) (any, error) {
	res, err := c.uc.Do(ctx, args...).Result()
	if err == goredis.Nil {
		return nil, ErrNil
	}
	return res, err
}

func (c *goRedisClient) Pipeline() Pipeliner {
	return &goRedisPipeline{p: c.uc.Pipeline()}
}

func (c *goRedisClient) TxPipeline() Pipeliner {
	return &goRedisPipeline{p: c.uc.TxPipeline()}
}

func (c *goRedisClient) Subscribe(ctx context.Context, channels ...string) Subscription {
	ps := c.uc.Subscribe(ctx, channels...)
	return newGoRedisSubscription(ps, false)
}

func (c *goRedisClient) PSubscribe(ctx context.Context, patterns ...string) Subscription {
	ps := c.uc.PSubscribe(ctx, patterns...)
	return newGoRedisSubscription(ps, true)
}

func (c *goRedisClient) Close() error { return c.uc.Close() }

// ForEachMaster implements ClusterAware when the underlying client is a
// cluster. Non-cluster clients fall back to invoking fn on themselves.
func (c *goRedisClient) ForEachMaster(ctx context.Context, fn func(Client) error) error {
	if cc, ok := c.uc.(*goredis.ClusterClient); ok {
		return cc.ForEachMaster(ctx, func(ctx context.Context, master *goredis.Client) error {
			return fn(NewGoRedisClient(master))
		})
	}
	return fn(c)
}

// -- pipeline -----------------------------------------------------------

type goRedisPipeline struct {
	p     goredis.Pipeliner
	queue []*goRedisCmd
}

func (p *goRedisPipeline) Do(ctx context.Context, args ...any) Cmder {
	cmd := p.p.Do(ctx, args...)
	wrapped := &goRedisCmd{cmd: cmd, args: args}
	p.queue = append(p.queue, wrapped)
	return wrapped
}

func (p *goRedisPipeline) Exec(ctx context.Context) ([]Cmder, error) {
	_, err := p.p.Exec(ctx)
	// Even on partial error, populate results so callers can inspect.
	out := make([]Cmder, len(p.queue))
	for i, c := range p.queue {
		out[i] = c
	}
	// go-redis returns an error when any queued command errors; surface it
	// but keep individual results available via Cmder.Err / Result.
	if err == goredis.Nil {
		err = nil
	}
	return out, err
}

func (p *goRedisPipeline) Discard() { p.p.Discard(); p.queue = nil }
func (p *goRedisPipeline) Len() int { return len(p.queue) }

type goRedisCmd struct {
	cmd  *goredis.Cmd
	args []any
}

func (c *goRedisCmd) Name() string { return c.cmd.Name() }
func (c *goRedisCmd) Args() []any  { return c.args }
func (c *goRedisCmd) Result() (any, error) {
	v, err := c.cmd.Result()
	if err == goredis.Nil {
		return nil, ErrNil
	}
	return v, err
}
func (c *goRedisCmd) Err() error {
	if err := c.cmd.Err(); err != nil {
		if err == goredis.Nil {
			return ErrNil
		}
		return err
	}
	return nil
}

// -- subscription -------------------------------------------------------

type goRedisSubscription struct {
	ps       *goredis.PubSub
	ch       chan Message
	done     chan struct{}
	isPatt   bool
}

func newGoRedisSubscription(ps *goredis.PubSub, pattern bool) *goRedisSubscription {
	s := &goRedisSubscription{
		ps:     ps,
		ch:     make(chan Message, 16),
		done:   make(chan struct{}),
		isPatt: pattern,
	}
	go s.pump()
	return s
}

func (s *goRedisSubscription) pump() {
	defer close(s.ch)
	src := s.ps.Channel()
	for {
		select {
		case <-s.done:
			return
		case m, ok := <-src:
			if !ok {
				return
			}
			msg := Message{Channel: m.Channel, Payload: m.Payload}
			if s.isPatt {
				msg.Pattern = m.Pattern
			}
			select {
			case s.ch <- msg:
			case <-s.done:
				return
			}
		}
	}
}

func (s *goRedisSubscription) Channel() <-chan Message { return s.ch }

func (s *goRedisSubscription) Close() error {
	select {
	case <-s.done:
	default:
		close(s.done)
	}
	return s.ps.Close()
}

// DialSingle opens a single-node go-redis client from a ConnectionConfig.
func DialSingle(cfg ConnectionConfig) (Client, error) {
	opts := &goredis.Options{
		Addr:         addr(cfg),
		Username:     cfg.Username,
		Password:     cfg.Password,
		DB:           cfg.Database,
		DialTimeout:  cfg.Timeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		MaxRetries:   cfg.MaxRetries,
		MinIdleConns: cfg.MinIdleConns,
		PoolSize:     cfg.PoolSize,
	}
	return NewGoRedisClient(goredis.NewClient(opts)), nil
}

// DialCluster opens a cluster client.
func DialCluster(cfg ConnectionConfig) (Client, error) {
	if cfg.Cluster == nil {
		return nil, ErrDriverNotFound
	}
	opts := &goredis.ClusterOptions{
		Addrs:          cfg.Cluster.Addrs,
		Username:       cfg.Cluster.Username,
		Password:       cfg.Cluster.Password,
		DialTimeout:    cfg.Timeout,
		ReadTimeout:    cfg.ReadTimeout,
		WriteTimeout:   cfg.WriteTimeout,
		MaxRetries:     cfg.MaxRetries,
		PoolSize:       cfg.PoolSize,
		ReadOnly:       cfg.Cluster.ReadOnly,
		RouteByLatency: cfg.Cluster.RouteByLatency,
		RouteRandomly:  cfg.Cluster.RouteRandomly,
	}
	return NewGoRedisClient(goredis.NewClusterClient(opts)), nil
}

// DialSentinel opens a failover (sentinel) client.
func DialSentinel(cfg ConnectionConfig) (Client, error) {
	if cfg.Sentinel == nil {
		return nil, ErrDriverNotFound
	}
	opts := &goredis.FailoverOptions{
		MasterName:       cfg.Sentinel.MasterName,
		SentinelAddrs:    cfg.Sentinel.SentinelAddrs,
		SentinelPassword: cfg.Sentinel.SentinelPassword,
		Username:         cfg.Sentinel.Username,
		Password:         cfg.Sentinel.Password,
		DB:               cfg.Sentinel.Database,
		DialTimeout:      cfg.Timeout,
		ReadTimeout:      cfg.ReadTimeout,
		WriteTimeout:     cfg.WriteTimeout,
	}
	return NewGoRedisClient(goredis.NewFailoverClient(opts)), nil
}

func addr(cfg ConnectionConfig) string {
	if cfg.URL != "" {
		return cfg.URL
	}
	host := cfg.Host
	if host == "" {
		host = "127.0.0.1"
	}
	port := cfg.Port
	if port == 0 {
		port = 6379
	}
	return formatAddr(host, port)
}

func formatAddr(host string, port int) string {
	return host + ":" + strconv.Itoa(port)
}
