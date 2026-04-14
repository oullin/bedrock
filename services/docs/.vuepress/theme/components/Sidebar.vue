<script setup lang="ts">
import { ref } from 'vue'
import { useRoute } from 'vue-router'
import { RouteLink } from 'vuepress/client'

interface SidebarItem {
  text: string
  link?: string
  collapsible?: boolean
  children?: SidebarItem[]
}

const props = defineProps<{ config: SidebarItem[] }>()
const emit = defineEmits<{ navigate: [] }>()

const route = useRoute()

function containsActive(items: SidebarItem[]): boolean {
  return items.some(
    (item) =>
      (item.link && isActive(item.link)) ||
      (item.children && containsActive(item.children)),
  )
}

const openSections = ref<Record<string, boolean>>({})

function isOpen(section: SidebarItem): boolean {
  if (section.text in openSections.value) return openSections.value[section.text]
  return !section.collapsible || containsActive(section.children ?? [])
}

function toggle(section: SidebarItem) {
  openSections.value[section.text] = !isOpen(section)
}

function isActive(link?: string): boolean {
  if (!link) return false
  const norm = link.replace(/\.html$/, '')
  return route.path === link || route.path === norm
}
</script>

<template>
  <!--
    Protocol sidebar: text-only active states (no pill, no bg, no border accent).
    Active links use text-emerald-500 dark:text-emerald-400.
    Dividers between groups: divide-zinc-900/5 dark:divide-white/5.
  -->
  <nav
    class="text-sm divide-y divide-zinc-900/5 dark:divide-white/5"
    aria-label="Sidebar navigation"
  >
    <div v-for="section in config" :key="section.text" class="py-4 first:pt-0 last:pb-0">

      <!-- ── Section header (collapsible) ──────────────────────────────── -->
      <button
        v-if="section.children"
        type="button"
        class="flex w-full items-center justify-between px-3 py-1 mb-1 gap-2
               text-xs font-semibold text-zinc-900 dark:text-white
               hover:text-zinc-700 dark:hover:text-zinc-200 transition-colors"
        @click="toggle(section)"
      >
        <span>{{ section.text }}</span>
        <svg
          :class="['h-2.5 w-2.5 shrink-0 transition-transform duration-200', isOpen(section) ? 'rotate-90' : '']"
          fill="none" viewBox="0 0 6 10" aria-hidden="true"
        >
          <path d="M1 1l4 4-4 4" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" />
        </svg>
      </button>

      <!-- ── Top-level link (no children) ──────────────────────────────── -->
      <RouteLink
        v-else-if="section.link"
        :to="section.link"
        :class="[
          'block px-3 py-1 transition-colors',
          isActive(section.link)
            ? 'text-emerald-500 dark:text-emerald-400 font-medium'
            : 'text-zinc-600 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-white',
        ]"
        @click="emit('navigate')"
      >
        {{ section.text }}
      </RouteLink>

      <!-- ── Children list ──────────────────────────────────────────────── -->
      <Transition
        enter-active-class="transition-all duration-200 ease-out overflow-hidden"
        enter-from-class="opacity-0 max-h-0"
        enter-to-class="opacity-100 max-h-[800px]"
        leave-active-class="transition-all duration-150 ease-in overflow-hidden"
        leave-from-class="opacity-100 max-h-[800px]"
        leave-to-class="opacity-0 max-h-0"
      >
        <ul v-if="section.children && isOpen(section)">
          <li v-for="item in section.children" :key="item.link ?? item.text">
            <RouteLink
              v-if="item.link"
              :to="item.link"
              :class="[
                'flex py-1 pr-3 pl-4 transition-colors',
                isActive(item.link)
                  ? 'text-emerald-500 dark:text-emerald-400 font-medium'
                  : 'text-zinc-600 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-white',
              ]"
              @click="emit('navigate')"
            >
              {{ item.text }}
            </RouteLink>
            <span
              v-else
              class="block py-1 pl-4 pr-3 text-zinc-400 dark:text-zinc-600 select-none"
            >
              {{ item.text }}
            </span>
          </li>
        </ul>
      </Transition>
    </div>
  </nav>
</template>
