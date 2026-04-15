<script setup lang="ts">
import { type HTMLAttributes } from 'vue'
import { DialogOverlay, type DialogOverlayProps, useForwardProps } from 'radix-vue'
import { cn } from '../../../lib/utils'

interface Props extends DialogOverlayProps {
  class?: HTMLAttributes['class']
}

const props = defineProps<Props>()

const delegatedProps = useForwardProps(props)
</script>

<template>
  <DialogOverlay
    v-bind="delegatedProps"
    :class="cn(
      'fixed inset-0 z-50 bg-black/80 data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0',
      props.class,
    )"
  />
</template>
