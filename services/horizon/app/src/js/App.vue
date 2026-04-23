<script setup>
import { computed, defineAsyncComponent } from "vue";
import { currentRoute } from "./router.js";

const routes = {
  "/": defineAsyncComponent(() => import("./pages/Dashboard.vue")),
  "/monitoring": defineAsyncComponent(() => import("./pages/Monitoring.vue")),
  "/metrics": defineAsyncComponent(() => import("./pages/Metrics.vue")),
  "/batches": defineAsyncComponent(() => import("./pages/Batches.vue")),
  "/jobs/pending": defineAsyncComponent(() => import("./pages/Jobs/Pending.vue")),
  "/jobs/completed": defineAsyncComponent(() => import("./pages/Jobs/Completed.vue")),
  "/jobs/failed": defineAsyncComponent(() => import("./pages/Jobs/Failed.vue")),
  "/jobs/silenced": defineAsyncComponent(() => import("./pages/Jobs/Silenced.vue")),
};

const page = computed(() => routes[currentRoute.value.path] || routes["/"]);
const isActive = (path) => currentRoute.value.path === path;
</script>

<template>
  <div class="layout">
    <aside class="sidebar">
      <h1>Bedrock Horizon</h1>
      <nav>
        <a href="#/" :class="{ active: isActive('/') }">Dashboard</a>
        <a href="#/monitoring" :class="{ active: isActive('/monitoring') }">Monitoring</a>
        <a href="#/metrics" :class="{ active: isActive('/metrics') }">Metrics</a>
        <a href="#/batches" :class="{ active: isActive('/batches') }">Batches</a>
        <a href="#/jobs/pending" :class="{ active: isActive('/jobs/pending') }">Pending Jobs</a>
        <a href="#/jobs/completed" :class="{ active: isActive('/jobs/completed') }">Completed Jobs</a>
        <a href="#/jobs/failed" :class="{ active: isActive('/jobs/failed') }">Failed Jobs</a>
        <a href="#/jobs/silenced" :class="{ active: isActive('/jobs/silenced') }">Silenced Jobs</a>
      </nav>
    </aside>
    <main class="main">
      <component :is="page" />
    </main>
  </div>
</template>
