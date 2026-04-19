import { ref } from "vue";
import { router } from "@inertiajs/vue3";

const isLoading = ref(false);
let timeout = null;
let registered = false;

function register() {
  if (registered) return;
  registered = true;

  router.on("start", (event) => {
    const visit = event.detail.visit;
    if (visit.only && visit.only.length > 0) return;

    timeout = setTimeout(() => {
      isLoading.value = true;
    }, 150);
  });

  router.on("finish", () => {
    clearTimeout(timeout);
    isLoading.value = false;
  });
}

export function useInertiaLoading() {
  register();
  return { isLoading };
}
