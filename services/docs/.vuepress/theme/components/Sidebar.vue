<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute } from 'vue-router'
import { RouteLink } from 'vuepress/client'
import { ChevronRight } from 'lucide-vue-next'
import { Button } from './ui/button'
import { Collapsible, CollapsibleTrigger, CollapsibleContent } from './ui/collapsible'
import { ScrollArea } from './ui/scroll-area'
import { Separator } from './ui/separator'
import { cn } from '../lib/utils'

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
  const norm = (p: string) => p.replace(/\.html$/, '').replace(/\/$/, '')
  return norm(route.path) === norm(link)
}
</script>

<template>
  <ScrollArea class="h-full">
    <nav class="text-sm px-3 py-4" aria-label="Sidebar navigation">
      <template v-for="(section, idx) in config" :key="section.text">
        <Separator v-if="idx > 0" class="my-3 bg-slate-100 dark:bg-slate-800" />

        <!-- Section with children (collapsible) -->
        <Collapsible
          v-if="section.children"
          :open="isOpen(section)"
          @update:open="() => toggle(section)"
        >
          <CollapsibleTrigger as-child>
            <Button
              variant="ghost"
              size="sm"
              class="w-full justify-between px-3 text-[11px] font-semibold uppercase tracking-widest
                     text-slate-500 dark:text-slate-400
                     hover:text-slate-700 dark:hover:text-slate-300
                     hover:bg-slate-900/5 dark:hover:bg-white/5"
            >
              {{ section.text }}
              <ChevronRight
                :class="cn(
                  'h-3 w-3 shrink-0 transition-transform duration-200',
                  isOpen(section) ? 'rotate-90' : '',
                )"
              />
            </Button>
          </CollapsibleTrigger>

          <CollapsibleContent>
            <ul class="mt-1 space-y-0.5">
              <li v-for="item in section.children" :key="item.link ?? item.text">
                <Button
                  v-if="item.link"
                  variant="ghost"
                  size="sm"
                  as-child
                  :class="cn(
                    'w-full justify-start pl-4 font-normal',
                    isActive(item.link)
                      ? 'text-sky-500 dark:text-sky-400 font-medium hover:text-sky-500 dark:hover:text-sky-400'
                      : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white',
                  )"
                >
                  <RouteLink :to="item.link" @click="emit('navigate')">
                    {{ item.text }}
                  </RouteLink>
                </Button>
                <span
                  v-else
                  class="block px-4 py-1.5 text-sm text-slate-400 dark:text-slate-600 select-none"
                >
                  {{ item.text }}
                </span>
              </li>
            </ul>
          </CollapsibleContent>
        </Collapsible>

        <!-- Top-level link (no children) -->
        <Button
          v-else-if="section.link"
          variant="ghost"
          size="sm"
          as-child
          :class="cn(
            'w-full justify-start font-normal',
            isActive(section.link)
              ? 'text-sky-500 dark:text-sky-400 font-medium hover:text-sky-500 dark:hover:text-sky-400'
              : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white',
          )"
        >
          <RouteLink :to="section.link" @click="emit('navigate')">
            {{ section.text }}
          </RouteLink>
        </Button>
      </template>
    </nav>
  </ScrollArea>
</template>
