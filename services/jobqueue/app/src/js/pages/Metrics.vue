<script setup>
import { ref, onMounted } from "vue";
import { apiGet, apiPost } from "../api.js";

const jobs = ref([]);
const queues = ref([]);
const error = ref("");

async function refresh() {
  try {
    const [j, q] = await Promise.all([apiGet("/api/metrics/jobs"), apiGet("/api/metrics/queues")]);
    jobs.value = j;
    queues.value = q;
    error.value = "";
  } catch (err) {
    error.value = err.message;
  }
}

async function snapshot() {
  await apiPost("/api/metrics/snapshot");
  refresh();
}

onMounted(refresh);
</script>

<template>
  <h2>Metrics <button @click="snapshot">Snapshot</button></h2>
  <div v-if="error" class="error">{{ error }}</div>
  <h3>Jobs</h3>
  <table v-if="jobs.length">
    <thead>
      <tr>
        <th>Job</th>
        <th>Throughput</th>
        <th>Avg Runtime (ms)</th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="row in jobs" :key="row.name">
        <td>{{ row.name }}</td>
        <td>{{ row.throughput }}</td>
        <td>{{ row.averageRuntime }}</td>
      </tr>
    </tbody>
  </table>
  <p v-else>No job metrics recorded.</p>

  <h3 style="margin-top: 2rem">Queues</h3>
  <table v-if="queues.length">
    <thead>
      <tr>
        <th>Queue</th>
        <th>Throughput</th>
        <th>Avg Runtime (ms)</th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="row in queues" :key="row.name">
        <td>{{ row.name }}</td>
        <td>{{ row.throughput }}</td>
        <td>{{ row.averageRuntime }}</td>
      </tr>
    </tbody>
  </table>
  <p v-else>No queue metrics yet — trigger a snapshot after recording jobs.</p>
</template>
