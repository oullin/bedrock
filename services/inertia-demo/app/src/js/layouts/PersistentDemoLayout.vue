<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue";
import AppLayout from "@/js/layouts/AppLayout.vue";

const counter = ref(0);
const elapsed = ref(0);
let interval: ReturnType<typeof setInterval> | undefined;

onMounted(() => {
  interval = setInterval(() => {
    elapsed.value++;
  }, 1000);
});

onUnmounted(() => {
  clearInterval(interval);
});

function formatTime(seconds: number): string {
  const m = Math.floor(seconds / 60);
  const s = seconds % 60;
  return `${m}:${String(s).padStart(2, "0")}`;
}
</script>

<template>
  <AppLayout title="Persistent Layout Demo">
    <div class="flex h-full flex-1 flex-col gap-4 p-4">
      <div class="flex items-center gap-6 rounded-lg border bg-muted/30 p-4">
        <div class="text-sm">
          <span class="text-muted-foreground">Stopwatch:</span>
          <span class="ml-1 font-mono font-medium">{{ formatTime(elapsed) }}</span>
        </div>
        <div class="text-sm">
          <span class="text-muted-foreground">Counter:</span>
          <span class="ml-1 font-mono font-medium">{{ counter }}</span>
          <button class="text-primary ml-2 text-xs underline" @click="counter++">increment</button>
        </div>
        <div class="text-muted-foreground text-xs">
          These values persist across page navigations because the layout is not remounted.
        </div>
      </div>

      <slot />
    </div>
  </AppLayout>
</template>
