<script setup lang="ts">
import { useRoute } from 'vue-router'
import { RouteLink } from 'vuepress/client'
import { Menu, Github, Search } from 'lucide-vue-next'
import { Button } from './ui/button'
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
  <header
    class="sticky top-0 z-30 flex h-14 items-center gap-4 border-b
           border-slate-200 dark:border-slate-800
           bg-white/80 dark:bg-slate-900/80 backdrop-blur-sm
           px-4 lg:px-6"
  >
    <!-- Mobile: hamburger -->
    <Button
      variant="ghost"
      size="icon"
      class="lg:hidden -ml-1 text-slate-500 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white"
      aria-label="Open sidebar"
      @click="emit('toggle-sidebar')"
    >
      <Menu class="h-5 w-5" />
    </Button>

    <!-- Desktop nav links -->
    <nav class="hidden lg:flex items-center gap-0.5 shrink-0" aria-label="Top navigation">
      <Button
        v-for="item in navbarItems"
        :key="item.link ?? item.text"
        variant="ghost"
        size="sm"
        as-child
        :class="[
          'text-sm',
          isActive(item.link)
            ? 'text-slate-900 dark:text-white bg-slate-900/5 dark:bg-white/5'
            : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white hover:bg-slate-900/5 dark:hover:bg-white/5',
        ]"
      >
        <RouteLink :to="item.link ?? '#'">{{ item.text }}</RouteLink>
      </Button>
    </nav>

    <!-- Search -->
    <button
      type="button"
      class="hidden lg:flex items-center gap-2 ml-4 h-8 rounded-full px-3
             ring-1 ring-slate-200 dark:ring-slate-700
             text-sm text-slate-400 dark:text-slate-500
             hover:ring-slate-300 dark:hover:ring-slate-600
             transition-all cursor-pointer min-w-[200px] bg-white dark:bg-slate-800/60"
      aria-label="Search documentation"
    >
      <Search class="h-3.5 w-3.5 shrink-0 text-slate-400 dark:text-slate-500" />
      <span class="flex-1 text-left text-sm">Search docs</span>
      <kbd class="font-mono text-[10px] text-slate-300 dark:text-slate-600">⌘K</kbd>
    </button>

    <!-- Right side -->
    <div class="ml-auto flex items-center gap-1">
      <Button
        variant="ghost"
        size="icon"
        as-child
        class="text-slate-400 hover:text-slate-600 dark:hover:text-slate-300"
        aria-label="GitHub"
      >
        <a href="https://github.com/gocanto/bedrock" target="_blank" rel="noopener noreferrer">
          <Github class="h-5 w-5" />
        </a>
      </Button>
      <DarkModeToggle />
    </div>
  </header>
</template>
