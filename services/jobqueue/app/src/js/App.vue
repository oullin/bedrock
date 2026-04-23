<script setup lang="ts">
import { computed, defineAsyncComponent, type Component } from "vue";
import { currentRoute } from "@/js/router";

const routes: Record<string, Component> = {
  "/": defineAsyncComponent(() => import("@/js/pages/Dashboard.vue")),
  "/monitoring": defineAsyncComponent(() => import("@/js/pages/Monitoring.vue")),
  "/metrics": defineAsyncComponent(() => import("@/js/pages/Metrics.vue")),
  "/batches": defineAsyncComponent(() => import("@/js/pages/Batches.vue")),
  "/jobs/pending": defineAsyncComponent(() => import("@/js/pages/Jobs/Pending.vue")),
  "/jobs/completed": defineAsyncComponent(() => import("@/js/pages/Jobs/Completed.vue")),
  "/jobs/failed": defineAsyncComponent(() => import("@/js/pages/Jobs/Failed.vue")),
  "/jobs/silenced": defineAsyncComponent(() => import("@/js/pages/Jobs/Silenced.vue")),
};

const page = computed<Component>(() => routes[currentRoute.value.path] ?? routes["/"]);
const isActive = (path: string): boolean => currentRoute.value.path === path;

const navLinkClass =
  "block px-3 py-2 rounded text-sm text-sidebar-foreground/80 hover:bg-sidebar-accent hover:text-sidebar-accent-foreground";
const activeClass = "bg-sidebar-accent text-sidebar-accent-foreground";
</script>

<template>
  <div class="dark bg-background text-foreground antialiased min-h-screen">
    <div class="grid grid-cols-[220px_1fr] min-h-screen">
      <aside class="bg-sidebar-background text-sidebar-foreground p-6 border-r border-sidebar-border">
        <h1 class="text-base font-semibold mb-6 text-sidebar-foreground">Bedrock JobQueue</h1>
        <nav class="flex flex-col gap-1">
          <a href="#/" :class="[navLinkClass, isActive('/') && activeClass]">Dashboard</a>
          <a href="#/monitoring" :class="[navLinkClass, isActive('/monitoring') && activeClass]">Monitoring</a>
          <a href="#/metrics" :class="[navLinkClass, isActive('/metrics') && activeClass]">Metrics</a>
          <a href="#/batches" :class="[navLinkClass, isActive('/batches') && activeClass]">Batches</a>
          <a href="#/jobs/pending" :class="[navLinkClass, isActive('/jobs/pending') && activeClass]">Pending Jobs</a>
          <a href="#/jobs/completed" :class="[navLinkClass, isActive('/jobs/completed') && activeClass]">Completed Jobs</a>
          <a href="#/jobs/failed" :class="[navLinkClass, isActive('/jobs/failed') && activeClass]">Failed Jobs</a>
          <a href="#/jobs/silenced" :class="[navLinkClass, isActive('/jobs/silenced') && activeClass]">Silenced Jobs</a>
        </nav>
      </aside>
      <main class="p-6 md:px-8">
        <component :is="page" />
      </main>
    </div>
  </div>
</template>
