<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount, watch } from 'vue'
import { usePageData } from 'vuepress/client'
import { useRoute } from 'vue-router'
import { ScrollArea } from './ui/scroll-area'

const page = usePageData()
const route = useRoute()
const activeSlug = ref('')

interface FlatHeader {
  level: number
  title: string
  slug: string
}

const headers = computed<FlatHeader[]>(() => {
  const result: FlatHeader[] = []
  const walk = (items: typeof page.value.headers) => {
    for (const h of items ?? []) {
      if (h.level === 2 || h.level === 3) {
        result.push({ level: h.level, title: h.title, slug: h.slug })
      }
      if (h.children?.length) walk(h.children)
    }
  }
  walk(page.value.headers)
  return result
})

let observer: IntersectionObserver | null = null

function initObserver() {
  observer?.disconnect()
  if (typeof document === 'undefined') return
  const nodes = Array.from(
    document.querySelectorAll('[vp-content] h2[id], [vp-content] h3[id]'),
  )
  if (!nodes.length) return

  observer = new IntersectionObserver(
    (entries) => {
      const visible = entries
        .filter((e) => e.isIntersecting)
        .sort((a, b) => a.boundingClientRect.top - b.boundingClientRect.top)
      if (visible.length) activeSlug.value = visible[0].target.id
    },
    { rootMargin: '-56px 0px -60% 0px' },
  )
  nodes.forEach((n) => observer!.observe(n))
}

onMounted(initObserver)
onBeforeUnmount(() => observer?.disconnect())
watch(
  () => route.path,
  () => {
    activeSlug.value = ''
    setTimeout(initObserver, 350)
  },
)
</script>

<template>
  <nav v-if="headers.length" class="text-sm" aria-label="On this page">
    <p
      class="mb-4 text-[11px] font-semibold uppercase tracking-widest
             text-slate-900 dark:text-white"
      style="font-family: var(--font-display);"
    >
      On this page
    </p>
    <ScrollArea class="max-h-[calc(100vh-12rem)]">
      <ul class="space-y-2 border-l border-slate-200 dark:border-slate-800 pr-2">
        <li v-for="h in headers" :key="h.slug">
          <a
            :href="`#${h.slug}`"
            :class="[
              'block -ml-px border-l pl-3 text-[13px] leading-snug transition-colors duration-150',
              h.level === 3 ? 'pl-5' : '',
              activeSlug === h.slug
                ? 'border-sky-500 dark:border-sky-400 text-sky-500 dark:text-sky-400 font-semibold'
                : 'border-transparent text-slate-500 dark:text-slate-400 hover:text-slate-700 dark:hover:text-slate-300 hover:border-slate-400 dark:hover:border-slate-500',
            ]"
          >
            {{ h.title }}
          </a>
        </li>
      </ul>
    </ScrollArea>
  </nav>
</template>
