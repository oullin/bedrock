<script setup lang="ts">
import { Head } from "@inertiajs/vue3";

defineProps<{ title?: string }>();
</script>

<template>
  <div class="min-h-screen flex flex-col items-center bg-gray-100 pt-6 sm:justify-center sm:pt-0">
    <Head :title="title" />

    <div class="mb-6">
      <a href="/" class="text-3xl font-bold text-gray-900">Bedrock</a>
    </div>

    <div class="w-full overflow-hidden bg-white px-6 py-4 shadow-md sm:max-w-md sm:rounded-lg">
      <slot />
    </div>
  </div>
</template>
