<script setup lang="ts">
import { usePageFrontmatter, useSiteData, RouteLink } from 'vuepress/client'
import { ChevronRight, BookOpen, Package, Lightbulb, FlaskConical, Github } from 'lucide-vue-next'
import Navbar from '../Navbar.vue'

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

// Quick links — matches the Syntax template card grid
const quickLinks = [
  {
    title: 'Getting started',
    description: 'Everything you need to install Bedrock and start building your first Go web application.',
    href: '/getting-started',
    icon: BookOpen,
  },
  {
    title: 'Browse packages',
    description: 'Explore 11 production-grade packages for auth, caching, routing, queues, and more.',
    href: '/packages/container',
    icon: Package,
  },
  {
    title: 'Core concepts',
    description: 'Understand the request lifecycle, service container, and foundational architecture.',
    href: '/concepts/request-lifecycle',
    icon: Lightbulb,
  },
  {
    title: 'Testing guide',
    description: 'Learn how to write reliable unit, integration, and feature tests for Bedrock apps.',
    href: '/concepts/testing',
    icon: FlaskConical,
  },
]

// Stats
const stats = [
  { value: '11',       label: 'Packages' },
  { value: 'Go 1.24+', label: 'Requires' },
  { value: 'MIT',      label: 'License' },
]
</script>

<template>
  <div class="flex min-h-full flex-col bg-white antialiased dark:bg-slate-900">
    <Navbar :navbar-items="themeOptions.navbar ?? []" @toggle-sidebar="() => {}" />

    <!-- ─── Hero (CacheAdvance-style dark header) ────────────────────────── -->
    <header class="relative overflow-hidden bg-[#0b1120]">

      <!-- Circuit-board grid -->
      <div
        class="pointer-events-none absolute inset-0"
        style="
          background-image:
            linear-gradient(rgba(56,189,248,0.06) 1px, transparent 1px),
            linear-gradient(90deg, rgba(56,189,248,0.06) 1px, transparent 1px);
          background-size: 64px 64px;
        "
        aria-hidden="true"
      />
      <!-- Radial glow -->
      <div
        class="pointer-events-none absolute inset-0"
        style="background: radial-gradient(ellipse 70% 50% at 40% 50%, rgba(56,189,248,0.07) 0%, transparent 70%);"
        aria-hidden="true"
      />
      <!-- Circuit nodes + lines -->
      <svg class="pointer-events-none absolute inset-0 h-full w-full" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
        <circle cx="512" cy="64"  r="3" fill="rgba(56,189,248,0.3)" />
        <circle cx="768" cy="192" r="3" fill="rgba(99,102,241,0.3)" />
        <circle cx="960" cy="96"  r="2" fill="rgba(56,189,248,0.2)" />
        <circle cx="320" cy="256" r="2" fill="rgba(99,102,241,0.2)" />
        <circle cx="640" cy="320" r="3" fill="rgba(56,189,248,0.2)" />
        <circle cx="1088" cy="224" r="2" fill="rgba(56,189,248,0.15)" />
        <line x1="512" y1="64"  x2="640" y2="64"  stroke="rgba(56,189,248,0.12)" stroke-width="1"/>
        <line x1="640" y1="64"  x2="640" y2="192" stroke="rgba(56,189,248,0.1)"  stroke-width="1"/>
        <line x1="768" y1="192" x2="960" y2="192" stroke="rgba(56,189,248,0.1)"  stroke-width="1"/>
        <line x1="960" y1="96"  x2="960" y2="192" stroke="rgba(99,102,241,0.1)"  stroke-width="1"/>
        <line x1="960" y1="192" x2="1088" y2="224" stroke="rgba(56,189,248,0.08)" stroke-width="1"/>
      </svg>

      <div class="relative mx-auto grid max-w-6xl grid-cols-1 items-center gap-12 px-8 py-20 md:grid-cols-2 lg:px-12">

        <!-- Left: headline + CTAs -->
        <div class="flex flex-col gap-6">
          <h1
            class="text-5xl font-bold leading-[1.1] tracking-tight text-white sm:text-6xl"
            style="font-family: var(--font-display);"
          >
            {{ frontmatter.heroText ?? site.title }}&thinsp;–<br>
            <span class="text-sky-400">the Upstream way.</span>
          </h1>

          <p class="max-w-md text-base leading-relaxed text-slate-400">
            {{ frontmatter.tagline ?? site.description }}
          </p>

          <div class="flex flex-wrap gap-3">
            <RouteLink
              to="/getting-started"
              class="inline-flex items-center gap-2 rounded-full px-5 py-2.5 text-sm font-semibold
                     bg-sky-500 text-white hover:bg-sky-400 shadow-lg shadow-sky-500/20 transition-colors"
              style="font-family: var(--font-display);"
            >
              Get started
              <ChevronRight class="h-3.5 w-3.5" />
            </RouteLink>
            <a
              href="https://github.com/gocanto/bedrock"
              target="_blank"
              rel="noopener noreferrer"
              class="inline-flex items-center gap-2 rounded-full px-5 py-2.5 text-sm font-semibold
                     ring-1 ring-white/15 text-slate-300 hover:ring-white/30 hover:text-white hover:bg-white/5 transition-all"
              style="font-family: var(--font-display);"
            >
              <Github class="h-4 w-4" />
              View on GitHub
            </a>
          </div>

          <!-- Stats -->
          <div class="flex gap-10 border-t border-white/10 pt-6 mt-2">
            <div v-for="stat in stats" :key="stat.label">
              <div class="text-xl font-semibold text-white" style="font-family: var(--font-display);">
                {{ stat.value }}
              </div>
              <div class="mt-0.5 text-xs text-slate-500">{{ stat.label }}</div>
            </div>
          </div>
        </div>

        <!-- Right: code editor card -->
        <div class="relative rounded-2xl overflow-hidden shadow-2xl ring-1 ring-white/10" style="background: #0d1b2e;">
          <!-- Traffic lights -->
          <div class="flex items-center gap-1.5 border-b border-white/[0.06] px-4 py-3">
            <span class="h-3 w-3 rounded-full bg-[#ff5f57]" />
            <span class="h-3 w-3 rounded-full bg-[#febc2e]" />
            <span class="h-3 w-3 rounded-full bg-[#28c840]" />
          </div>
          <!-- File tabs -->
          <div class="flex border-b border-white/[0.06]">
            <div class="px-4 py-2 text-xs font-mono text-slate-200 bg-white/[0.04] border-r border-white/[0.06] border-b-2 border-b-sky-400 cursor-default">
              container.go
            </div>
            <div class="px-4 py-2 text-xs font-mono text-slate-500 border-r border-white/[0.06] cursor-default hover:text-slate-400">
              routes.go
            </div>
          </div>
          <!-- Code -->
          <div class="overflow-x-auto p-5">
            <div class="flex gap-5 font-mono text-sm leading-6">
              <div class="select-none text-right text-slate-600 shrink-0" aria-hidden="true">
                <div v-for="n in 18" :key="n">{{ n }}</div>
              </div>
              <pre class="flex-1 min-w-0 text-slate-300 overflow-x-auto [&::-webkit-scrollbar]:hidden [scrollbar-width:none]"><span class="text-violet-400">package </span><span class="text-sky-400">main</span>

<span class="text-violet-400">import </span><span class="text-slate-300">(</span>
    <span class="text-green-400">"github.com/gocanto/bedrock/</span><span class="text-sky-400">container</span><span class="text-green-400">"</span>
    <span class="text-green-400">"github.com/gocanto/bedrock/</span><span class="text-sky-400">routing</span><span class="text-green-400">"</span>
    <span class="text-green-400">"github.com/gocanto/bedrock/</span><span class="text-sky-400">auth</span><span class="text-green-400">"</span>
<span class="text-slate-300">)</span>

<span class="text-violet-400">func </span><span class="text-sky-400">main</span><span class="text-slate-300">() {</span>
    <span class="text-slate-300">app </span><span class="text-cyan-400">:= </span><span class="text-sky-400">container</span><span class="text-slate-300">.New()</span>

    <span class="text-slate-300">app.Register(</span><span class="text-violet-400">func</span><span class="text-slate-300">(c </span><span class="text-yellow-300">*container.Container</span><span class="text-slate-300">) {</span>
        <span class="text-slate-300">c.Singleton(</span><span class="text-green-400">"router"</span><span class="text-slate-300">, routing.New())</span>
        <span class="text-slate-300">c.Singleton(</span><span class="text-green-400">"auth"</span><span class="text-slate-300">,   auth.New(c))</span>
    <span class="text-slate-300">})</span>

    <span class="text-slate-300">app.</span><span class="text-sky-400">Boot</span><span class="text-slate-300">()</span>
<span class="text-slate-300">}</span></pre>
            </div>
          </div>
        </div>

      </div>
    </header>

    <!-- ─── Quick links grid ─────────────────────────────────────────────── -->
    <div class="relative mx-auto max-w-5xl px-6 pb-16 lg:px-8">
      <div class="grid grid-cols-1 gap-6 sm:grid-cols-2">
        <RouteLink
          v-for="card in quickLinks"
          :key="card.title"
          :to="card.href"
          class="group relative rounded-xl border border-slate-200 p-6 dark:border-slate-800
                 hover:border-transparent dark:hover:border-transparent transition-colors duration-300"
          style="--card-bg: #ffffff;"
        >
          <!-- Syntax-style gradient border on hover -->
          <div
            class="pointer-events-none absolute inset-0 rounded-xl opacity-0 group-hover:opacity-100 transition-opacity duration-300 dark:hidden"
            style="background: linear-gradient(#ffffff, #ffffff) padding-box,
                   linear-gradient(115deg, transparent 20%, #6366f1 20% 25%, transparent 25% 35%, #06b6d4 35% 40%, transparent 40% 60%, #6366f1 60% 65%, transparent 65% 75%, #0ea5e9 75% 80%, transparent 80%) border-box;
                   border: 1px solid transparent;"
            aria-hidden="true"
          />
          <div
            class="pointer-events-none absolute inset-0 rounded-xl opacity-0 group-hover:opacity-100 transition-opacity duration-300 hidden dark:block"
            style="background: linear-gradient(rgb(15 23 42), rgb(15 23 42)) padding-box,
                   linear-gradient(115deg, transparent 20%, #6366f1 20% 25%, transparent 25% 35%, #06b6d4 35% 40%, transparent 40% 60%, #6366f1 60% 65%, transparent 65% 75%, #0ea5e9 75% 80%, transparent 80%) border-box;
                   border: 1px solid transparent;"
            aria-hidden="true"
          />

          <!-- Icon container with gradient fill -->
          <div
            class="relative mb-4 flex h-12 w-12 items-center justify-center rounded-lg"
            style="background: radial-gradient(circle at 30% 30%, #06b6d4, #6366f1);"
          >
            <component :is="card.icon" class="h-5 w-5 text-white" />
          </div>

          <h3
            class="relative text-sm font-semibold text-slate-900 dark:text-white"
            style="font-family: var(--font-display);"
          >
            {{ card.title }}
          </h3>
          <p class="relative mt-1 text-sm leading-relaxed text-slate-500 dark:text-slate-400">
            {{ card.description }}
          </p>
        </RouteLink>
      </div>
    </div>

    <!-- ─── Packages grid ───────────────────────────────────────────────── -->
    <div
      v-if="frontmatter.features?.length"
      class="border-t border-slate-200 dark:border-slate-800 bg-slate-50 dark:bg-slate-800/20"
    >
      <div class="mx-auto max-w-5xl px-6 py-16 lg:px-8">
        <p class="mb-1 text-[11px] font-semibold uppercase tracking-widest text-slate-400 dark:text-slate-500">
          Packages
        </p>
        <h2
          class="mb-10 text-2xl font-normal text-slate-900 dark:text-white"
          style="font-family: var(--font-display);"
        >
          11 packages, ready to install individually
        </h2>

        <div class="grid grid-cols-1 divide-y divide-slate-200 dark:divide-slate-800 rounded-2xl border border-slate-200 dark:border-slate-800 overflow-hidden sm:grid-cols-2 sm:divide-x xl:grid-cols-3">
          <RouteLink
            v-for="(feature, idx) in frontmatter.features"
            :key="feature.title"
            :to="`/packages/${feature.title}`"
            :class="[
              'group relative block bg-white dark:bg-slate-900 p-6 transition-colors hover:bg-slate-50 dark:hover:bg-slate-800/60',
              'sm:[&:nth-child(2n)]:border-r-0 xl:[&:nth-child(2n)]:border-r xl:[&:nth-child(3n)]:border-r-0',
            ]"
          >
            <!-- Sky left-border accent on hover -->
            <div
              class="pointer-events-none absolute inset-y-0 left-0 w-0.5 bg-sky-500 opacity-0 group-hover:opacity-100 transition-opacity"
              aria-hidden="true"
            />
            <div class="flex items-start justify-between">
              <span class="font-mono text-sm font-semibold text-slate-900 dark:text-white">
                {{ feature.title }}
              </span>
              <ChevronRight
                class="mt-0.5 h-4 w-4 shrink-0 text-slate-300 dark:text-slate-600 group-hover:text-sky-500 transition-colors"
              />
            </div>
            <p class="mt-2 text-sm leading-relaxed text-slate-500 dark:text-slate-400">
              {{ feature.details }}
            </p>
          </RouteLink>
        </div>
      </div>
    </div>

    <!-- ─── Footer CTA ───────────────────────────────────────────────────── -->
    <div class="border-t border-slate-200 dark:border-slate-800">
      <div class="mx-auto max-w-5xl px-6 py-20 lg:px-8 text-center">
        <p
          class="mx-auto max-w-2xl text-2xl font-normal leading-snug text-slate-900 dark:text-white"
          style="font-family: var(--font-display);"
        >
          Prepare for a development environment that can finally keep pace with the speed of your mind.
        </p>
        <RouteLink
          to="/getting-started"
          class="mt-8 inline-flex items-center gap-2 rounded-full px-6 py-3 text-sm font-semibold
                 bg-sky-500 text-white hover:bg-sky-600 shadow-sm shadow-sky-500/30 transition-colors"
          style="font-family: var(--font-display);"
        >
          Get started
          <ChevronRight class="h-4 w-4" />
        </RouteLink>
      </div>
    </div>

    <!-- Footer -->
    <footer class="border-t border-slate-200 dark:border-slate-800 py-8 text-center text-xs text-slate-400 dark:text-slate-500">
      {{ frontmatter.footer ?? '© Bedrock contributors' }}
    </footer>
  </div>
</template>
