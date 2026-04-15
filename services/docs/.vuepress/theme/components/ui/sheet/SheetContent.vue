<script setup lang="ts">
import { type HTMLAttributes } from 'vue'
import { DialogContent, type DialogContentEmits, type DialogContentProps, DialogPortal, useForwardPropsEmits } from 'radix-vue'
import SheetOverlay from './SheetOverlay.vue'
import { cn } from '../../../lib/utils'
import { cva } from 'class-variance-authority'

interface Props extends DialogContentProps {
  class?: HTMLAttributes['class']
  side?: 'top' | 'right' | 'bottom' | 'left'
}

const sheetVariants = cva(
  'fixed z-50 gap-4 bg-background p-6 shadow-lg transition ease-in-out data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:duration-300 data-[state=open]:duration-500',
  {
    variants: {
      side: {
        top: 'inset-x-0 top-0 border-b data-[state=closed]:slide-out-to-top data-[state=open]:slide-in-from-top',
        bottom: 'inset-x-0 bottom-0 border-t data-[state=closed]:slide-out-to-bottom data-[state=open]:slide-in-from-bottom',
        left: 'inset-y-0 left-0 h-full w-72 border-r data-[state=closed]:slide-out-to-left data-[state=open]:slide-in-from-left',
        right: 'inset-y-0 right-0 h-full w-72 border-l data-[state=closed]:slide-out-to-right data-[state=open]:slide-in-from-right',
      },
    },
    defaultVariants: {
      side: 'right',
    },
  },
)

const props = withDefaults(defineProps<Props>(), {
  side: 'right',
})
const emits = defineEmits<DialogContentEmits>()

const delegatedProps = useForwardPropsEmits(props, emits)
</script>

<template>
  <DialogPortal>
    <SheetOverlay />
    <DialogContent
      :class="cn(sheetVariants({ side }), props.class)"
      v-bind="delegatedProps"
    >
      <slot />
    </DialogContent>
  </DialogPortal>
</template>
