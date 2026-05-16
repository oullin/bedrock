<script setup lang="ts">
import type { GraphData, NodeType } from "@/types/graph";
import { computed } from "vue";
import { nodeColor } from "@/lib/colors";

const props = defineProps<{
  graph: GraphData;
  enabledTypes: Set<NodeType>;
}>();
const emit = defineEmits<{ (e: "toggle", t: NodeType): void }>();

const typeCounts = computed(() => {
  const c = new Map<NodeType, number>();
  for (const n of props.graph.nodes) {
    c.set(n.type, (c.get(n.type) ?? 0) + 1);
  }
  return [...c.entries()].sort(([a], [b]) => a.localeCompare(b));
});
</script>

<template>
  <div class="filter-sidebar">
    <h2>Node types</h2>
    <ul>
      <li v-for="[t, n] in typeCounts" :key="t">
        <label>
          <input
            type="checkbox"
            :checked="enabledTypes.has(t)"
            @change="emit('toggle', t)"
          />
          <span class="swatch" :style="{ background: nodeColor(t) }" />
          <span class="name">{{ t }}</span>
          <span class="count">{{ n }}</span>
        </label>
      </li>
    </ul>
  </div>
</template>

<style scoped>
.filter-sidebar h2 {
  font-size: 12px;
  color: #7d8590;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  margin: 0 0 8px;
}
ul { list-style: none; padding: 0; margin: 0; }
li { margin-bottom: 4px; }
label { display: flex; align-items: center; gap: 8px; cursor: pointer; font-size: 12px; }
.swatch {
  width: 10px;
  height: 10px;
  border-radius: 2px;
  flex: none;
}
.name { flex: 1; }
.count { color: #7d8590; font-family: monospace; }
</style>
