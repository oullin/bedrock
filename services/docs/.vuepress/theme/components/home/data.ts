// Landing-page content. Ported from the design mock.
// `href` points at an existing docs route (three design IDs were remapped to
// real package pages: https→httpx, mail→mailx, forge→prompts).

export interface Package {
  id: string
  name: string
  tier: 'http' | 'data' | 'infra'
  icon: string
  tag: string
  blurb: string
  v: string
  downloads: string
  since: string
  href: string
}

export interface Tier {
  id: 'http' | 'data' | 'infra'
  label: string
  sub: string
}

export interface HeroTab {
  id: string
  label: string
  code: string
}

export const PACKAGES: Package[] = [
  // --- HTTP LAYER ---
  { id: 'container', name: 'container', tier: 'http', icon: 'IconContainer', tag: 'core',
    blurb: 'IoC service container with bindings, singletons, contextual resolvers and a full lifecycle.',
    v: '1.4.0', downloads: '812k', since: '2023', href: '/packages/container' },
  { id: 'routing', name: 'routing', tier: 'http', icon: 'IconRouting', tag: 'core',
    blurb: 'Expressive HTTP router — path params, groups, resource routes, middleware pipelines.',
    v: '1.7.2', downloads: '1.1M', since: '2023', href: '/packages/routing' },
  { id: 'httpx', name: 'httpx', tier: 'http', icon: 'IconHttp', tag: 'io',
    blurb: 'Fluent request/response primitives. Content negotiation, cookies, signed URLs, streaming.',
    v: '1.3.0', downloads: '940k', since: '2023', href: '/packages/httpx' },
  { id: 'validation', name: 'validation', tier: 'http', icon: 'IconValidation', tag: 'io',
    blurb: 'Tag-based and programmatic validators. Composable rules, friendly errors, localizable.',
    v: '0.9.5', downloads: '420k', since: '2024', href: '/packages/validation' },
  { id: 'auth', name: 'auth', tier: 'http', icon: 'IconAuth', tag: 'security',
    blurb: 'Sessions, tokens, guards and gates. Pluggable drivers for JWT, OAuth and passkeys.',
    v: '1.1.0', downloads: '560k', since: '2023', href: '/packages/auth' },
  { id: 'session', name: 'session', tier: 'http', icon: 'IconShield', tag: 'security',
    blurb: 'HTTP sessions with flash data, CSRF tokens and pluggable stores — cookie, file, database, cache.',
    v: '1.0.2', downloads: '390k', since: '2023', href: '/packages/session' },

  // --- DATA & MESSAGING ---
  { id: 'cache', name: 'cache', tier: 'data', icon: 'IconCache', tag: 'data',
    blurb: 'Uniform cache API with tagging, atomic locks, stampede protection and drivers for Redis, memcached, file.',
    v: '1.2.1', downloads: '780k', since: '2023', href: '/packages/cache' },
  { id: 'redis', name: 'redis', tier: 'data', icon: 'IconRedis', tag: 'data',
    blurb: 'Production Redis client with pipelines, transactions, Lua scripts, cluster and Sentinel support.',
    v: '1.0.4', downloads: '610k', since: '2023', href: '/packages/redis' },
  { id: 'events', name: 'events', tier: 'data', icon: 'IconEvents', tag: 'messaging',
    blurb: 'Pub/sub with wildcard listeners, subscribers, queued delivery and transactional outbox.',
    v: '1.1.3', downloads: '330k', since: '2023', href: '/packages/events' },
  { id: 'queue', name: 'queue', tier: 'data', icon: 'IconQueue', tag: 'messaging',
    blurb: 'Durable background jobs with retries, backoff, batching and workers across Redis, SQS and NATS.',
    v: '1.2.0', downloads: '450k', since: '2023', href: '/packages/queue' },
  { id: 'bus', name: 'bus', tier: 'data', icon: 'IconBus', tag: 'messaging',
    blurb: 'Command and query bus with middleware, pipelines, and sync/async dispatch semantics.',
    v: '0.8.1', downloads: '190k', since: '2024', href: '/packages/bus' },
  { id: 'mailx', name: 'mailx', tier: 'data', icon: 'IconMail', tag: 'io',
    blurb: 'Template-based email with attachments, markdown, mailers and local preview tooling.',
    v: '1.0.2', downloads: '240k', since: '2023', href: '/packages/mailx' },

  // --- INFRA & UTIL ---
  { id: 'encryption', name: 'encryption', tier: 'infra', icon: 'IconEncryption', tag: 'security',
    blurb: 'AES-GCM, ChaCha20 and key rotation. MAC-verified payloads with zero-copy decryption.',
    v: '1.0.1', downloads: '210k', since: '2024', href: '/packages/encryption' },
  { id: 'hashing', name: 'hashing', tier: 'infra', icon: 'IconHash', tag: 'security',
    blurb: 'Password hashing with bcrypt, argon2id and scrypt. Hash-then-verify with timing-safe equality.',
    v: '1.0.0', downloads: '330k', since: '2023', href: '/packages/hashing' },
  { id: 'filesystem', name: 'filesystem', tier: 'infra', icon: 'IconFile', tag: 'io',
    blurb: 'One filesystem interface across local, S3, GCS and memory. Atomic writes, signed URLs, glob.',
    v: '1.1.0', downloads: '380k', since: '2023', href: '/packages/filesystem' },
  { id: 'pagination', name: 'pagination', tier: 'infra', icon: 'IconPagination', tag: 'data',
    blurb: 'Offset and cursor paginators for SQL, REST and GraphQL. JSON:API compatible.',
    v: '0.7.0', downloads: '120k', since: '2024', href: '/packages/pagination' },
  { id: 'translation', name: 'translation', tier: 'infra', icon: 'IconTranslation', tag: 'util',
    blurb: 'ICU-MessageFormat and YAML translations with pluralization, gender and in-memory hot reload.',
    v: '0.9.0', downloads: '95k', since: '2024', href: '/packages/translation' },
  { id: 'support', name: 'support', tier: 'infra', icon: 'IconSupport', tag: 'util',
    blurb: 'Tiny utilities — Optional[T], Map/Slice helpers, string pipelines, time and Result[T, E].',
    v: '1.3.0', downloads: '1.4M', since: '2023', href: '/packages/support' },
  { id: 'log', name: 'log', tier: 'infra', icon: 'IconLog', tag: 'util',
    blurb: 'Structured logging with stacked handlers, channels, context and sampling. slog-compatible.',
    v: '1.2.0', downloads: '720k', since: '2023', href: '/packages/log' },
  { id: 'inception', name: 'inception', tier: 'infra', icon: 'IconInception', tag: 'core',
    blurb: 'Unified application scaffolding — bootstrap, service providers, lifecycle hooks, env and config.',
    v: '1.0.3', downloads: '510k', since: '2023', href: '/packages/inception' },
  { id: 'prompts', name: 'prompts', tier: 'infra', icon: 'IconBilling', tag: 'dx',
    blurb: 'Artisan-style CLI — codegen, migrations, interactive prompts, batteries included.',
    v: '1.1.1', downloads: '640k', since: '2023', href: '/packages/prompts' },
]

export const TIERS: Tier[] = [
  { id: 'http', label: 'HTTP', sub: 'request, route, respond' },
  { id: 'data', label: 'Data & messaging', sub: 'cache, queues, events' },
  { id: 'infra', label: 'Infrastructure', sub: 'security, storage, utils' },
]

// HERO_TABS — code samples use the real public API of each package, taken
// from the source in /packages/{container,routing,queue}.
export const HERO_TABS: HeroTab[] = [
  {
    id: 'container',
    label: 'container.go',
    code: `package main

import "github.com/bedrock/packages/container"

type Mailer struct{ Host string }

func main() {
  c := container.New()

  c.Bind("mailer", func(_ *container.Container) (any, error) {
    return &Mailer{Host: "smtp.example.com"}, nil
  }, false)

  c.Singleton("config", func(_ *container.Container) (any, error) {
    return loadConfig(), nil
  })

  mailer, _ := c.Make("mailer")
  _ = mailer
}`,
  },
  {
    id: 'routing',
    label: 'routes.go',
    code: `package main

import "github.com/bedrock/packages/routing"

func main() {
  r := routing.NewRouter(nil, nil)

  r.Get("/posts", listPosts)
  r.Get("/posts/{id}", showPost)

  r.Group(map[string]any{"prefix": "api", "as": "api."},
    func(g *routing.Router) {
      g.Post("/posts",       createPost).Name("posts.store")
      g.Patch("/posts/{id}", updatePost).Name("posts.update")
      g.Delete("/posts/{id}", deletePost).Name("posts.destroy")
    },
  )
}`,
  },
  {
    id: 'queue',
    label: 'jobs.go',
    code: `package jobs

import (
  "context"

  "github.com/bedrock/packages/queue"
)

type SendWelcome struct {
  queue.InteractsWithQueue
  UserID string
}

func (j *SendWelcome) Handle(ctx context.Context, job queue.Job) error {
  if err := sendWelcomeEmail(ctx, j.UserID); err != nil {
    return j.Release(30 * time.Second)
  }
  return j.Delete()
}`,
  },
]
