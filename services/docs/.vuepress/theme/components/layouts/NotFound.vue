<script setup lang="ts">
import { RouteLink } from 'vuepress/client'
</script>

<template>
  <div class="flex min-h-screen flex-col items-center justify-center bg-white dark:bg-slate-900 px-6 text-center">
    <p
      class="text-8xl font-normal tracking-tight bg-clip-text text-transparent"
      style="font-family: var(--font-display); background-image: linear-gradient(to right, #06b6d4, #6366f1);"
    >
      404
    </p>

    <h1
      class="mt-6 text-3xl font-normal text-slate-900 dark:text-white"
      style="font-family: var(--font-display);"
    >
      Page not found
    </h1>

    <p class="mt-4 max-w-sm text-slate-500 dark:text-slate-400">
      The page you're looking for doesn't exist or has been moved.
    </p>

    <RouteLink
      to="/"
      class="mt-8 inline-flex items-center gap-2 rounded-full px-5 py-2.5 text-sm font-semibold
             bg-sky-500 text-white hover:bg-sky-600 shadow-sm shadow-sky-500/30 transition-colors"
      style="font-family: var(--font-display);"
    >
      Back to home
    </RouteLink>
  </div>
</template>
