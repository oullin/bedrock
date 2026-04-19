<script setup>
import { nextTick, watch, ref } from "vue";
import { useInertiaLoading } from "@/js/composables/useInertiaLoading";

const { isLoading } = useInertiaLoading();

const width = ref("0%");
const opacity = ref(0);
const transition = ref("none");

let finishTimer = null;
let fadeTimer = null;

function clearTimers() {
  clearTimeout(finishTimer);
  clearTimeout(fadeTimer);
}

watch(isLoading, async (loading) => {
  clearTimers();

  if (loading) {
    transition.value = "none";
    width.value = "0%";
    opacity.value = 1;

    await nextTick();

    transition.value = "width 2s cubic-bezier(0.4, 0.0, 0.2, 1)";
    width.value = "80%";
  } else {
    transition.value = "width 200ms ease-out";
    width.value = "100%";

    finishTimer = setTimeout(() => {
      transition.value = "opacity 300ms ease-out";
      opacity.value = 0;

      fadeTimer = setTimeout(() => {
        transition.value = "none";
        width.value = "0%";
      }, 300);
    }, 200);
  }
});
</script>

<template>
  <div class="pointer-events-none fixed inset-x-0 top-0 z-50 h-0.5">
    <div class="h-full bg-primary" :style="{ width, opacity, transition }" />
  </div>
</template>
