<script setup lang="ts">
import { usePageFrontmatter, useSiteData, RouteLink } from 'vuepress/client'
import { ChevronRight } from 'lucide-vue-next'
import Navbar from '../Navbar.vue'
import { Button } from '../ui/button'
import { Badge } from '../ui/badge'
import { Tabs, TabsList, TabsTrigger, TabsContent } from '../ui/tabs'
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '../ui/card'
import { cn } from '../../lib/utils'

interface Action {
  text: string
  link: string
  type?: 'primary' | 'secondary'
}

interface Feature {
  title: string
  details: string
}

interface HomeFrontmatter {
  heroText?: string
  tagline?: string
  actions?: Action[]
  features?: Feature[]
  footer?: string
}

const frontmatter = usePageFrontmatter<HomeFrontmatter>()
const site = useSiteData()

declare const __BEDROCK_THEME_OPTIONS__: string
const themeOptions = JSON.parse(__BEDROCK_THEME_OPTIONS__)

// Install tabs
const managers = [
  { label: 'auth',    cmd: 'go get github.com/gocanto/bedrock/packages/auth@latest' },
  { label: 'cache',   cmd: 'go get github.com/gocanto/bedrock/packages/cache@latest' },
  { label: 'routing', cmd: 'go get github.com/gocanto/bedrock/packages/routing@latest' },
  { label: 'queue',   cmd: 'go get github.com/gocanto/bedrock/packages/queue@latest' },
]

// Stats
const stats = [
  { value: '11',       label: 'Packages' },
  { value: 'Go 1.24+', label: 'Go version' },
  { value: 'MIT',      label: 'License' },
]
</script>

<template>
  <div class="flex min-h-full flex-col bg-[#0f0f12] antialiased text-white">
    <Navbar :navbar-items="themeOptions.navbar ?? []" @toggle-sidebar="() => {}" />

    <!-- ─── Hero ─────────────────────────────────────────────────────────── -->
    <section class="relative overflow-hidden border-b border-white/10">
      <!-- Background glow blobs -->
      <div class="pointer-events-none absolute inset-0 overflow-hidden" aria-hidden="true">
        <div
          class="absolute -top-40 left-1/2 h-[600px] w-[800px] -translate-x-1/2 rounded-full opacity-20"
          style="background: radial-gradient(ellipse at center, #36b49f55 0%, transparent 70%);"
        />
        <div
          class="absolute -bottom-20 right-0 h-[400px] w-[500px] rounded-full opacity-10"
          style="background: radial-gradient(ellipse at center, #DBFF7555 0%, transparent 70%);"
        />
      </div>

      <div class="relative mx-auto grid max-w-6xl grid-cols-1 gap-0 md:grid-cols-2 md:divide-x md:divide-white/10">

        <!-- Left: text -->
        <div class="flex flex-col justify-center gap-8 p-10 md:py-20 md:pr-16">
          <!-- Attribution badge -->
          <Badge
            variant="outline"
            as-child
            class="w-fit font-mono text-xs uppercase tracking-widest text-white/50 border-white/20 bg-transparent hover:text-white/80 transition-colors cursor-pointer"
          >
            <a
              href="https://github.com/gocanto/bedrock"
              target="_blank"
              rel="noopener noreferrer"
            >
              Open Source · Go · MIT
            </a>
          </Badge>

          <div>
            <h1 class="text-4xl font-bold leading-tight tracking-tight text-white sm:text-5xl lg:text-6xl">
              {{ frontmatter.heroText ?? site.title }}
            </h1>
            <p class="mt-5 max-w-md text-base leading-relaxed text-white/60 md:text-lg">
              {{ frontmatter.tagline ?? site.description }}
            </p>
          </div>

          <!-- CTA buttons -->
          <div v-if="frontmatter.actions?.length" class="flex flex-wrap gap-3">
            <RouteLink
              v-for="action in frontmatter.actions"
              :key="action.link"
              :to="action.link"
              :class="cn(
                'inline-flex items-center gap-1.5 rounded-lg px-5 py-2.5 text-sm font-semibold transition-all',
                action.type === 'primary'
                  ? 'bg-white text-zinc-900 hover:bg-zinc-100 shadow-sm'
                  : 'border border-white/20 text-white hover:border-white/40 hover:bg-white/5',
              )"
            >
              {{ action.text }}
              <ChevronRight v-if="action.type === 'primary'" class="h-3.5 w-3.5" />
            </RouteLink>
          </div>
        </div>

        <!-- Right: install code group -->
        <div class="flex flex-col justify-center p-10 md:py-20 md:pl-16">
          <Tabs :default-value="managers[0].label">
            <TabsList
              class="w-full justify-start gap-1 border-b border-white/10
                     bg-transparent rounded-none h-auto p-0 pb-3 mb-3
                     text-white/40"
            >
              <TabsTrigger
                v-for="m in managers"
                :key="m.label"
                :value="m.label"
                class="rounded px-3 py-1 font-mono text-xs transition-colors
                       bg-transparent shadow-none
                       text-white/40 hover:text-white/70
                       data-[state=active]:bg-white/10
                       data-[state=active]:text-white
                       data-[state=active]:shadow-none"
              >
                {{ m.label }}
              </TabsTrigger>
            </TabsList>

            <TabsContent
              v-for="m in managers"
              :key="m.label"
              :value="m.label"
              class="mt-0"
            >
              <!-- Command block -->
              <div class="group rounded-xl bg-white/5 ring-1 ring-white/10 p-5 flex items-start gap-3">
                <pre class="flex-1 min-w-0 overflow-x-auto font-mono text-sm leading-relaxed [&::-webkit-scrollbar]:hidden [scrollbar-width:none]"><span class="text-white/40 select-none">$ </span><span class="text-emerald-400">{{ m.cmd }}</span></pre>
                <Button
                  variant="ghost"
                  size="icon"
                  class="shrink-0 h-7 w-7 text-white/30 opacity-0 group-hover:opacity-100
                         hover:bg-white/10 hover:text-white"
                  aria-label="Copy"
                  @click="() => navigator.clipboard.writeText(m.cmd)"
                >
                  <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5" aria-hidden="true">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M15.666 3.888A2.25 2.25 0 0013.5 2.25h-3c-1.03 0-1.9.693-2.166 1.638m7.332 0c.055.194.084.4.084.612v0a.75.75 0 01-.75.75H9a.75.75 0 01-.75-.75v0c0-.212.03-.418.084-.612m7.332 0c.646.049 1.288.11 1.927.184 1.1.128 1.907 1.077 1.907 2.185V19.5a2.25 2.25 0 01-2.25 2.25H6.75A2.25 2.25 0 014.5 19.5V6.257c0-1.108.806-2.057 1.907-2.185a48.208 48.208 0 011.927-.184" />
                  </svg>
                </Button>
              </div>
            </TabsContent>
          </Tabs>

          <!-- Stats row -->
          <div class="mt-8 grid grid-cols-3 gap-4 border-t border-white/10 pt-8">
            <div v-for="stat in stats" :key="stat.label" class="text-center">
              <div class="text-xl font-bold text-white">{{ stat.value }}</div>
              <div class="mt-0.5 text-xs text-white/40">{{ stat.label }}</div>
            </div>
          </div>
        </div>

      </div>
    </section>

    <!-- ─── Feature grid ─────────────────────────────────────────────────── -->
    <section v-if="frontmatter.features?.length" class="mx-auto w-full max-w-6xl px-10 py-20">
      <div class="mb-12">
        <p class="font-mono text-xs uppercase tracking-widest text-white/40">Packages</p>
        <h2 class="mt-2 text-3xl font-bold text-white">
          Redefining Go web development
        </h2>
      </div>

      <div class="rounded-2xl border border-white/10 overflow-hidden">
      <div class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3">
        <RouteLink
          v-for="feature in frontmatter.features"
          :key="feature.title"
          :to="`/packages/${feature.title}`"
          class="block border-b border-white/10 sm:border-r sm:[&:nth-child(2n)]:border-r-0 xl:border-r xl:[&:nth-child(3n)]:border-r-0"
        >
          <Card
            class="group relative rounded-none border-0 bg-[#0f0f12]
                   hover:bg-white/[0.03] transition-colors cursor-pointer h-full"
          >
            <!-- Subtle left accent on hover -->
            <div
              class="pointer-events-none absolute inset-y-0 left-0 w-px bg-emerald-400 opacity-0 group-hover:opacity-100 transition-opacity"
              aria-hidden="true"
            />
            <CardHeader class="pb-2">
              <CardTitle class="flex items-start justify-between">
                <span class="font-mono text-sm font-semibold text-white">
                  {{ feature.title }}
                </span>
                <ChevronRight
                  class="mt-0.5 h-4 w-4 shrink-0 text-white/20 group-hover:text-emerald-400 transition-colors"
                />
              </CardTitle>
            </CardHeader>
            <CardContent>
              <CardDescription class="text-sm leading-relaxed text-white/50">
                {{ feature.details }}
              </CardDescription>
            </CardContent>
          </Card>
        </RouteLink>
      </div>
      </div>
    </section>

    <!-- ─── Footer CTA ───────────────────────────────────────────────────── -->
    <section class="border-t border-white/10 bg-[#0c0c0f]">
      <div class="mx-auto max-w-6xl px-10 py-20 text-center">
        <p class="mx-auto max-w-2xl text-2xl font-semibold leading-snug text-white/80">
          Prepare for a development environment that can finally keep pace with the speed of your mind.
        </p>
        <Button
          variant="default"
          as-child
          class="mt-8 bg-white text-zinc-900 hover:bg-zinc-100"
        >
          <RouteLink to="/getting-started">
            Get started
            <ChevronRight class="h-4 w-4" />
          </RouteLink>
        </Button>
      </div>
    </section>

    <!-- Footer -->
    <footer class="border-t border-white/10 py-8 text-center text-xs text-white/30">
      {{ frontmatter.footer ?? '© Bedrock contributors' }}
    </footer>
  </div>
</template>
