import fs from 'node:fs'
import path from 'node:path'

const root = path.resolve(new URL('../../../', import.meta.url).pathname)
const docsRoot = path.join(root, 'services/docs')
const packagesRoot = path.join(root, 'packages')
const packageDocsRoot = path.join(docsRoot, 'packages')
const docsInventoryPath = path.join(root, 'services/compliance/docs-inventories/laravel-docs.txt')
const validLaravelSources = new Set(read(docsInventoryPath).split('\n').filter((line) => line && !line.startsWith('#')))

const laravelSources = {
  'ai/boost': ['boost.md#introduction', 'boost.md#installation', 'boost.md#mcp-server', 'boost.md#ai-guidelines', 'boost.md#agent-skills', 'boost.md#documentation-api', 'boost.md#extending-boost'],
  'ai/mcp': ['mcp.md#introduction', 'mcp.md#installation', 'mcp.md#creating-servers', 'mcp.md#web-server', 'mcp.md#creating-resources', 'mcp.md#prompts', 'mcp.md#testing-servers', 'mcp.md#tools'],
  'ai/sdk': ['ai-sdk.md#introduction', 'ai-sdk.md#installation', 'ai-sdk.md#agents', 'ai-sdk.md#images', 'ai-sdk.md#audio', 'ai-sdk.md#transcription', 'ai-sdk.md#embeddings', 'ai-sdk.md#reranking', 'ai-sdk.md#files', 'ai-sdk.md#vector-stores', 'ai-sdk.md#failover', 'ai-sdk.md#testing', 'ai-sdk.md#events'],
  auth: ['authentication.md#authentication', 'authentication.md#manually-authenticating-users', 'authentication.md#http-basic-authentication', 'authorization.md#authorization', 'passwords.md#resetting-passwords', 'verification.md#email-verification'],
  broadcasting: ['broadcasting.md#introduction', 'broadcasting.md#quickstart', 'broadcasting.md#client-side-installation'],
  bus: ['queues.md#creating-jobs', 'queues.md#job-batching', 'queues.md#chains-and-batches'],
  cache: ['cache.md#cache', 'cache.md#cache-usage', 'cache.md#cache-tags', 'cache.md#atomic-locks', 'cache.md#cache-failover', 'cache.md#events'],
  collection: ['collections.md#introduction', 'collections.md#method-listing', 'collections.md#higher-order-messages'],
  concurrency: ['concurrency.md#introduction', 'concurrency.md#running-concurrent-tasks'],
  conditionable: ['collections.md#higher-order-messages', 'helpers.md#other-utilities'],
  config: ['configuration.md#introduction', 'configuration.md#environment-configuration', 'configuration.md#accessing-configuration-values'],
  console: ['artisan.md#introduction', 'artisan.md#defining-input-expectations', 'scheduling.md#task-scheduling'],
  container: ['container.md#introduction', 'container.md#binding', 'container.md#resolving', 'container.md#method-invocation-and-injection'],
  contracts: ['contracts.md#introduction', 'providers.md#service-providers'],
  cookie: ['requests.md#cookies', 'responses.md#attaching-cookies-to-responses', 'encryption.md#introduction'],
  database: ['database.md#database', 'database.md#introduction', 'queries.md#queries', 'queries.md#running-database-queries', 'migrations.md#migrations', 'migrations.md#introduction', 'eloquent.md#eloquent', 'eloquent.md#introduction', 'seeding.md#seeding', 'seeding.md#introduction'],
  echo: ['broadcasting.md#client-side-installation', 'reverb.md#introduction'],
  encryption: ['encryption.md#introduction', 'encryption.md#using-the-encrypter'],
  envoy: ['envoy.md#introduction', 'envoy.md#writing-tasks', 'envoy.md#notifications'],
  events: ['events.md#introduction', 'events.md#defining-events', 'events.md#defining-listeners', 'events.md#queued-event-listeners'],
  facades: ['facades.md#introduction', 'facades.md#facade-class-reference'],
  filesystem: ['filesystem.md#filesystem', 'filesystem.md#introduction', 'filesystem.md#configuration', 'filesystem.md#obtaining-disk-instances', 'filesystem.md#retrieving-files', 'filesystem.md#storing-files', 'filesystem.md#deleting-files', 'filesystem.md#directories'],
  fortify: ['fortify.md#introduction', 'fortify.md#authentication', 'fortify.md#registration', 'fortify.md#password-reset'],
  hashing: ['hashing.md#introduction', 'hashing.md#basic-usage', 'hashing.md#determining-if-a-password-needs-to-be-rehashed'],
  horizon: ['horizon.md#introduction', 'horizon.md#configuration', 'horizon.md#monitoring', 'horizon.md#balancing-strategies'],
  httpx: ['requests.md#input', 'responses.md#creating-responses', 'http-client.md#introduction', 'http-tests.md#making-requests'],
  inception: ['fortify.md#introduction', 'sanctum.md#introduction', 'passport.md#introduction'],
  inertia: ['frontend.md#using-react-svelte-or-vue', 'responses.md#redirects', 'csrf.md#csrf-protection'],
  jetstream: ['fortify.md#introduction', 'sanctum.md#spa-authentication', 'sanctum.md#api-token-authentication'],
  jsonx: ['validation.md#working-with-validated-input', 'mcp.md#creating-tools'],
  log: ['logging.md#introduction', 'logging.md#building-log-stacks', 'logging.md#writing-log-messages'],
  lottery: ['helpers.md#other-utilities'],
  mailx: ['mail.md#introduction', 'mail.md#generating-mailables', 'mail.md#attachments', 'mail.md#events'],
  money: ['validation.md#available-validation-rules', 'localization.md#introduction'],
  notifications: ['notifications.md#introduction', 'notifications.md#sending-notifications', 'notifications.md#mail-notifications', 'notifications.md#notification-events'],
  pagination: ['pagination.md#introduction', 'pagination.md#cursor-pagination', 'pagination.md#customizing-the-pagination-view'],
  pail: ['pail.md#introduction', 'logging.md#writing-log-messages'],
  passport: ['passport.md#introduction', 'passport.md#password-grant', 'passport.md#personal-access-tokens', 'passport.md#protecting-routes'],
  pennant: ['pennant.md#introduction', 'pennant.md#defining-features', 'pennant.md#scope', 'pennant.md#adding-custom-pennant-drivers'],
  pipeline: ['helpers.md#other-utilities', 'middleware.md#middleware'],
  precognition: ['precognition.md#introduction', 'precognition.md#handling-file-uploads', 'precognition.md#managing-side-effects'],
  process: ['processes.md#introduction', 'processes.md#invoking-processes', 'processes.md#concurrent-processes', 'processes.md#testing'],
  prompts: ['prompts.md#introduction', 'prompts.md#forms', 'prompts.md#tables', 'prompts.md#progress-bars'],
  queue: ['queues.md#introduction', 'queues.md#creating-jobs', 'queues.md#job-batching', 'queues.md#running-the-queue-worker', 'queues.md#failed-jobs'],
  redis: ['redis.md#introduction', 'redis.md#interacting-with-redis', 'redis.md#pipelining-commands', 'redis.md#pubsub'],
  reverb: ['reverb.md#introduction', 'reverb.md#application-credentials', 'reverb.md#running-reverb-in-production', 'broadcasting.md#client-side-installation'],
  routing: ['routing.md#routing', 'routing.md#basic-routing', 'routing.md#route-parameters', 'routing.md#named-routes', 'routing.md#middleware', 'routing.md#route-groups', 'routing.md#route-model-binding', 'routing.md#fallback-routes', 'routing.md#rate-limiting'],
  scout: ['scout.md#introduction', 'scout.md#searchable-models', 'scout.md#querying', 'search.md#search'],
  seo: ['responses.md#creating-responses', 'frontend.md#frontend'],
  session: ['session.md#session', 'session.md#introduction', 'session.md#using-the-session', 'session.md#flash-data'],
  socialite: ['socialite.md#introduction', 'socialite.md#routing', 'socialite.md#authentication-and-storage'],
  spark: ['spark.md#introduction', 'cashier.md#subscriptions'],
  str: ['strings.md#introduction', 'helpers.md#available-methods'],
  support: ['helpers.md#introduction', 'collections.md#introduction', 'contracts.md#introduction'],
  telescope: ['telescope.md#introduction', 'telescope.md#watchers', 'telescope.md#dashboard-authorization'],
  translation: ['localization.md#introduction', 'localization.md#retrieving-translation-strings', 'localization.md#pluralization'],
  validation: ['validation.md#validation', 'validation.md#introduction', 'validation.md#quickstart', 'validation.md#manually-creating-validators', 'validation.md#available-validation-rules'],
  wayfinder: ['urls.md#urls-for-controller-actions', 'routing.md#named-routes', 'frontend.md#using-react-svelte-or-vue'],
}

const packageTitles = {
  httpx: 'HTTP',
  jsonx: 'JSON Schema',
  mailx: 'Mail',
  str: 'Strings',
}

function walk(dir, predicate = () => true) {
  if (!fs.existsSync(dir)) return []
  const entries = fs.readdirSync(dir, { withFileTypes: true })
  return entries.flatMap((entry) => {
    const file = path.join(dir, entry.name)
    if (entry.isDirectory()) return walk(file, predicate)
    return predicate(file) ? [file] : []
  })
}

function packageSourcePath(slug) {
  return path.join(packagesRoot, slug)
}

function docsPath(slug) {
  return path.join(packageDocsRoot, `${slug}.md`)
}

function read(file) {
  return fs.existsSync(file) ? fs.readFileSync(file, 'utf8') : ''
}

function titleFor(slug) {
  return packageTitles[slug] ?? slug.split('/').at(-1).replace(/-/g, ' ')
}

function packageCommentFrom(file) {
  const body = read(file)
  const beforePackage = body.split(/\npackage\s+/)[0] ?? ''
  const lines = beforePackage.split('\n').map((line) => line.replace(/^\/\/\s?/, '').trimEnd())
  while (lines.length > 0 && lines[0].trim() === '') lines.shift()
  while (lines.length > 0 && lines.at(-1).trim() === '') lines.pop()
  return lines.join('\n')
}

function firstParagraph(text, fallback) {
  const paragraph = text.split(/\n\s*\n/).find((part) => part.trim().length > 0)
  return paragraph ? sanitizeMarkdownText(paragraph.replace(/\s+/g, ' ').trim()) : fallback
}

function sanitizeMarkdownText(text) {
  return text
    .replaceAll('{{', '{&#123;')
    .replaceAll('}}', '&#125;}')
    .replace(/<([^>\n]+)>/g, '&lt;$1&gt;')
}

function exportedSymbols(files) {
  const types = new Set()
  const funcs = new Set()
  const vars = new Set()
  const consts = new Set()

  for (const file of files) {
    const body = read(file)
    for (const match of body.matchAll(/^type\s+([A-Z][A-Za-z0-9_]*)/gm)) types.add(match[1])
    for (const match of body.matchAll(/^func\s+(?:\([^)]*\)\s*)?([A-Z][A-Za-z0-9_]*)\s*\(/gm)) funcs.add(match[1])
    for (const match of body.matchAll(/^var\s+([A-Z][A-Za-z0-9_]*)/gm)) vars.add(match[1])
    for (const match of body.matchAll(/^const\s+([A-Z][A-Za-z0-9_]*)/gm)) consts.add(match[1])
    for (const block of body.matchAll(/^const\s+\(([\s\S]*?)^\)/gm)) {
      for (const match of block[1].matchAll(/^\s*([A-Z][A-Za-z0-9_]*)\b/gm)) consts.add(match[1])
    }
    for (const block of body.matchAll(/^var\s+\(([\s\S]*?)^\)/gm)) {
      for (const match of block[1].matchAll(/^\s*([A-Z][A-Za-z0-9_]*)\b/gm)) vars.add(match[1])
    }
  }

  return {
    types: [...types].sort(),
    funcs: [...funcs].sort(),
    vars: [...vars].sort(),
    consts: [...consts].sort(),
  }
}

function testInventory(sourceRoot) {
  return walk(sourceRoot, (file) => /(_test\.go)$/.test(file))
    .filter((file) => /inventory|laravel|compliance|parity/i.test(path.basename(file)))
    .map((file) => path.relative(root, file))
    .sort()
}

function subpackageDocs(sourceRoot) {
  const docDirs = walk(sourceRoot, (file) => path.basename(file) === 'doc.go').map((file) => path.dirname(file))
  const goDirs = walk(sourceRoot, (file) => file.endsWith('.go') && !file.endsWith('_test.go')).map((file) => path.dirname(file))
  return [...new Set([...docDirs, ...goDirs])]
    .filter((dir) => !path.relative(sourceRoot, dir).split(path.sep).includes('tests'))
    .map((dir) => {
      const docFile = path.join(dir, 'doc.go')
      return {
        rel: path.relative(sourceRoot, dir) || '.',
        comment: fs.existsSync(docFile) ? packageCommentFrom(docFile) : '',
      }
    })
    .sort((a, b) => a.rel.localeCompare(b.rel))
}

function detectCapabilities(sourceRoot, symbols) {
  const names = walk(sourceRoot, (file) => file.endsWith('.go') && !file.endsWith('_test.go'))
    .map((file) => path.basename(file).toLowerCase())
  const joined = `${names.join(' ')} ${symbols.types.join(' ')} ${symbols.funcs.join(' ')}`.toLowerCase()
  const capabilities = []

  const checks = [
    ['Drivers and managers', /\b(manager|driver|connector|store|provider|repository)\b/],
    ['HTTP middleware or handlers', /\b(middleware|handler|request|response|router)\b/],
    ['Events and listeners', /\b(event|listener|dispatcher|emitter)\b/],
    ['Queue, async, or background work', /\b(queue|worker|job|batch|async|deferred)\b/],
    ['Database-backed persistence', /\b(database|migration|query|schema|eloquent|sql|sqlite|mysql|postgres)\b/],
    ['Redis or distributed coordination', /\b(redis|lock|limiter|pubsub|cluster|sentinel)\b/],
    ['Testing fakes or null implementations', /\b(fake|null|array|testing|assert)\b/],
    ['Security-sensitive behavior', /\b(auth|token|password|encrypt|hash|csrf|sign|signature|oauth)\b/],
    ['Serialization or transport formats', /\b(json|xml|schema|marshal|unmarshal|encode|decode)\b/],
  ]

  for (const [label, pattern] of checks) {
    if (pattern.test(joined)) capabilities.push(label)
  }

  return capabilities
}

function existingLaravelComments(existing, slug) {
  const existingSources = (existing.match(/^<!-- laravel-docs:[\s\S]*?-->$/gm) ?? [])
    .map((comment) => comment.replace(/^<!--\s*laravel-docs:\s*/, '').replace(/\s*-->$/, ''))
  const preferred = laravelSources[slug] ?? existingSources
  const sources = preferred.filter((source) => validLaravelSources.has(source))

  if (sources.length === 0) {
    const fallback = `${slug}.md#introduction`
    if (validLaravelSources.has(fallback)) sources.push(fallback)
  }

  return sources.map((source) => `<!-- laravel-docs: ${source} -->`).join('\n')
}

function codeList(items, limit = 20) {
  const shown = items.slice(0, limit).map((item) => `\`${item}\``).join(', ')
  const rest = items.length > limit ? `, and ${items.length - limit} more` : ''
  return shown ? `${shown}${rest}` : 'None exported from this package root.'
}

function tableRows(items, empty) {
  if (items.length === 0) return `| ${empty} | |\n`
  return items.map((item) => `| \`${item}\` | Source-backed public surface. See the Go package for exact signature and behavior. |`).join('\n') + '\n'
}

function monorepoTestCommand(slug) {
  return `GOWORK=./storage/.cache/go.work go test -count=1 ./packages/${slug}/...`
}

function coreConceptsFor(title, packageName) {
  return `The ${title} reference is organized around the exported Go surface for package \`${packageName}\`. Start from the source coverage and public surface tables to identify the constructors, managers, interfaces, sentinel errors, and helper functions available to callers. Use the package tests as executable wiring examples for collaborators, default behavior, and Laravel parity expectations.`
}

function packageMapTable(subdocs, sourceRoot) {
  if (subdocs.length === 0) {
    return '| Package | Purpose |\n| --- | --- |\n| Root package | Primary public API for this module. |\n'
  }
  return [
    '| Package | Purpose |',
    '| --- | --- |',
    ...subdocs.map(({ rel, comment }) => {
      const name = rel === '.' ? path.basename(sourceRoot) : rel
      const purpose = firstParagraph(comment, `Public ${name} API surface for this module.`)
      return `| \`${name}\` | ${purpose.replace(/\|/g, '\\|')} |`
    }),
  ].join('\n') + '\n'
}

function guideFor(slug) {
  const sourceRoot = packageSourcePath(slug)
  const existing = read(docsPath(slug))
  const comments = existingLaravelComments(existing, slug)
  const goFiles = walk(sourceRoot, (file) => {
    const relParts = path.relative(sourceRoot, file).split(path.sep)
    return file.endsWith('.go') && !file.endsWith('_test.go') && !relParts.includes('tests')
  })
  const symbols = exportedSymbols(goFiles)
  const subdocs = subpackageDocs(sourceRoot)
  const rootDoc = subdocs.find((entry) => entry.rel === '.')?.comment ?? ''
  const summary = rootDoc || firstParagraph(rootDoc, `The ${titleFor(slug)} package provides Bedrock's Go implementation for this Laravel-aligned surface.`)
  const inventories = testInventory(sourceRoot)
  const capabilities = detectCapabilities(sourceRoot, symbols)
  const modulePath = `github.com/bedrock/packages/${slug}`
  const packageName = slug.split('/').at(-1)
  const title = titleFor(slug)
  const coreConcepts = coreConceptsFor(title, packageName)

  return `# ${title}

${comments}

${summary}

<div class="docs-callout docs-callout-laravel">
  <strong>Laravel baseline.</strong>
  This page follows the Laravel 13.x documentation structure for the matching feature area, then rewrites the examples and edge cases for Bedrock's Go packages.
</div>

<div class="docs-callout docs-callout-go">
  <strong>Go adaptation.</strong>
  Bedrock replaces Laravel facades, service container magic, PHP traits, and Artisan commands with explicit Go constructors, interfaces, structs, context propagation, and ordinary package tests.
</div>

## Installation

Install this module directly in applications that consume packages independently:

\`\`\`bash
go get ${modulePath}@latest
\`\`\`

When working inside this monorepo, use the repository workspace:

\`\`\`bash
${monorepoTestCommand(slug)}
\`\`\`

## Source Coverage

${packageMapTable(subdocs, sourceRoot)}
## Core Concepts

${coreConcepts}

### Public Surface

| Surface | Exported API |
| --- | --- |
| Types | ${codeList(symbols.types)} |
| Constructors and functions | ${codeList(symbols.funcs)} |
| Variables | ${codeList(symbols.vars)} |
| Constants | ${codeList(symbols.consts)} |

### Capability Matrix

| Capability | Documentation note |
| --- | --- |
${capabilities.length > 0 ? capabilities.map((capability) => `| ${capability} | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |`).join('\n') : '| Core package API | The root constructors and exported types are the primary integration surface. |'}

## Usage

Start with the package constructor or manager type when one is exported. Bedrock keeps dependencies explicit, so callers should pass repositories, stores, handlers, dispatchers, clocks, or clients directly instead of relying on global framework state.

\`\`\`go
package main

import (
    _ "${modulePath}"
)

func main() {
    // Import the package you use, then wire the exported constructors,
    // managers, stores, handlers, or helpers required by your application.
}
\`\`\`

Use package tests as executable examples when the exact constructor requires collaborators. The tests under \`packages/${slug}\` cover the supported creation paths, default values, and Laravel parity behavior.

## Configuration

Laravel documents many features through configuration files. Bedrock documents the equivalent behavior through Go options and constructor arguments:

| Laravel shape | Bedrock shape |
| --- | --- |
| Config file keys | Typed config structs, options, or constructor parameters |
| Facade defaults | Explicit manager/default-driver setup |
| Service providers | Go service-provider structs or direct application wiring |
| Runtime helpers | Package functions and interfaces |

Prefer narrow interfaces at package boundaries. When a package exposes a manager, register drivers or providers at startup, set the default once, and resolve named instances per request or job.

## Advanced Features

The package reference should be read through these Laravel parity lenses:

| Area | Documentation coverage |
| --- | --- |
| Drivers/providers | Available implementations, default selection, custom registration, and failure behavior |
| Events | Emitted structs, dispatcher hooks, listener timing, transaction or queue interaction |
| Errors | Exported sentinel errors, wrapping, and \`errors.Is\` compatibility |
| Context | Which operations accept \`context.Context\` and how cancellation/deadlines propagate |
| Testing | Fakes, null implementations, assertion helpers, and deterministic clocks/stores |

## Edge Cases

- Do not translate PHP-only behavior literally. If Laravel depends on PHP traits, request globals, Blade, Artisan, or Eloquent magic, document the Bedrock Go equivalent instead.
- Preserve error identity when the package exports sentinel errors; callers should be able to use \`errors.Is\` where the package promises it.
- Treat driver compatibility as observable behavior. Unsupported store/driver combinations should be documented as errors or explicit no-ops, never as silent omissions.
- For I/O paths, document cancellation and timeout behavior whenever the package accepts a \`context.Context\`.
- For test fakes, document whether assertions inspect recorded calls, stored payloads, emitted events, or rendered output.

## Testing

Run the package tests before changing examples:

\`\`\`bash
${monorepoTestCommand(slug)}
\`\`\`

${inventories.length > 0 ? `Laravel parity is tracked by these tests:\n\n${inventories.map((file) => `- \`${file}\``).join('\n')}` : 'No dedicated Laravel inventory test was detected for this package. Use the ordinary package tests and exported API as the documentation source of truth.'}

## API Reference

### Exported Types

| Type | Notes |
| --- | --- |
${tableRows(symbols.types, 'No exported types detected')}
### Exported Functions

| Function | Notes |
| --- | --- |
${tableRows(symbols.funcs, 'No exported functions detected')}
### Exported Errors, Variables, and Constants

| Name | Notes |
| --- | --- |
${tableRows([...symbols.vars, ...symbols.consts].sort(), 'No exported variables or constants detected')}
## Laravel Parity Notes

This page should stay aligned with the official Laravel 13.x documentation for the corresponding feature while keeping the Go API explicit. If Bedrock implements a Laravel feature, document the user-facing behavior, the Go entry points, supported drivers, emitted events, error behavior, and the tests that prove parity. If a Laravel feature is PHP-only, record the exclusion in \`services/compliance/docs-status.yml\` instead of inventing a Go API.
`
}

function main() {
  const topLevel = fs.readdirSync(packagesRoot, { withFileTypes: true })
    .filter((entry) => entry.isDirectory())
    .map((entry) => entry.name)
    .filter((name) => name !== 'ai')
    .sort()
  const ai = fs.readdirSync(path.join(packagesRoot, 'ai'), { withFileTypes: true })
    .filter((entry) => entry.isDirectory())
    .map((entry) => `ai/${entry.name}`)
    .sort()
  const slugs = [...ai, ...topLevel]

  for (const slug of slugs) {
    const file = docsPath(slug)
    fs.mkdirSync(path.dirname(file), { recursive: true })
    fs.writeFileSync(file, guideFor(slug))
  }
}

main()
