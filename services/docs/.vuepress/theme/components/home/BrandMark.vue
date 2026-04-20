<script setup lang="ts">
withDefaults(defineProps<{ size?: number }>(), { size: 28 })
</script>

<template>
  <svg :width="size" :height="size" viewBox="0 0 32 32" fill="none" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
    <defs>
      <linearGradient id="bdrk-g1" x1="0" y1="0" x2="32" y2="32" gradientUnits="userSpaceOnUse">
        <stop offset="0" stop-color="var(--accent)" />
        <stop offset="1" stop-color="var(--accent-2)" />
      </linearGradient>
    </defs>
    <path d="M3 22 L16 29 L29 22 L16 15 Z" fill="url(#bdrk-g1)" />
    <path d="M3 16 L16 23 L29 16 L16 9 Z" fill="var(--accent)" opacity="0.75" />
    <path d="M3 10 L16 17 L29 10 L16 3 Z" fill="var(--ink)" />
    <path d="M3 10 L16 17 L29 10" stroke="var(--accent)" stroke-width="1.2" stroke-linejoin="round" fill="none" opacity="0.9" />
  </svg>
</template>
