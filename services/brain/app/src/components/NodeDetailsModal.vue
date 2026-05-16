<script setup lang="ts">
import type { GraphNode } from "@/types/graph";
import { computed } from "vue";

const props = defineProps<{ node: GraphNode | null }>();
const emit = defineEmits<{ (e: "close"): void }>();

const rows = computed(() => {
  if (!props.node) return [] as { k: string; v: string }[];
  return Object.entries(props.node.data).map(([k, v]) => ({
    k,
    v: typeof v === "string" ? v : JSON.stringify(v),
  }));
});
</script>

<template>
  <div v-if="node" class="modal-overlay" @click.self="emit('close')">
    <div class="modal-card">
      <header>
        <div class="type">{{ node.type }}</div>
        <div class="label">{{ node.label }}</div>
        <button class="close" @click="emit('close')" aria-label="Close">×</button>
      </header>
      <dl>
        <template v-for="row in rows" :key="row.k">
          <dt>{{ row.k }}</dt>
          <dd>{{ row.v }}</dd>
        </template>
      </dl>
    </div>
  </div>
</template>

<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.55);
  display: grid;
  place-items: center;
  z-index: 50;
}
.modal-card {
  background: #161b22;
  border: 1px solid #30363d;
  border-radius: 8px;
  min-width: 360px;
  max-width: 600px;
  max-height: 80vh;
  overflow: auto;
  padding: 16px;
}
header { display: flex; align-items: center; gap: 8px; margin-bottom: 12px; }
.type { color: #7d8590; font-size: 11px; text-transform: uppercase; letter-spacing: 0.04em; }
.label { font-weight: 600; font-size: 16px; flex: 1; }
.close {
  background: none; border: none; color: #e6edf3; font-size: 24px;
  cursor: pointer; line-height: 1; padding: 0 4px;
}
dl { display: grid; grid-template-columns: 110px 1fr; gap: 4px 12px; margin: 0; }
dt { color: #7d8590; font-family: monospace; }
dd { margin: 0; word-break: break-all; font-family: monospace; font-size: 12px; }
</style>
