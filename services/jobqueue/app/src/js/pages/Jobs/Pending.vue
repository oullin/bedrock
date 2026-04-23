<script setup>
import { ref, onMounted } from "vue";
import JobTable from "../../components/JobTable.vue";
import { apiGet } from "../../api.js";

const jobs = ref([]);
const queue = ref("default");
const error = ref("");

async function refresh() {
  try {
    jobs.value = await apiGet(`/api/jobs/pending?queue=${encodeURIComponent(queue.value)}`);
    error.value = "";
  } catch (err) {
    error.value = err.message;
  }
}

onMounted(refresh);
</script>

<template>
  <h2>Pending Jobs</h2>
  <div v-if="error" class="error">{{ error }}</div>
  <form @submit.prevent="refresh" style="margin-bottom: 1rem">
    <input v-model="queue" placeholder="Queue" />
    <button type="submit">Filter</button>
  </form>
  <JobTable :jobs="jobs" />
</template>
