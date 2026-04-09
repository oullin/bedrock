<script setup lang="ts">
import AppLayout from "@/Layouts/AppLayout.vue";
</script>

<template>
  <AppLayout title="Dashboard">
    <template #header>
      <h2 class="text-xl font-semibold leading-tight text-gray-800">Dashboard</h2>
    </template>

    <div class="rounded-lg bg-white p-6 shadow">
      <p class="text-gray-600">You're logged in!</p>
    </div>
  </AppLayout>
</template>
