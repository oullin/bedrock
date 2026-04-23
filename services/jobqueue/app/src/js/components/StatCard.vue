<script setup lang="ts">
defineProps<{
  label: string;
  value: string | number;
}>();
</script>

<template>
  <div class="bg-card text-card-foreground border border-border rounded-md p-4">
    <div class="text-xs uppercase tracking-wider text-muted-foreground">{{ label }}</div>
    <div class="mt-1 text-3xl font-semibold">{{ value }}</div>
  </div>
</template>
