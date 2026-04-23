<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute } from 'vue-router'
import { Content, RouteLink, useSiteData } from 'vuepress/client'
import { ArrowLeft, ArrowRight, Github } from 'lucide-vue-next'
import Sidebar from '../Sidebar.vue'
import Navbar from '../Navbar.vue'
import TableOfContents from '../TableOfContents.vue'
import BrandMark from '../home/BrandMark.vue'
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
  <div class="bedrock-docs relative flex min-h-full antialiased" style="background: var(--bg); color: var(--ink);">

    <!-- ── Mobile sidebar drawer ─────────────────────────────────────────── -->
    <Sheet v-model:open="sidebarOpen">
      <SheetContent
        side="left"
        class="bedrock-docs flex w-72 flex-col p-0 border-r"
        style="background: var(--panel); border-color: var(--line); color: var(--ink);"
      >
        <SheetHeader class="sr-only">
          <SheetTitle>Navigation</SheetTitle>
        </SheetHeader>

        <div class="flex h-14 shrink-0 items-center border-b px-6" style="border-color: var(--line);">
          <RouteLink
            to="/"
            class="flex items-center gap-2.5 text-sm font-semibold hover:opacity-75 transition-opacity"
            style="color: var(--ink);"
            @click="sidebarOpen = false"
          >
            <span class="brand-mark shrink-0" aria-hidden="true">
              <BrandMark :size="24" />
            </span>
            <span class="font-semibold lowercase" style="letter-spacing: -0.01em;">{{ site.title }}</span>
          </RouteLink>
        </div>

        <div class="flex-1 overflow-hidden">
          <Sidebar :config="sidebarConfig" @navigate="sidebarOpen = false" />
        </div>

        <div
          v-if="themeOptions.repo"
          class="shrink-0 border-t px-6 py-4"
          style="border-color: var(--line);"
        >
          <a
            :href="themeOptions.repo"
            target="_blank"
            rel="noopener noreferrer"
            class="flex items-center gap-2 text-sm transition-colors hover:opacity-100"
            style="color: var(--ink-3); opacity: 0.9;"
          >
            <Github class="h-4 w-4 shrink-0" />
            GitHub
          </a>
        </div>
      </SheetContent>
    </Sheet>

    <!-- ── Desktop sidebar ───────────────────────────────────────────────── -->
    <aside
      class="hidden lg:fixed lg:inset-y-0 lg:left-0 lg:z-50 lg:flex lg:w-72 lg:flex-col border-r"
      style="background: var(--panel); border-color: var(--line);"
    >
      <div class="flex h-14 shrink-0 items-center border-b px-6" style="border-color: var(--line);">
        <RouteLink
          to="/"
          class="flex items-center gap-2.5 text-sm font-semibold hover:opacity-75 transition-opacity"
          style="color: var(--ink);"
        >
          <span class="brand-mark shrink-0" aria-hidden="true">
            <BrandMark :size="24" />
          </span>
          <span class="font-semibold lowercase" style="letter-spacing: -0.01em;">{{ site.title }}</span>
        </RouteLink>
      </div>

      <div class="flex-1 overflow-hidden">
        <Sidebar :config="sidebarConfig" />
      </div>

      <div
        v-if="themeOptions.repo"
        class="shrink-0 border-t px-6 py-4"
        style="border-color: var(--line);"
      >
        <a
          :href="themeOptions.repo"
          target="_blank"
          rel="noopener noreferrer"
          class="flex items-center gap-2 text-sm transition-colors hover:opacity-100"
          style="color: var(--ink-3); opacity: 0.9;"
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
                   prose-headings:scroll-mt-28 lg:prose-headings:scroll-mt-[8.5rem]"
            vp-content
          >
            <Content />

            <!-- Prev / Next -->
            <div
              v-if="prevPage || nextPage"
              class="mt-16 flex items-stretch gap-3 border-t pt-8 not-prose"
              style="border-color: var(--line);"
            >
              <Button
                v-if="prevPage"
                variant="outline"
                as-child
                class="flex-1 h-auto flex-col items-start gap-1 px-5 py-4 text-left hover:bg-transparent"
                style="border-color: var(--line);"
                @mouseenter="(e: MouseEvent) => (e.currentTarget as HTMLElement).style.borderColor = 'var(--accent)'"
                @mouseleave="(e: MouseEvent) => (e.currentTarget as HTMLElement).style.borderColor = 'var(--line)'"
              >
                <RouteLink :to="prevPage.link" class="w-full">
                  <span class="flex items-center gap-1 text-xs" style="color: var(--ink-3);">
                    <ArrowLeft class="h-3 w-3" />
                    Previous
                  </span>
                  <span class="mt-1 block text-sm font-medium" style="color: var(--ink);">
                    {{ prevPage.text }}
                  </span>
                </RouteLink>
              </Button>

              <div v-if="!prevPage" class="flex-1" />

              <Button
                v-if="nextPage"
                variant="outline"
                as-child
                class="flex-1 h-auto flex-col items-end gap-1 px-5 py-4 text-right hover:bg-transparent"
                style="border-color: var(--line);"
                @mouseenter="(e: MouseEvent) => (e.currentTarget as HTMLElement).style.borderColor = 'var(--accent)'"
                @mouseleave="(e: MouseEvent) => (e.currentTarget as HTMLElement).style.borderColor = 'var(--line)'"
              >
                <RouteLink :to="nextPage.link" class="w-full text-right">
                  <span class="flex items-center justify-end gap-1 text-xs" style="color: var(--ink-3);">
                    Next
                    <ArrowRight class="h-3 w-3" />
                  </span>
                  <span class="mt-1 block text-sm font-medium" style="color: var(--ink);">
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
