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
    class="sticky top-0 z-30 flex h-14 items-center gap-4 border-b backdrop-blur-sm px-4 lg:px-6"
    style="border-color: var(--line); background-color: color-mix(in srgb, var(--bg) 85%, transparent);"
  >
    <!-- Mobile: hamburger -->
    <Button
      variant="ghost"
      size="icon"
      class="lg:hidden -ml-1 hover:bg-transparent"
      style="color: var(--ink-3);"
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
        class="text-sm hover:bg-transparent"
        :style="isActive(item.link)
          ? 'color: var(--ink);'
          : 'color: var(--ink-2);'"
      >
        <RouteLink :to="item.link ?? '#'">{{ item.text }}</RouteLink>
      </Button>
    </nav>

    <!-- Search -->
    <button
      type="button"
      class="hidden lg:flex items-center gap-2 ml-4 h-8 rounded-full px-3
             transition-all cursor-pointer min-w-[200px]"
      style="background: var(--panel); box-shadow: inset 0 0 0 1px var(--line); color: var(--ink-3);"
      aria-label="Search documentation"
    >
      <Search class="h-3.5 w-3.5 shrink-0" style="color: var(--ink-3);" />
      <span class="flex-1 text-left text-sm">Search docs</span>
      <kbd class="font-mono text-[10px]" style="color: var(--ink-3);">⌘K</kbd>
    </button>

    <!-- Right side -->
    <div class="ml-auto flex items-center gap-1">
      <Button
        variant="ghost"
        size="icon"
        as-child
        class="hover:bg-transparent"
        style="color: var(--ink-3);"
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
