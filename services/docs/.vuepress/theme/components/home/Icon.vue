<script setup lang="ts">
withDefaults(defineProps<{
  size?: number
  stroke?: number
  fill?: string
  vb?: string
}>(), {
  size: 18,
  stroke: 1.6,
  fill: 'none',
  vb: '0 0 24 24',
})
</script>

<template>
  <svg
    :width="size"
    :height="size"
    :viewBox="vb"
    :fill="fill"
    stroke="currentColor"
    :stroke-width="stroke"
    stroke-linecap="round"
    stroke-linejoin="round"
    aria-hidden="true"
  >
    <slot />
  </svg>
</template>
