<script setup lang="ts">
import { ref, onMounted } from "vue";
import { apiGet, apiPost } from "@/js/api";
import type { MetricRow } from "@/js/types/horizon";

const jobs = ref<MetricRow[]>([]);
const queues = ref<MetricRow[]>([]);
const error = ref("");

async function refresh(): Promise<void> {
  try {
    const [j, q] = await Promise.all([
      apiGet<MetricRow[]>("/api/metrics/jobs"),
      apiGet<MetricRow[]>("/api/metrics/queues"),
    ]);
    jobs.value = j;
    queues.value = q;
    error.value = "";
  } catch (err) {
    error.value = err instanceof Error ? err.message : String(err);
  }
}

async function snapshot(): Promise<void> {
  await apiPost("/api/metrics/snapshot");
  refresh();
}

onMounted(refresh);

const thClass = "px-3 py-2 text-left text-xs uppercase text-muted-foreground font-medium border-b border-border";
const tdClass = "px-3 py-2 border-b border-border";
</script>

<template>
  <h2 class="text-2xl font-semibold mb-4 flex items-center gap-3">
    Metrics
    <button
      @click="snapshot"
      class="bg-primary text-primary-foreground rounded px-3 py-1.5 text-sm hover:opacity-90"
    >Snapshot</button>
  </h2>
  <div v-if="error" class="bg-destructive text-destructive-foreground px-4 py-3 rounded mb-4">{{ error }}</div>
  <h3 class="text-lg font-semibold mt-4">Jobs</h3>
  <table v-if="jobs.length" class="w-full border-collapse mt-4 text-sm">
    <thead>
      <tr>
        <th :class="thClass">Job</th>
        <th :class="thClass">Throughput</th>
        <th :class="thClass">Avg Runtime (ms)</th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="row in jobs" :key="row.name">
        <td :class="tdClass">{{ row.name }}</td>
        <td :class="tdClass">{{ row.throughput }}</td>
        <td :class="tdClass">{{ row.averageRuntime }}</td>
      </tr>
    </tbody>
  </table>
  <p v-else class="text-muted-foreground mt-2">No job metrics recorded.</p>

  <h3 class="mt-8 text-lg font-semibold">Queues</h3>
  <table v-if="queues.length" class="w-full border-collapse mt-4 text-sm">
    <thead>
      <tr>
        <th :class="thClass">Queue</th>
        <th :class="thClass">Throughput</th>
        <th :class="thClass">Avg Runtime (ms)</th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="row in queues" :key="row.name">
        <td :class="tdClass">{{ row.name }}</td>
        <td :class="tdClass">{{ row.throughput }}</td>
        <td :class="tdClass">{{ row.averageRuntime }}</td>
      </tr>
    </tbody>
  </table>
  <p v-else class="text-muted-foreground mt-2">No queue metrics yet — trigger a snapshot after recording jobs.</p>
</template>
