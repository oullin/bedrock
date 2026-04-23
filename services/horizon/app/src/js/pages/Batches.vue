<script setup lang="ts">
import { ref, onMounted } from "vue";
import { apiGet } from "@/js/api";
import type { Batch, BatchListResponse } from "@/js/types/horizon";

const batches = ref<Batch[]>([]);
const nextCursor = ref("");
const search = ref("");
const error = ref("");

async function load(after = ""): Promise<void> {
  try {
    const params = new URLSearchParams();
    if (search.value) params.set("name", search.value);
    if (after) params.set("after", after);
    params.set("limit", "25");
    const result = await apiGet<BatchListResponse>(`/api/batches?${params.toString()}`);
    batches.value = after ? batches.value.concat(result.batches) : result.batches;
    nextCursor.value = result.nextCursor || "";
    error.value = "";
  } catch (err) {
    error.value = err instanceof Error ? err.message : String(err);
  }
}

onMounted(() => load());

const thClass = "px-3 py-2 text-left text-xs uppercase text-muted-foreground font-medium border-b border-border";
const tdClass = "px-3 py-2 border-b border-border";
</script>

<template>
  <h2 class="text-2xl font-semibold mb-4">Batches</h2>
  <div v-if="error" class="bg-destructive text-destructive-foreground px-4 py-3 rounded mb-4">{{ error }}</div>
  <form @submit.prevent="load()" class="mb-4 flex gap-2">
    <input
      v-model="search"
      placeholder="Search by name"
      class="bg-input text-foreground border border-border rounded px-3 py-2 text-sm flex-1"
    />
    <button
      type="submit"
      class="bg-primary text-primary-foreground rounded px-3 py-2 text-sm hover:opacity-90"
    >Search</button>
  </form>
  <table v-if="batches.length" class="w-full border-collapse mt-4 text-sm">
    <thead>
      <tr>
        <th :class="thClass">ID</th>
        <th :class="thClass">Name</th>
        <th :class="thClass">Total</th>
        <th :class="thClass">Pending</th>
        <th :class="thClass">Failed</th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="batch in batches" :key="batch.id">
        <td :class="tdClass"><code class="font-mono text-xs">{{ batch.id }}</code></td>
        <td :class="tdClass">{{ batch.name }}</td>
        <td :class="tdClass">{{ batch.totalJobs }}</td>
        <td :class="tdClass">{{ batch.pendingJobs }}</td>
        <td :class="tdClass">{{ batch.failedJobs }}</td>
      </tr>
    </tbody>
  </table>
  <p v-else class="text-muted-foreground mt-2">No batches.</p>
  <button
    v-if="nextCursor"
    @click="load(nextCursor)"
    class="mt-4 bg-secondary text-secondary-foreground rounded px-3 py-2 text-sm hover:opacity-90"
  >Load more</button>
</template>
