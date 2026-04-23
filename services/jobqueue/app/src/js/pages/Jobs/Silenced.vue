<script setup>
import { ref, onMounted } from "vue";
import JobTable from "../../components/JobTable.vue";
import { apiGet } from "../../api.js";

const jobs = ref([]);
const error = ref("");

async function refresh() {
  try {
    jobs.value = await apiGet("/api/jobs/silenced");
    error.value = "";
  } catch (err) {
    error.value = err.message;
  }
}

onMounted(refresh);
</script>

<template>
  <h2>Silenced Jobs</h2>
  <div v-if="error" class="error">{{ error }}</div>
  <JobTable :jobs="jobs" />
</template>
