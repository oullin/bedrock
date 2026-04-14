<script setup lang="ts">
import { useRoute } from 'vue-router'
import { RouteLink } from 'vuepress/client'
import DarkModeToggle from './DarkModeToggle.vue'

interface NavbarItem {
  text: string
  link?: string
}

defineProps<{ navbarItems: NavbarItem[] }>()
const emit = defineEmits<{ 'toggle-sidebar': [] }>()

const route = useRoute()

function isActive(link?: string): boolean {
  if (!link || link === '/') return route.path === '/'
  return route.path.startsWith(link)
}
</script>

<template>
  <!--
    Protocol navbar: bg-white/50 dark:bg-zinc-900/50 + backdrop-blur-sm.
    Border: border-zinc-900/10 dark:border-white/10.
    Height: h-14 (56px — matches Protocol).
  -->
  <header
    class="sticky top-0 z-30 flex h-14 items-center gap-3 border-b
           border-zinc-900/10 dark:border-white/10
           bg-white/[0.5] dark:bg-zinc-900/[0.5] backdrop-blur-sm
           px-4 lg:px-6"
  >
    <!-- Mobile: hamburger -->
    <button
      type="button"
      class="lg:hidden -ml-1 rounded-md p-2
             text-zinc-500 dark:text-zinc-400
             hover:text-zinc-900 dark:hover:text-white transition-colors"
      aria-label="Open sidebar"
      @click="emit('toggle-sidebar')"
    >
      <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2" aria-hidden="true">
        <path stroke-linecap="round" stroke-linejoin="round" d="M4 6h16M4 12h16M4 18h16" />
      </svg>
    </button>

    <!-- Desktop nav links -->
    <nav class="hidden lg:flex items-center gap-0.5" aria-label="Top navigation">
      <RouteLink
        v-for="item in navbarItems"
        :key="item.link ?? item.text"
        :to="item.link ?? '#'"
        :class="[
          'rounded-lg px-3 py-1.5 text-sm transition-colors',
          isActive(item.link)
            ? 'text-zinc-900 dark:text-white bg-zinc-900/5 dark:bg-white/5'
            : 'text-zinc-600 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-white hover:bg-zinc-900/5 dark:hover:bg-white/5',
        ]"
      >
        {{ item.text }}
      </RouteLink>
    </nav>

    <!-- Right side -->
    <div class="ml-auto flex items-center gap-3">
      <!-- GitHub icon link -->
      <a
        v-if="true"
        href="https://github.com/gocanto/bedrock"
        target="_blank"
        rel="noopener noreferrer"
        class="text-zinc-400 hover:text-zinc-600 dark:hover:text-zinc-300 transition-colors"
        aria-label="GitHub"
      >
        <svg class="h-5 w-5" fill="currentColor" viewBox="0 0 24 24" aria-hidden="true">
          <path fill-rule="evenodd" d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z" clip-rule="evenodd" />
        </svg>
      </a>
      <DarkModeToggle />
    </div>
  </header>
</template>
