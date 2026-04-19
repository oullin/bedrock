<script setup lang="ts">
withDefaults(defineProps<{ title: string; description?: string }>(), {
  description: "",
});
</script>

<template>
  <div class="mb-6">
    <h2 class="text-2xl font-semibold tracking-tight">{{ title }}</h2>
    <p v-if="description" class="text-muted-foreground mt-1 text-sm">{{ description }}</p>
  </div>
</template>
