<script setup lang="ts">
import { ref, onMounted } from "vue";
import JobTable from "@/js/components/JobTable.vue";
import { apiGet } from "@/js/api";
import type { Job } from "@/js/types/jobqueue";

const jobs = ref<Job[]>([]);
const tag = ref("");
const error = ref("");

async function refresh(): Promise<void> {
  try {
    const url = tag.value
      ? `/api/jobs/failed?tag=${encodeURIComponent(tag.value)}`
      : "/api/jobs/failed";
    jobs.value = await apiGet<Job[]>(url);
    error.value = "";
  } catch (err) {
    error.value = err instanceof Error ? err.message : String(err);
  }
}

onMounted(refresh);
</script>

<template>
  <h2 class="text-2xl font-semibold mb-4">Failed Jobs</h2>
  <div v-if="error" class="bg-destructive text-destructive-foreground px-4 py-3 rounded mb-4">{{ error }}</div>
  <form @submit.prevent="refresh" class="mb-4 flex gap-2">
    <input
      v-model="tag"
      placeholder="Filter by tag"
      class="bg-input text-foreground border border-border rounded px-3 py-2 text-sm flex-1"
    />
    <button
      type="submit"
      class="bg-primary text-primary-foreground rounded px-3 py-2 text-sm hover:opacity-90"
    >Filter</button>
  </form>
  <JobTable :jobs="jobs" />
</template>
