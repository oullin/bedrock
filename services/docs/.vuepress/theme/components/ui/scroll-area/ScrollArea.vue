<script setup lang="ts">
import { type HTMLAttributes } from 'vue'
import { ScrollAreaRoot, ScrollAreaViewport, ScrollAreaCorner, type ScrollAreaRootProps } from 'radix-vue'
import ScrollBar from './ScrollBar.vue'
import { cn } from '../../../lib/utils'

interface Props extends ScrollAreaRootProps {
  class?: HTMLAttributes['class']
}

const props = defineProps<Props>()
</script>

<template>
  <ScrollAreaRoot
    :class="cn('relative overflow-hidden', props.class)"
    v-bind="props"
  >
    <ScrollAreaViewport class="h-full w-full rounded-[inherit]">
      <slot />
    </ScrollAreaViewport>
    <ScrollBar />
    <ScrollAreaCorner />
  </ScrollAreaRoot>
</template>
