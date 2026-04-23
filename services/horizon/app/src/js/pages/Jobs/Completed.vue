<script setup lang="ts">
import { ref, onMounted } from "vue";
import JobTable from "@/js/components/JobTable.vue";
import { apiGet } from "@/js/api";
import type { Job } from "@/js/types/horizon";

const jobs = ref<Job[]>([]);
const error = ref("");

async function refresh(): Promise<void> {
  try {
    jobs.value = await apiGet<Job[]>("/api/jobs/completed?limit=50");
    error.value = "";
  } catch (err) {
    error.value = err instanceof Error ? err.message : String(err);
  }
}

onMounted(refresh);
</script>

<template>
  <h2 class="text-2xl font-semibold mb-4">Completed Jobs</h2>
  <div v-if="error" class="bg-destructive text-destructive-foreground px-4 py-3 rounded mb-4">{{ error }}</div>
  <JobTable :jobs="jobs" />
</template>
