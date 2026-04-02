<script setup lang="ts">
import { describeAuthWorkspace } from "@gollin/auth";

import { Button } from "@/components/ui/button";

const authStatus = describeAuthWorkspace();
</script>

<template>
  <main class="min-h-screen bg-slate-950 px-6 py-12 text-slate-50">
    <section class="mx-auto grid max-w-5xl gap-10 lg:grid-cols-[1.3fr_0.7fr]">
      <div class="space-y-8">
        <div
          class="inline-flex items-center rounded-full border border-emerald-400/30 bg-emerald-400/10 px-3 py-1 text-sm text-emerald-200"
        >
          Monorepo bootstrap complete
        </div>
        <div class="space-y-4">
          <p class="text-sm uppercase tracking-[0.35em] text-slate-400">Dashboard</p>
          <h1 class="max-w-2xl text-4xl font-semibold tracking-tight text-white sm:text-5xl">
            Workspace-aware Vue apps on top of pnpm and Turborepo.
          </h1>
          <p class="max-w-2xl text-lg leading-8 text-slate-300">
            This app consumes an internal package through the workspace graph, uses Tailwind v4
            styling, and renders a `shadcn-vue` button component from local source.
          </p>
        </div>
        <div class="flex flex-wrap gap-4">
          <Button size="lg">Open workspace</Button>
          <Button size="lg" variant="secondary">Review package graph</Button>
        </div>
      </div>

      <aside
        class="rounded-3xl border border-white/10 bg-white/5 p-6 shadow-2xl shadow-slate-950/30"
      >
        <p class="text-sm uppercase tracking-[0.3em] text-slate-400">Internal package</p>
        <p class="mt-4 text-2xl font-medium text-white">{{ authStatus.label }}</p>
        <p class="mt-3 text-sm leading-7 text-slate-300">
          {{ authStatus.message }}
        </p>
        <dl class="mt-6 grid gap-4 text-sm text-slate-300">
          <div class="rounded-2xl border border-white/10 bg-slate-900/60 p-4">
            <dt class="text-slate-400">Workspace</dt>
            <dd class="mt-1 font-medium text-white">{{ authStatus.workspace }}</dd>
          </div>
          <div class="rounded-2xl border border-white/10 bg-slate-900/60 p-4">
            <dt class="text-slate-400">Runtime</dt>
            <dd class="mt-1 font-medium text-white">Vite + Vue 3 + TypeScript</dd>
          </div>
        </dl>
      </aside>
    </section>

    <article
      class="prose prose-invert prose-slate mx-auto mt-12 max-w-3xl rounded-3xl border border-white/10 bg-white/5 p-8"
    >
      <h2>Typography plugin smoke test</h2>
      <p>
        Tailwind Typography is active in this workspace. The prose block is intentionally present so
        `@tailwindcss/typography` is exercised in the first build instead of being a dead
        dependency.
      </p>
    </article>
  </main>
</template>
