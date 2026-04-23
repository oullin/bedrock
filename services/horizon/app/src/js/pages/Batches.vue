<script setup>
import { ref, onMounted } from "vue";
import { apiGet } from "../api.js";

const batches = ref([]);
const nextCursor = ref("");
const search = ref("");
const error = ref("");

async function load(after = "") {
  try {
    const params = new URLSearchParams();
    if (search.value) params.set("name", search.value);
    if (after) params.set("after", after);
    params.set("limit", "25");
    const result = await apiGet(`/api/batches?${params.toString()}`);
    batches.value = after ? batches.value.concat(result.batches) : result.batches;
    nextCursor.value = result.nextCursor || "";
    error.value = "";
  } catch (err) {
    error.value = err.message;
  }
}

onMounted(() => load());
</script>

<template>
  <h2>Batches</h2>
  <div v-if="error" class="error">{{ error }}</div>
  <form @submit.prevent="load()" style="margin-bottom: 1rem">
    <input v-model="search" placeholder="Search by name" />
    <button type="submit">Search</button>
  </form>
  <table v-if="batches.length">
    <thead>
      <tr>
        <th>ID</th>
        <th>Name</th>
        <th>Total</th>
        <th>Pending</th>
        <th>Failed</th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="batch in batches" :key="batch.id">
        <td><code>{{ batch.id }}</code></td>
        <td>{{ batch.name }}</td>
        <td>{{ batch.totalJobs }}</td>
        <td>{{ batch.pendingJobs }}</td>
        <td>{{ batch.failedJobs }}</td>
      </tr>
    </tbody>
  </table>
  <p v-else>No batches.</p>
  <button v-if="nextCursor" @click="load(nextCursor)">Load more</button>
</template>
