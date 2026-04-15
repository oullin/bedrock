<script setup lang="ts">
import { type HTMLAttributes } from 'vue'
import { Separator, type SeparatorProps } from 'radix-vue'
import { cn } from '../../../lib/utils'

interface Props extends SeparatorProps {
  class?: HTMLAttributes['class']
}

const props = withDefaults(defineProps<Props>(), {
  orientation: 'horizontal',
  decorative: true,
})
</script>

<template>
  <Separator
    v-bind="props"
    :class="cn(
      'shrink-0 bg-border',
      props.orientation === 'horizontal' ? 'h-px w-full' : 'h-full w-px',
      props.class,
    )"
  />
</template>
