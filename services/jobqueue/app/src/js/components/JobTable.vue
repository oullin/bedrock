<script setup>
defineProps({
  jobs: { type: Array, required: true },
});
</script>

<template>
  <table v-if="jobs.length">
    <thead>
      <tr>
        <th>ID</th>
        <th>Name</th>
        <th>Queue</th>
        <th>Status</th>
        <th>Tags</th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="job in jobs" :key="job.id">
        <td><code>{{ job.id }}</code></td>
        <td>{{ job.name }}</td>
        <td>{{ job.queue }}</td>
        <td :class="`status-${job.status}`">{{ job.status }}</td>
        <td>
          <span v-for="tag in job.tags || []" :key="tag" class="tag">{{ tag }}</span>
        </td>
      </tr>
    </tbody>
  </table>
  <p v-else>No jobs.</p>
</template>
