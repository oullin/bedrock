<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue";
import StatCard from "@/js/components/StatCard.vue";
import { apiGet } from "@/js/api";
import type { Stats, Supervisor } from "@/js/types/horizon";

const stats = ref<Stats | null>(null);
const supervisors = ref<Supervisor[]>([]);
const error = ref("");
let timer: ReturnType<typeof setInterval> | null = null;

async function refresh(): Promise<void> {
  try {
    const [s, sv] = await Promise.all([
      apiGet<Stats>("/api/stats"),
      apiGet<Supervisor[]>("/api/master-supervisors"),
    ]);
    stats.value = s;
    supervisors.value = sv;
    error.value = "";
  } catch (err) {
    error.value = err instanceof Error ? err.message : String(err);
  }
}

function supervisorStatusClass(status: string): string {
  switch (status) {
    case "running":
      return "text-green-500";
    case "paused":
      return "text-yellow-500";
    default:
      return "text-muted-foreground";
  }
}

onMounted(() => {
  refresh();
  timer = setInterval(refresh, 5000);
});

onUnmounted(() => {
  if (timer) clearInterval(timer);
});
</script>

<template>
  <h2 class="text-2xl font-semibold mb-4">Dashboard</h2>
  <div v-if="error" class="bg-destructive text-destructive-foreground px-4 py-3 rounded mb-4">{{ error }}</div>
  <div v-if="stats" class="grid gap-4 grid-cols-[repeat(auto-fit,minmax(200px,1fr))]">
    <StatCard label="Pending Jobs" :value="stats.pendingJobs" />
    <StatCard label="Failed Jobs" :value="stats.failedJobs" />
    <StatCard label="Jobs / Minute" :value="stats.jobsPerMinute" />
    <StatCard label="Status" :value="stats.status" />
    <StatCard label="Max Runtime Queue" :value="stats.queueWithMaxRuntime || 'n/a'" />
    <StatCard label="Max Throughput Queue" :value="stats.queueWithMaxThroughput || 'n/a'" />
  </div>

  <h3 class="mt-8 text-lg font-semibold">Master Supervisors</h3>
  <table v-if="supervisors.length" class="w-full border-collapse mt-4 text-sm">
    <thead>
      <tr>
        <th class="px-3 py-2 text-left text-xs uppercase text-muted-foreground font-medium border-b border-border">Name</th>
        <th class="px-3 py-2 text-left text-xs uppercase text-muted-foreground font-medium border-b border-border">PID</th>
        <th class="px-3 py-2 text-left text-xs uppercase text-muted-foreground font-medium border-b border-border">Status</th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="sv in supervisors" :key="sv.name">
        <td class="px-3 py-2 border-b border-border">{{ sv.name }}</td>
        <td class="px-3 py-2 border-b border-border">{{ sv.pid }}</td>
        <td class="px-3 py-2 border-b border-border" :class="supervisorStatusClass(sv.status)">{{ sv.status }}</td>
      </tr>
    </tbody>
  </table>
  <p v-else class="text-muted-foreground mt-2">No supervisors registered.</p>
</template>
