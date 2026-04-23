<script setup lang="ts">
import type { Job, JobStatus } from "@/js/types/horizon";

defineProps<{
  jobs: Job[];
}>();

function statusClass(status: JobStatus): string {
  switch (status) {
    case "running":
    case "completed":
      return "text-green-500";
    case "paused":
    case "pending":
      return "text-yellow-500";
    case "failed":
      return "text-red-500";
    default:
      return "text-muted-foreground";
  }
}
</script>

<template>
  <table v-if="jobs.length" class="w-full border-collapse mt-4 text-sm">
    <thead>
      <tr>
        <th class="px-3 py-2 text-left text-xs uppercase text-muted-foreground font-medium border-b border-border">ID</th>
        <th class="px-3 py-2 text-left text-xs uppercase text-muted-foreground font-medium border-b border-border">Name</th>
        <th class="px-3 py-2 text-left text-xs uppercase text-muted-foreground font-medium border-b border-border">Queue</th>
        <th class="px-3 py-2 text-left text-xs uppercase text-muted-foreground font-medium border-b border-border">Status</th>
        <th class="px-3 py-2 text-left text-xs uppercase text-muted-foreground font-medium border-b border-border">Tags</th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="job in jobs" :key="job.id">
        <td class="px-3 py-2 border-b border-border"><code class="font-mono text-xs">{{ job.id }}</code></td>
        <td class="px-3 py-2 border-b border-border">{{ job.name }}</td>
        <td class="px-3 py-2 border-b border-border">{{ job.queue }}</td>
        <td class="px-3 py-2 border-b border-border" :class="statusClass(job.status)">{{ job.status }}</td>
        <td class="px-3 py-2 border-b border-border">
          <span
            v-for="tag in job.tags || []"
            :key="tag"
            class="inline-block bg-muted px-2 py-0.5 rounded text-xs mr-1"
          >{{ tag }}</span>
        </td>
      </tr>
    </tbody>
  </table>
  <p v-else class="text-muted-foreground mt-4">No jobs.</p>
</template>
