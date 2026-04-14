<script setup lang="ts">
import { RouteLink } from 'vuepress/client'
</script>

<template>
  <div class="flex min-h-screen flex-col items-center justify-center bg-white dark:bg-zinc-900 px-6 text-center">
    <p
      class="text-8xl font-extrabold tracking-tight bg-clip-text text-transparent"
      style="background-image: linear-gradient(to right, #36b49f, #DBFF75);"
    >
      404
    </p>

    <h1 class="mt-6 text-3xl font-bold text-zinc-900 dark:text-white">
      Page not found
    </h1>

    <p class="mt-4 max-w-sm text-zinc-500 dark:text-zinc-400">
      The page you're looking for doesn't exist or has been moved.
    </p>

    <RouteLink
      to="/"
      class="mt-8 inline-flex items-center rounded-full bg-zinc-900 px-5 py-2.5
             text-sm font-semibold text-white transition
             hover:bg-zinc-700
             dark:bg-white dark:text-zinc-900 dark:hover:bg-zinc-200"
    >
      Back to home
    </RouteLink>
  </div>
</template>
