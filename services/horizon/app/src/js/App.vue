<script setup lang="ts">
import { computed, defineAsyncComponent, type Component } from "vue";
import { currentRoute } from "@/js/router";
import { ROUTES, route } from "@/js/lib/routes";

const pages: Record<string, Component> = {
  [ROUTES.dashboard]:      defineAsyncComponent(() => import("@/js/pages/Dashboard.vue")),
  [ROUTES.monitoring]:     defineAsyncComponent(() => import("@/js/pages/Monitoring.vue")),
  [ROUTES.metrics]:        defineAsyncComponent(() => import("@/js/pages/Metrics.vue")),
  [ROUTES.batches]:        defineAsyncComponent(() => import("@/js/pages/Batches.vue")),
  [ROUTES.jobs.pending]:   defineAsyncComponent(() => import("@/js/pages/Jobs/Pending.vue")),
  [ROUTES.jobs.completed]: defineAsyncComponent(() => import("@/js/pages/Jobs/Completed.vue")),
  [ROUTES.jobs.failed]:    defineAsyncComponent(() => import("@/js/pages/Jobs/Failed.vue")),
  [ROUTES.jobs.silenced]:  defineAsyncComponent(() => import("@/js/pages/Jobs/Silenced.vue")),
};

const page = computed<Component>(() => pages[currentRoute.value.path] ?? pages[ROUTES.dashboard]);
const isActive = (path: string): boolean => currentRoute.value.path === path;

const navLinkClass =
  "block px-3 py-2 rounded text-sm text-sidebar-foreground/80 hover:bg-sidebar-accent hover:text-sidebar-accent-foreground";
const activeClass = "bg-sidebar-accent text-sidebar-accent-foreground";
</script>

<template>
  <div class="dark bg-background text-foreground antialiased min-h-screen">
    <div class="grid grid-cols-[220px_1fr] min-h-screen">
      <aside class="bg-sidebar-background text-sidebar-foreground p-6 border-r border-sidebar-border">
        <h1 class="text-base font-semibold mb-6 text-sidebar-foreground">Bedrock Horizon</h1>
        <nav class="flex flex-col gap-1">
          <a :href="route(ROUTES.dashboard)"      :class="[navLinkClass, isActive(ROUTES.dashboard)      && activeClass]">Dashboard</a>
          <a :href="route(ROUTES.monitoring)"     :class="[navLinkClass, isActive(ROUTES.monitoring)     && activeClass]">Monitoring</a>
          <a :href="route(ROUTES.metrics)"        :class="[navLinkClass, isActive(ROUTES.metrics)        && activeClass]">Metrics</a>
          <a :href="route(ROUTES.batches)"        :class="[navLinkClass, isActive(ROUTES.batches)        && activeClass]">Batches</a>
          <a :href="route(ROUTES.jobs.pending)"   :class="[navLinkClass, isActive(ROUTES.jobs.pending)   && activeClass]">Pending Jobs</a>
          <a :href="route(ROUTES.jobs.completed)" :class="[navLinkClass, isActive(ROUTES.jobs.completed) && activeClass]">Completed Jobs</a>
          <a :href="route(ROUTES.jobs.failed)"    :class="[navLinkClass, isActive(ROUTES.jobs.failed)    && activeClass]">Failed Jobs</a>
          <a :href="route(ROUTES.jobs.silenced)"  :class="[navLinkClass, isActive(ROUTES.jobs.silenced)  && activeClass]">Silenced Jobs</a>
        </nav>
      </aside>
      <main class="p-6 md:px-8">
        <component :is="page" />
      </main>
    </div>
  </div>
</template>
