<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute } from 'vue-router'
import { Content, RouteLink, useSiteData } from 'vuepress/client'
import { ArrowLeft, ArrowRight, Github } from 'lucide-vue-next'
import Sidebar from '../Sidebar.vue'
import Navbar from '../Navbar.vue'
import TableOfContents from '../TableOfContents.vue'
import { Button } from '../ui/button'
import { Sheet, SheetContent, SheetTrigger, SheetTitle, SheetHeader } from '../ui/sheet'

declare const __BEDROCK_THEME_OPTIONS__: string
const themeOptions = JSON.parse(__BEDROCK_THEME_OPTIONS__)

const site = useSiteData()
const route = useRoute()
const sidebarOpen = ref(false)

interface SidebarItem {
  text: string
  link?: string
  collapsible?: boolean
  children?: SidebarItem[]
}

const sidebarConfig = computed<SidebarItem[]>(() => {
  const sidebar: Record<string, SidebarItem[]> = themeOptions.sidebar ?? {}
  const keys = Object.keys(sidebar).sort((a, b) => b.length - a.length)
  for (const key of keys) {
    if (route.path.startsWith(key)) return sidebar[key]
  }
  return sidebar['/'] ?? []
})

const flatLinks = computed<Array<{ text: string; link: string }>>(() => {
  const links: Array<{ text: string; link: string }> = []
  const walk = (items: SidebarItem[]) => {
    for (const item of items) {
      if (item.link) links.push({ text: item.text, link: item.link })
      if (item.children) walk(item.children)
    }
  }
  walk(sidebarConfig.value)
  return links
})

const currentIdx = computed(() =>
  flatLinks.value.findIndex(({ link }) => {
    const norm = link.replace(/\.html$/, '')
    return route.path === link || route.path === norm
  }),
)

const prevPage = computed(() =>
  currentIdx.value > 0 ? flatLinks.value[currentIdx.value - 1] : null,
)

const nextPage = computed(() =>
  currentIdx.value >= 0 && currentIdx.value < flatLinks.value.length - 1
    ? flatLinks.value[currentIdx.value + 1]
    : null,
)
</script>

<template>
  <div class="relative flex min-h-full bg-white antialiased dark:bg-slate-900">

    <!-- ── Mobile sidebar drawer ─────────────────────────────────────────── -->
    <Sheet v-model:open="sidebarOpen">
      <SheetContent
        side="left"
        class="flex w-72 flex-col p-0 bg-white dark:bg-slate-900
               border-r border-slate-200 dark:border-slate-800"
      >
        <SheetHeader class="sr-only">
          <SheetTitle>Navigation</SheetTitle>
        </SheetHeader>

        <div class="flex h-14 shrink-0 items-center border-b border-slate-200 dark:border-slate-800 px-6">
          <RouteLink
            to="/"
            class="flex items-center gap-2.5 text-sm font-semibold text-slate-900 dark:text-white hover:opacity-75 transition-opacity"
            @click="sidebarOpen = false"
          >
            <span
              class="flex h-6 w-6 items-center justify-center rounded text-xs font-bold text-white shrink-0"
              style="background: linear-gradient(135deg, #06b6d4, #6366f1);"
              aria-hidden="true"
            >B</span>
            {{ site.title }}
          </RouteLink>
        </div>

        <div class="flex-1 overflow-hidden">
          <Sidebar :config="sidebarConfig" @navigate="sidebarOpen = false" />
        </div>

        <div
          v-if="themeOptions.repo"
          class="shrink-0 border-t border-slate-200 dark:border-slate-800 px-6 py-4"
        >
          <a
            :href="themeOptions.repo"
            target="_blank"
            rel="noopener noreferrer"
            class="flex items-center gap-2 text-sm text-slate-500 hover:text-slate-900 dark:text-slate-400 dark:hover:text-white transition-colors"
          >
            <Github class="h-4 w-4 shrink-0" />
            GitHub
          </a>
        </div>
      </SheetContent>
    </Sheet>

    <!-- ── Desktop sidebar ───────────────────────────────────────────────── -->
    <aside
      class="hidden lg:fixed lg:inset-y-0 lg:left-0 lg:z-50 lg:flex lg:w-72 lg:flex-col
             bg-white dark:bg-slate-900
             border-r border-slate-200 dark:border-slate-800"
    >
      <div class="flex h-14 shrink-0 items-center border-b border-slate-200 dark:border-slate-800 px-6">
        <RouteLink
          to="/"
          class="flex items-center gap-2.5 text-sm font-semibold text-slate-900 dark:text-white hover:opacity-75 transition-opacity"
        >
          <span
            class="flex h-6 w-6 items-center justify-center rounded text-xs font-bold text-white shrink-0"
            style="background: linear-gradient(135deg, #06b6d4, #6366f1);"
            aria-hidden="true"
          >B</span>
          {{ site.title }}
        </RouteLink>
      </div>

      <div class="flex-1 overflow-hidden">
        <Sidebar :config="sidebarConfig" />
      </div>

      <div
        v-if="themeOptions.repo"
        class="shrink-0 border-t border-slate-200 dark:border-slate-800 px-6 py-4"
      >
        <a
          :href="themeOptions.repo"
          target="_blank"
          rel="noopener noreferrer"
          class="flex items-center gap-2 text-sm text-slate-500 hover:text-slate-900 dark:text-slate-400 dark:hover:text-white transition-colors"
        >
          <Github class="h-4 w-4 shrink-0" />
          GitHub
        </a>
      </div>
    </aside>

    <!-- ── Main content area ─────────────────────────────────────────────── -->
    <div class="relative z-10 flex min-w-0 flex-1 flex-col lg:pl-72">

      <Navbar
        :navbar-items="themeOptions.navbar ?? []"
        @toggle-sidebar="sidebarOpen = true"
      />

      <div class="flex flex-1 justify-center">

        <!-- Article -->
        <main class="min-w-0 flex-1 px-4 py-16 lg:pl-8 xl:px-16">
          <article
            class="mx-auto max-w-2xl min-w-0
                   prose prose-slate dark:prose-invert
                   prose-headings:font-display prose-headings:font-normal
                   prose-headings:scroll-mt-28 lg:prose-headings:scroll-mt-[8.5rem]
                   prose-lead:text-slate-500 dark:prose-lead:text-slate-400
                   prose-a:font-semibold dark:prose-a:text-sky-400
                   prose-a:no-underline
                   prose-pre:rounded-xl prose-pre:bg-slate-900 prose-pre:shadow-lg
                   dark:prose-pre:bg-slate-800/60 dark:prose-pre:ring-1 dark:prose-pre:ring-slate-300/10
                   dark:prose-hr:border-slate-800"
            vp-content
          >
            <Content />

            <!-- Prev / Next -->
            <div
              v-if="prevPage || nextPage"
              class="mt-16 flex items-stretch gap-3 border-t border-slate-200 dark:border-slate-800 pt-8 not-prose"
            >
              <Button
                v-if="prevPage"
                variant="outline"
                as-child
                class="flex-1 h-auto flex-col items-start gap-1 px-5 py-4 text-left
                       border-slate-200 dark:border-slate-800 hover:border-sky-400 dark:hover:border-sky-500 hover:bg-transparent"
              >
                <RouteLink :to="prevPage.link" class="w-full">
                  <span class="flex items-center gap-1 text-xs text-slate-400 dark:text-slate-500">
                    <ArrowLeft class="h-3 w-3" />
                    Previous
                  </span>
                  <span class="mt-1 block text-sm font-medium text-slate-900 dark:text-white">
                    {{ prevPage.text }}
                  </span>
                </RouteLink>
              </Button>

              <div v-if="!prevPage" class="flex-1" />

              <Button
                v-if="nextPage"
                variant="outline"
                as-child
                class="flex-1 h-auto flex-col items-end gap-1 px-5 py-4 text-right
                       border-slate-200 dark:border-slate-800 hover:border-sky-400 dark:hover:border-sky-500 hover:bg-transparent"
              >
                <RouteLink :to="nextPage.link" class="w-full text-right">
                  <span class="flex items-center justify-end gap-1 text-xs text-slate-400 dark:text-slate-500">
                    Next
                    <ArrowRight class="h-3 w-3" />
                  </span>
                  <span class="mt-1 block text-sm font-medium text-slate-900 dark:text-white">
                    {{ nextPage.text }}
                  </span>
                </RouteLink>
              </Button>
            </div>
          </article>
        </main>

        <!-- Right TOC — xl+ only -->
        <aside class="hidden xl:block w-56 shrink-0 pr-6">
          <div class="sticky top-[calc(3.5rem+2rem)] pt-1">
            <TableOfContents />
          </div>
        </aside>

      </div>
    </div>
  </div>
</template>
