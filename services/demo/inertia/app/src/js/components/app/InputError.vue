<script setup lang="ts">
withDefaults(defineProps<{ message?: string }>(), {
  message: "",
});
</script>

<template>
  <p v-if="message" role="alert" class="text-destructive text-sm">
    {{ message }}
  </p>
</template>
