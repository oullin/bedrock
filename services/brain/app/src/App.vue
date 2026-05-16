<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import type { GraphData, GraphNode, Manifest, NodeType } from "@/types/graph";
import GraphView from "@/components/GraphView.vue";
import FilterSidebar from "@/components/FilterSidebar.vue";
import NodeDetailsModal from "@/components/NodeDetailsModal.vue";

const manifest = ref<Manifest | null>(null);
const graph = ref<GraphData | null>(null);
const error = ref<string | null>(null);
const loading = ref(true);
const selected = ref<GraphNode | null>(null);
const enabledTypes = ref<Set<NodeType>>(new Set());

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
    enabledTypes.value = new Set(g.nodes.map((n: GraphNode) => n.type));
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

const filteredGraph = computed<GraphData | null>(() => {
  if (!graph.value) return null;
  const allowed = enabledTypes.value;
  const nodes = graph.value.nodes.filter((n) => allowed.has(n.type));
  const idSet = new Set(nodes.map((n) => n.id));
  const edges = graph.value.edges.filter((e) => idSet.has(e.source) && idSet.has(e.target));
  return { ...graph.value, nodes, edges, meta: { ...graph.value.meta, nodeCount: nodes.length, edgeCount: edges.length } };
});

function toggleType(t: NodeType) {
  const s = new Set(enabledTypes.value);
  if (s.has(t)) s.delete(t);
  else s.add(t);
  enabledTypes.value = s;
}

onMounted(fetchAll);
</script>

<template>
  <div class="brain-shell">
    <header class="brain-topbar">
      <h1>brain</h1>
      <span v-if="manifest" class="meta">
        {{ manifest.project }} · {{ filteredGraph?.meta.nodeCount ?? 0 }}/{{ manifest.totalNodes }} nodes · {{ filteredGraph?.meta.edgeCount ?? 0 }}/{{ manifest.totalEdges }} edges
      </span>
      <span v-else class="meta">loading…</span>
      <button @click="rescan" style="margin-left: auto">Rescan</button>
    </header>
    <main class="brain-main">
      <aside class="brain-sidebar">
        <p v-if="loading" class="meta">Loading…</p>
        <p v-else-if="error">Error: {{ error }}</p>
        <FilterSidebar
          v-else-if="graph"
          :graph="graph"
          :enabled-types="enabledTypes"
          @toggle="toggleType"
        />
      </aside>
      <section class="brain-viewport">
        <GraphView
          v-if="filteredGraph && filteredGraph.nodes.length > 0"
          :graph="filteredGraph"
          @select="(n) => (selected = n)"
        />
        <div v-else class="brain-empty">
          <p v-if="!graph">Loading graph…</p>
          <p v-else>No nodes match the current filter.</p>
        </div>
      </section>
    </main>
    <NodeDetailsModal :node="selected" @close="selected = null" />
  </div>
</template>
