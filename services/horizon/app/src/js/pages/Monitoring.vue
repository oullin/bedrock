<script setup lang="ts">
import { ref, onMounted } from "vue";
import { apiGet, apiPost, apiDelete } from "@/js/api";
import type { MonitoringEntry } from "@/js/types/horizon";

const tags = ref<MonitoringEntry[]>([]);
const newTag = ref("");
const error = ref("");

async function refresh(): Promise<void> {
  try {
    tags.value = await apiGet<MonitoringEntry[]>("/api/monitoring");
    error.value = "";
  } catch (err) {
    error.value = err instanceof Error ? err.message : String(err);
  }
}

async function addTag(): Promise<void> {
  if (!newTag.value.trim()) return;
  await apiPost("/api/monitoring", { tag: newTag.value.trim() });
  newTag.value = "";
  refresh();
}

async function removeTag(tag: string): Promise<void> {
  await apiDelete(`/api/monitoring/${encodeURIComponent(tag)}`);
  refresh();
}

onMounted(refresh);
</script>

<template>
  <h2 class="text-2xl font-semibold mb-4">Monitoring</h2>
  <div v-if="error" class="bg-destructive text-destructive-foreground px-4 py-3 rounded mb-4">{{ error }}</div>
  <form @submit.prevent="addTag" class="mb-4 flex gap-2">
    <input
      v-model="newTag"
      placeholder="Tag to monitor"
      class="bg-input text-foreground border border-border rounded px-3 py-2 text-sm flex-1"
    />
    <button
      type="submit"
      class="bg-primary text-primary-foreground rounded px-3 py-2 text-sm hover:opacity-90"
    >Monitor</button>
  </form>
  <table v-if="tags.length" class="w-full border-collapse mt-4 text-sm">
    <thead>
      <tr>
        <th class="px-3 py-2 text-left text-xs uppercase text-muted-foreground font-medium border-b border-border">Tag</th>
        <th class="px-3 py-2 text-left text-xs uppercase text-muted-foreground font-medium border-b border-border">Pending Jobs</th>
        <th class="px-3 py-2 text-left text-xs uppercase text-muted-foreground font-medium border-b border-border"></th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="entry in tags" :key="entry.tag">
        <td class="px-3 py-2 border-b border-border">{{ entry.tag }}</td>
        <td class="px-3 py-2 border-b border-border">{{ entry.count }}</td>
        <td class="px-3 py-2 border-b border-border">
          <button
            @click="removeTag(entry.tag)"
            class="bg-secondary text-secondary-foreground rounded px-2 py-1 text-xs hover:opacity-90"
          >Stop</button>
        </td>
      </tr>
    </tbody>
  </table>
  <p v-else class="text-muted-foreground mt-2">No tags monitored.</p>
</template>
