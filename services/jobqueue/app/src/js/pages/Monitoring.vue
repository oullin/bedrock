<script setup>
import { ref, onMounted } from "vue";
import { apiGet, apiPost, apiDelete } from "../api.js";

const tags = ref([]);
const newTag = ref("");
const error = ref("");

async function refresh() {
  try {
    tags.value = await apiGet("/api/monitoring");
    error.value = "";
  } catch (err) {
    error.value = err.message;
  }
}

async function addTag() {
  if (!newTag.value.trim()) return;
  await apiPost("/api/monitoring", { tag: newTag.value.trim() });
  newTag.value = "";
  refresh();
}

async function removeTag(tag) {
  await apiDelete(`/api/monitoring/${encodeURIComponent(tag)}`);
  refresh();
}

onMounted(refresh);
</script>

<template>
  <h2>Monitoring</h2>
  <div v-if="error" class="error">{{ error }}</div>
  <form @submit.prevent="addTag" style="margin-bottom: 1rem">
    <input v-model="newTag" placeholder="Tag to monitor" />
    <button type="submit">Monitor</button>
  </form>
  <table v-if="tags.length">
    <thead>
      <tr>
        <th>Tag</th>
        <th>Pending Jobs</th>
        <th></th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="entry in tags" :key="entry.tag">
        <td>{{ entry.tag }}</td>
        <td>{{ entry.count }}</td>
        <td><button @click="removeTag(entry.tag)">Stop</button></td>
      </tr>
    </tbody>
  </table>
  <p v-else>No tags monitored.</p>
</template>
