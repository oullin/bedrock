<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute } from 'vue-router'
import { Content, RouteLink, useSiteData } from 'vuepress/client'
import Sidebar from '../Sidebar.vue'
import Navbar from '../Navbar.vue'
import TableOfContents from '../TableOfContents.vue'

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
  <!--
    Root: relative so the absolute decoration is contained.
    bg-white / dark:bg-zinc-900 matches Protocol exactly.
  -->
  <div class="relative flex min-h-full bg-white antialiased dark:bg-zinc-900">

    <!-- ── Protocol gradient + skewed grid decoration ─────────────────────
         Light: mask fades the gradient out toward the bottom.
         Dark:  full-height glow at reduced colour opacity.
    ──────────────────────────────────────────────────────────────────────── -->
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
        <!-- Skewed grid overlay — same pattern as Protocol's GridPattern component -->
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
          <!-- Highlighted cells for visual interest -->
          <rect x="144" y="56"  width="72" height="56" />
          <rect x="288" y="112" width="72" height="56" />
          <rect x="360" y="0"   width="72" height="56" />
          <rect x="576" y="112" width="72" height="56" />
          <rect x="720" y="56"  width="72" height="56" />
          <rect x="864" y="168" width="72" height="56" />
        </svg>
      </div>
    </div>

    <!-- ── Mobile overlay ─────────────────────────────────────────────────── -->
    <Transition
      enter-active-class="transition-opacity duration-200"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition-opacity duration-150"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div
        v-if="sidebarOpen"
        class="fixed inset-0 z-40 bg-zinc-400/20 backdrop-blur-sm dark:bg-black/40 lg:hidden"
        aria-hidden="true"
        @click="sidebarOpen = false"
      />
    </Transition>

    <!-- ── Sidebar ────────────────────────────────────────────────────────── -->
    <aside
      :class="[
        'fixed inset-y-0 left-0 z-50 flex w-72 flex-col',
        'bg-white dark:bg-zinc-900',
        'border-r border-zinc-900/10 dark:border-white/10',
        'transition-transform duration-300 ease-in-out',
        sidebarOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0',
      ]"
    >
      <!-- Logo / site title (same height as navbar) -->
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
      <div class="flex-1 overflow-y-auto py-4">
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
          <svg class="h-4 w-4 shrink-0" fill="currentColor" viewBox="0 0 24 24" aria-hidden="true">
            <path fill-rule="evenodd" d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z" clip-rule="evenodd" />
          </svg>
          GitHub
        </a>
      </div>
    </aside>

    <!-- ── Main content area ───────────────────────────────────────────────
         relative z-10 keeps it above the gradient decoration (z-0).
         No explicit bg — the root bg + decoration show through naturally.
    ──────────────────────────────────────────────────────────────────────── -->
    <div class="relative z-10 flex min-w-0 flex-1 flex-col lg:pl-72">

      <Navbar
        :navbar-items="themeOptions.navbar ?? []"
        @toggle-sidebar="sidebarOpen = !sidebarOpen"
      />

      <!-- Content row: article + right TOC -->
      <div class="flex flex-1 justify-center">

        <!-- Article -->
        <main class="min-w-0 flex-1 px-8 py-10 lg:px-10 xl:px-14">
          <article class="mx-auto max-w-[752px]">
            <Content vp-content />

            <!-- ── Prev / Next page navigation ─────────────────────────── -->
            <div
              v-if="prevPage || nextPage"
              class="mt-16 flex items-stretch gap-3 border-t border-zinc-900/5 dark:border-white/5 pt-8"
            >
              <!-- Prev -->
              <RouteLink
                v-if="prevPage"
                :to="prevPage.link"
                class="group flex flex-1 flex-col gap-1 rounded-xl
                       ring-1 ring-zinc-900/10 dark:ring-white/10
                       px-5 py-4 text-left transition-colors
                       hover:ring-zinc-900/20 dark:hover:ring-white/20
                       hover:bg-zinc-900/2.5 dark:hover:bg-white/5"
              >
                <span class="flex items-center gap-1 text-xs text-zinc-400 dark:text-zinc-500">
                  <svg class="h-3 w-3" fill="none" viewBox="0 0 12 12" aria-hidden="true">
                    <path d="M8.5 3L5 6.5 8.5 10" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" />
                  </svg>
                  Previous
                </span>
                <span class="text-sm font-medium text-zinc-900 dark:text-white group-hover:text-emerald-500 transition-colors">
                  {{ prevPage.text }}
                </span>
              </RouteLink>

              <!-- Spacer when only next exists -->
              <div v-if="!prevPage" class="flex-1" />

              <!-- Next -->
              <RouteLink
                v-if="nextPage"
                :to="nextPage.link"
                class="group flex flex-1 flex-col gap-1 rounded-xl
                       ring-1 ring-zinc-900/10 dark:ring-white/10
                       px-5 py-4 text-right transition-colors
                       hover:ring-zinc-900/20 dark:hover:ring-white/20
                       hover:bg-zinc-900/2.5 dark:hover:bg-white/5"
              >
                <span class="flex items-center justify-end gap-1 text-xs text-zinc-400 dark:text-zinc-500">
                  Next
                  <svg class="h-3 w-3" fill="none" viewBox="0 0 12 12" aria-hidden="true">
                    <path d="M3.5 3L7 6.5 3.5 10" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" />
                  </svg>
                </span>
                <span class="text-sm font-medium text-zinc-900 dark:text-white group-hover:text-emerald-500 transition-colors">
                  {{ nextPage.text }}
                </span>
              </RouteLink>
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
