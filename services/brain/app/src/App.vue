<script setup lang="ts">
import { onMounted, ref } from "vue";
import type { GraphData, Manifest } from "@/types/graph";

const manifest = ref<Manifest | null>(null);
const graph = ref<GraphData | null>(null);
const error = ref<string | null>(null);
const loading = ref(true);

async function fetchAll() {
  loading.value = true;
  error.value = null;
  try {
    const [m, g] = await Promise.all([
      fetch("/_brain/api/manifest").then((r) => r.json()),
      fetch("/_brain/api/graph").then((r) => r.json()),
    ]);
    manifest.value = m;
    graph.value = g;
  } catch (e) {
    error.value = e instanceof Error ? e.message : String(e);
  } finally {
    loading.value = false;
  }
}

async function rescan() {
  await fetch("/_brain/api/scan", { method: "POST" });
  await fetchAll();
}

onMounted(fetchAll);
</script>

<template>
  <div class="brain-shell">
    <header class="brain-topbar">
      <h1>brain</h1>
      <span v-if="manifest" class="meta">
        {{ manifest.project }} · {{ manifest.totalNodes }} nodes · {{ manifest.totalEdges }} edges · scanned {{ manifest.analyzedAt }}
      </span>
      <span v-else class="meta">loading…</span>
      <button @click="rescan" style="margin-left: auto">Rescan</button>
    </header>
    <main class="brain-main">
      <aside class="brain-sidebar">
        <p v-if="loading" class="meta">Loading…</p>
        <p v-else-if="error">Error: {{ error }}</p>
        <template v-else-if="graph">
          <div v-for="n in graph.nodes" :key="n.id" class="brain-card">
            <div class="type">{{ n.type }}</div>
            <div class="label">{{ n.label }}</div>
          </div>
          <p v-if="graph.nodes.length === 0" class="meta">No nodes yet. Click Rescan after the analyzer pipeline grows.</p>
        </template>
      </aside>
      <section class="brain-viewport">
        <div class="brain-empty">
          <div>
            <p>GraphView lands in phase 9.</p>
            <p>API is live at <code>/_brain/api/*</code>.</p>
          </div>
        </div>
      </section>
    </main>
  </div>
</template>
