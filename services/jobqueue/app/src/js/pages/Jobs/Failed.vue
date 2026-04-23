<script setup>
import { ref, onMounted } from "vue";
import JobTable from "../../components/JobTable.vue";
import { apiGet } from "../../api.js";

const jobs = ref([]);
const tag = ref("");
const error = ref("");

async function refresh() {
  try {
    const url = tag.value ? `/api/jobs/failed?tag=${encodeURIComponent(tag.value)}` : "/api/jobs/failed";
    jobs.value = await apiGet(url);
    error.value = "";
  } catch (err) {
    error.value = err.message;
  }
}

onMounted(refresh);
</script>

<template>
  <h2>Failed Jobs</h2>
  <div v-if="error" class="error">{{ error }}</div>
  <form @submit.prevent="refresh" style="margin-bottom: 1rem">
    <input v-model="tag" placeholder="Filter by tag" />
    <button type="submit">Filter</button>
  </form>
  <JobTable :jobs="jobs" />
</template>
