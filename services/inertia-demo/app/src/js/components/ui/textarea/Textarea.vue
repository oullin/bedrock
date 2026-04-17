<script setup lang="ts">
import { cn } from "@/js/lib/utils";

const model = defineModel({
  default: "",
});

withDefaults(
  defineProps<{
    class?: any;
  }>(),
  {
    class: "",
  },
);
</script>

<template>
  <textarea
    v-model="model"
    data-slot="textarea"
    :class="
      cn(
        'border-input bg-background ring-offset-background placeholder:text-muted-foreground focus-visible:ring-ring/50 flex min-h-24 w-full rounded-md border px-3 py-2 text-sm shadow-xs transition-[color,box-shadow] outline-none focus-visible:ring-[3px] disabled:cursor-not-allowed disabled:opacity-50',
        $props.class,
      )
    "
  />
</template>
