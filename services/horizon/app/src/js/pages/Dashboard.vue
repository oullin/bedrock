<script setup>
import { ref, onMounted, onUnmounted } from "vue";
import StatCard from "../components/StatCard.vue";
import { apiGet } from "../api.js";

const stats = ref(null);
const supervisors = ref([]);
const error = ref("");
let timer = null;

async function refresh() {
  try {
    const [s, sv] = await Promise.all([apiGet("/api/stats"), apiGet("/api/master-supervisors")]);
    stats.value = s;
    supervisors.value = sv;
    error.value = "";
  } catch (err) {
    error.value = err.message;
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
  <h2>Dashboard</h2>
  <div v-if="error" class="error">{{ error }}</div>
  <div class="grid" v-if="stats">
    <StatCard label="Pending Jobs" :value="stats.pendingJobs" />
    <StatCard label="Failed Jobs" :value="stats.failedJobs" />
    <StatCard label="Jobs / Minute" :value="stats.jobsPerMinute" />
    <StatCard label="Status" :value="stats.status" />
    <StatCard label="Max Runtime Queue" :value="stats.queueWithMaxRuntime || 'n/a'" />
    <StatCard label="Max Throughput Queue" :value="stats.queueWithMaxThroughput || 'n/a'" />
  </div>

  <h3 style="margin-top: 2rem">Master Supervisors</h3>
  <table v-if="supervisors.length">
    <thead>
      <tr>
        <th>Name</th>
        <th>PID</th>
        <th>Status</th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="sv in supervisors" :key="sv.name">
        <td>{{ sv.name }}</td>
        <td>{{ sv.pid }}</td>
        <td :class="`status-${sv.status}`">{{ sv.status }}</td>
      </tr>
    </tbody>
  </table>
  <p v-else>No supervisors registered.</p>
</template>
