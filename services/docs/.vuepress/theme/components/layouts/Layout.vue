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

// Resolve sidebar config by longest-prefix match on current path
const sidebarConfig = computed<SidebarItem[]>(() => {
  const sidebar: Record<string, SidebarItem[]> = themeOptions.sidebar ?? {}
  const keys = Object.keys(sidebar).sort((a, b) => b.length - a.length)
  for (const key of keys) {
    if (route.path.startsWith(key)) return sidebar[key]
  }
  return sidebar['/'] ?? []
})

// Flatten sidebar to ordered list for prev/next navigation
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
  <div class="relative flex min-h-full bg-white antialiased dark:bg-zinc-900">

    <!-- ── Protocol gradient + skewed grid decoration ───────────────────── -->
    <div
      class="pointer-events-none absolute inset-x-0 top-0 z-0 h-[26rem] overflow-hidden
             [mask-image:linear-gradient(white,transparent)] dark:[mask-image:none]"
      aria-hidden="true"
    >
      <div
        class="absolute inset-0 bg-gradient-to-r from-[#36b49f] to-[#DBFF75]
               opacity-40 dark:from-[#36b49f]/30 dark:to-[#DBFF75]/30 dark:opacity-100
               [mask-image:radial-gradient(farthest-side_at_top_left,white,transparent)]"
      >
        <svg
          class="absolute inset-x-0 inset-y-[-50%] h-[200%] w-full skew-y-[-18deg]
                 fill-black/40 stroke-black/50
                 dark:fill-white/2.5 dark:stroke-white/5
                 mix-blend-overlay"
          aria-hidden="true"
        >
          <defs>
            <pattern id="gp" x="-12" y="4" width="72" height="56" patternUnits="userSpaceOnUse">
              <path d="M.5 56V.5H72" fill="none" />
            </pattern>
          </defs>
          <rect width="100%" height="100%" fill="url(#gp)" />
          <rect x="144" y="56"  width="72" height="56" />
          <rect x="288" y="112" width="72" height="56" />
          <rect x="360" y="0"   width="72" height="56" />
          <rect x="576" y="112" width="72" height="56" />
          <rect x="720" y="56"  width="72" height="56" />
          <rect x="864" y="168" width="72" height="56" />
        </svg>
      </div>
    </div>

    <!-- ── Mobile sidebar: shadcn Sheet (drawer) ────────────────────────── -->
    <Sheet v-model:open="sidebarOpen">
      <SheetContent
        side="left"
        class="flex w-72 flex-col p-0 bg-white dark:bg-zinc-900
               border-r border-zinc-900/10 dark:border-white/10"
      >
        <SheetHeader class="sr-only">
          <SheetTitle>Navigation</SheetTitle>
        </SheetHeader>

        <!-- Logo / site title -->
        <div
          class="flex h-14 shrink-0 items-center border-b border-zinc-900/10 dark:border-white/10 px-6"
        >
          <RouteLink
            to="/"
            class="flex items-center gap-2.5 text-sm font-semibold text-zinc-900 dark:text-white hover:opacity-75 transition-opacity"
            @click="sidebarOpen = false"
          >
            <span
              class="flex h-6 w-6 items-center justify-center rounded text-xs font-bold text-white shrink-0"
              style="background: linear-gradient(135deg, #36b49f, #DBFF75);"
              aria-hidden="true"
            >B</span>
            {{ site.title }}
          </RouteLink>
        </div>

        <!-- Nav links -->
        <div class="flex-1 overflow-hidden">
          <Sidebar :config="sidebarConfig" @navigate="sidebarOpen = false" />
        </div>

        <!-- GitHub link -->
        <div
          v-if="themeOptions.repo"
          class="shrink-0 border-t border-zinc-900/10 dark:border-white/10 px-6 py-4"
        >
          <a
            :href="themeOptions.repo"
            target="_blank"
            rel="noopener noreferrer"
            class="flex items-center gap-2 text-sm text-zinc-500 hover:text-zinc-900 dark:text-zinc-400 dark:hover:text-white transition-colors"
          >
            <Github class="h-4 w-4 shrink-0" />
            GitHub
          </a>
        </div>
      </SheetContent>
    </Sheet>

    <!-- ── Desktop sidebar (always visible at lg+) ──────────────────────── -->
    <aside
      class="hidden lg:fixed lg:inset-y-0 lg:left-0 lg:z-50 lg:flex lg:w-72 lg:flex-col
             bg-white dark:bg-zinc-900
             border-r border-zinc-900/10 dark:border-white/10"
    >
      <!-- Logo / site title -->
      <div
        class="flex h-14 shrink-0 items-center border-b border-zinc-900/10 dark:border-white/10 px-6"
      >
        <RouteLink
          to="/"
          class="flex items-center gap-2.5 text-sm font-semibold text-zinc-900 dark:text-white hover:opacity-75 transition-opacity"
        >
          <span
            class="flex h-6 w-6 items-center justify-center rounded text-xs font-bold text-white shrink-0"
            style="background: linear-gradient(135deg, #36b49f, #DBFF75);"
            aria-hidden="true"
          >B</span>
          {{ site.title }}
        </RouteLink>
      </div>

      <!-- Nav links -->
      <div class="flex-1 overflow-hidden">
        <Sidebar :config="sidebarConfig" />
      </div>

      <!-- GitHub link -->
      <div
        v-if="themeOptions.repo"
        class="shrink-0 border-t border-zinc-900/10 dark:border-white/10 px-6 py-4"
      >
        <a
          :href="themeOptions.repo"
          target="_blank"
          rel="noopener noreferrer"
          class="flex items-center gap-2 text-sm text-zinc-500 hover:text-zinc-900 dark:text-zinc-400 dark:hover:text-white transition-colors"
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

      <!-- Content row: article + right TOC -->
      <div class="flex flex-1 justify-center">

        <!-- Article -->
        <main class="min-w-0 flex-1 px-8 py-10 lg:px-10 xl:px-14">
          <article
            class="mx-auto max-w-[752px]
                   prose prose-zinc max-w-none dark:prose-invert
                   prose-a:text-emerald-500 dark:prose-a:text-emerald-400
                   prose-a:no-underline hover:prose-a:underline
                   prose-pre:bg-zinc-900 dark:prose-pre:bg-zinc-800/60
                   prose-pre:ring-1 prose-pre:ring-white/10
                   prose-h1:text-[28px] prose-h1:font-bold prose-h1:tracking-tight prose-h1:mt-0
                   prose-h2:text-[20px] prose-h2:font-semibold prose-h2:tracking-tight prose-h2:mt-12
                   prose-h3:text-[17px] prose-h3:font-semibold prose-h3:tracking-tight prose-h3:mt-8
                   prose-h4:text-[15px] prose-h4:font-semibold prose-h4:mt-6"
            vp-content
          >
            <Content />

            <!-- ── Prev / Next page navigation ──────────────────────────── -->
            <div
              v-if="prevPage || nextPage"
              class="mt-16 flex items-stretch gap-3 border-t border-zinc-900/5 dark:border-white/5 pt-8 not-prose"
            >
              <!-- Prev -->
              <Button
                v-if="prevPage"
                variant="outline"
                as-child
                class="flex-1 h-auto flex-col items-start gap-1 px-5 py-4 text-left
                       ring-zinc-900/10 dark:ring-white/10 hover:ring-zinc-900/20 dark:hover:ring-white/20"
              >
                <RouteLink :to="prevPage.link" class="w-full">
                  <span class="flex items-center gap-1 text-xs text-zinc-400 dark:text-zinc-500">
                    <ArrowLeft class="h-3 w-3" />
                    Previous
                  </span>
                  <span class="mt-1 block text-sm font-medium text-zinc-900 dark:text-white group-hover:text-emerald-500 transition-colors">
                    {{ prevPage.text }}
                  </span>
                </RouteLink>
              </Button>

              <div v-if="!prevPage" class="flex-1" />

              <!-- Next -->
              <Button
                v-if="nextPage"
                variant="outline"
                as-child
                class="flex-1 h-auto flex-col items-end gap-1 px-5 py-4 text-right
                       ring-zinc-900/10 dark:ring-white/10 hover:ring-zinc-900/20 dark:hover:ring-white/20"
              >
                <RouteLink :to="nextPage.link" class="w-full text-right">
                  <span class="flex items-center justify-end gap-1 text-xs text-zinc-400 dark:text-zinc-500">
                    Next
                    <ArrowRight class="h-3 w-3" />
                  </span>
                  <span class="mt-1 block text-sm font-medium text-zinc-900 dark:text-white group-hover:text-emerald-500 transition-colors">
                    {{ nextPage.text }}
                  </span>
                </RouteLink>
              </Button>
            </div>
          </article>
        </main>

        <!-- ── Right TOC — xl+ only ──────────────────────────────────────── -->
        <aside class="hidden xl:block w-56 shrink-0 pr-6">
          <div class="sticky top-[calc(3.5rem+2rem)] pt-1">
            <TableOfContents />
          </div>
        </aside>

      </div>
    </div>
  </div>
</template>
