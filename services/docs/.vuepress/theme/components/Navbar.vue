<script setup lang="ts">
import { useRoute } from 'vue-router'
import { RouteLink } from 'vuepress/client'
import { Menu, Github } from 'lucide-vue-next'
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
    class="sticky top-0 z-30 flex h-14 items-center gap-3 border-b
           border-zinc-900/10 dark:border-white/10
           bg-white/[0.5] dark:bg-zinc-900/[0.5] backdrop-blur-sm
           px-4 lg:px-6"
  >
    <!-- Mobile: hamburger -->
    <Button
      variant="ghost"
      size="icon"
      class="lg:hidden -ml-1 text-zinc-500 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-white"
      aria-label="Open sidebar"
      @click="emit('toggle-sidebar')"
    >
      <Menu class="h-5 w-5" />
    </Button>

    <!-- Desktop nav links -->
    <nav class="hidden lg:flex items-center gap-0.5" aria-label="Top navigation">
      <Button
        v-for="item in navbarItems"
        :key="item.link ?? item.text"
        variant="ghost"
        size="sm"
        as-child
        :class="[
          isActive(item.link)
            ? 'text-zinc-900 dark:text-white bg-zinc-900/5 dark:bg-white/5'
            : 'text-zinc-600 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-white hover:bg-zinc-900/5 dark:hover:bg-white/5',
        ]"
      >
        <RouteLink :to="item.link ?? '#'">
          {{ item.text }}
        </RouteLink>
      </Button>
    </nav>

    <!-- Right side -->
    <div class="ml-auto flex items-center gap-1">
      <Button
        variant="ghost"
        size="icon"
        as-child
        class="text-zinc-400 hover:text-zinc-600 dark:hover:text-zinc-300"
        aria-label="GitHub"
      >
        <a
          href="https://github.com/gocanto/bedrock"
          target="_blank"
          rel="noopener noreferrer"
        >
          <Github class="h-5 w-5" />
        </a>
      </Button>
      <DarkModeToggle />
    </div>
  </header>
</template>
